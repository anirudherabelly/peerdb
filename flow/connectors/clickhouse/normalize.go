package connclickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PeerDB-io/peer-flow/generated/protos"
	"github.com/PeerDB-io/peer-flow/model"
	"github.com/PeerDB-io/peer-flow/model/numeric"
	"github.com/PeerDB-io/peer-flow/model/qvalue"
)

const (
	signColName    = "_peerdb_is_deleted"
	signColType    = "Int8"
	versionColName = "_peerdb_version"
	versionColType = "Int64"
)

func (c *ClickhouseConnector) StartSetupNormalizedTables(_ context.Context) (interface{}, error) {
	return nil, nil
}

func (c *ClickhouseConnector) FinishSetupNormalizedTables(_ context.Context, _ interface{}) error {
	return nil
}

func (c *ClickhouseConnector) CleanupSetupNormalizedTables(_ context.Context, _ interface{}) {
}

func (c *ClickhouseConnector) SetupNormalizedTable(
	ctx context.Context,
	tx interface{},
	tableIdentifier string,
	tableSchema *protos.TableSchema,
	softDeleteColName string,
	syncedAtColName string,
) (bool, error) {
	tableAlreadyExists, err := c.checkIfTableExists(ctx, c.config.Database, tableIdentifier)
	if err != nil {
		return false, fmt.Errorf("error occurred while checking if normalized table exists: %w", err)
	}
	if tableAlreadyExists {
		return true, nil
	}

	normalizedTableCreateSQL, err := generateCreateTableSQLForNormalizedTable(
		tableIdentifier,
		tableSchema,
		softDeleteColName,
		syncedAtColName,
	)
	if err != nil {
		return false, fmt.Errorf("error while generating create table sql for normalized table: %w", err)
	}

	_, err = c.database.ExecContext(ctx, normalizedTableCreateSQL)
	if err != nil {
		return false, fmt.Errorf("[ch] error while creating normalized table: %w", err)
	}
	return false, nil
}

func generateCreateTableSQLForNormalizedTable(
	normalizedTable string,
	tableSchema *protos.TableSchema,
	_ string, // softDeleteColName
	syncedAtColName string,
) (string, error) {
	var stmtBuilder strings.Builder
	stmtBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", normalizedTable))

	for _, column := range tableSchema.Columns {
		colName := column.Name
		colType := qvalue.QValueKind(column.Type)
		clickhouseType, err := qValueKindToClickhouseType(colType)
		if err != nil {
			return "", fmt.Errorf("error while converting column type to clickhouse type: %w", err)
		}

		switch colType {
		case qvalue.QValueKindNumeric:
			precision, scale := numeric.ParseNumericTypmod(column.TypeModifier)
			if column.TypeModifier == -1 || precision > 76 || scale > precision {
				precision = numeric.PeerDBClickhousePrecision
				scale = numeric.PeerDBClickhouseScale
			}
			stmtBuilder.WriteString(fmt.Sprintf("`%s` DECIMAL(%d, %d), ",
				colName, precision, scale))
		default:
			stmtBuilder.WriteString(fmt.Sprintf("`%s` %s, ", colName, clickhouseType))
		}
	}
	// TODO support soft delete
	// synced at column will be added to all normalized tables
	if syncedAtColName != "" {
		colName := strings.ToLower(syncedAtColName)
		stmtBuilder.WriteString(fmt.Sprintf("`%s` %s, ", colName, "DateTime64(9) DEFAULT now()"))
	}

	// add sign and version columns
	stmtBuilder.WriteString(fmt.Sprintf("`%s` %s, ", signColName, signColType))
	stmtBuilder.WriteString(fmt.Sprintf("`%s` %s", versionColName, versionColType))

	stmtBuilder.WriteString(fmt.Sprintf(") ENGINE = ReplacingMergeTree(`%s`) ", versionColName))

	pkeys := tableSchema.PrimaryKeyColumns
	if len(pkeys) > 0 {
		pkeyStr := strings.Join(pkeys, ",")

		stmtBuilder.WriteString("PRIMARY KEY (")
		stmtBuilder.WriteString(pkeyStr)
		stmtBuilder.WriteString(") ")

		stmtBuilder.WriteString("ORDER BY (")
		stmtBuilder.WriteString(pkeyStr)
		stmtBuilder.WriteString(")")
	}

	return stmtBuilder.String(), nil
}

func (c *ClickhouseConnector) NormalizeRecords(ctx context.Context, req *model.NormalizeRecordsRequest) (*model.NormalizeResponse, error) {
	normBatchID, err := c.GetLastNormalizeBatchID(ctx, req.FlowJobName)
	if err != nil {
		c.logger.Error("[clickhouse] error while getting last sync and normalize batch id", "error", err)
		return nil, err
	}

	// normalize has caught up with sync, chill until more records are loaded.
	if normBatchID >= req.SyncBatchID {
		return &model.NormalizeResponse{
			Done:         false,
			StartBatchID: normBatchID,
			EndBatchID:   req.SyncBatchID,
		}, nil
	}

	destinationTableNames, err := c.getDistinctTableNamesInBatch(
		ctx,
		req.FlowJobName,
		req.SyncBatchID,
		normBatchID,
	)
	if err != nil {
		c.logger.Error("[clickhouse] error while getting distinct table names in batch", "error", err)
		return nil, err
	}

	rawTbl := c.getRawTableName(req.FlowJobName)

	// model the raw table data as inserts.
	for _, tbl := range destinationTableNames {
		// SELECT projection FROM raw_table WHERE _peerdb_batch_id > normalize_batch_id AND _peerdb_batch_id <= sync_batch_id
		selectQuery := strings.Builder{}
		selectQuery.WriteString("SELECT ")

		colSelector := strings.Builder{}
		colSelector.WriteString("(")

		schema := req.TableNameSchemaMapping[tbl]

		projection := strings.Builder{}

		for _, column := range schema.Columns {
			cn := column.Name
			ct := column.Type

			colSelector.WriteString(fmt.Sprintf("`%s`,", cn))
			colType := qvalue.QValueKind(ct)
			clickhouseType, err := qValueKindToClickhouseType(colType)
			if err != nil {
				return nil, fmt.Errorf("error while converting column type to clickhouse type: %w", err)
			}

			switch clickhouseType {
			case "Date":
				projection.WriteString(fmt.Sprintf(
					"toDate(parseDateTime64BestEffortOrNull(JSONExtractString(_peerdb_data, '%s'))) AS `%s`,",
					cn,
					cn,
				))
			case "DateTime64(6)":
				projection.WriteString(fmt.Sprintf(
					"parseDateTime64BestEffortOrNull(JSONExtractString(_peerdb_data, '%s')) AS `%s`,",
					cn,
					cn,
				))
			default:
				projection.WriteString(fmt.Sprintf("JSONExtract(_peerdb_data, '%s', '%s') AS `%s`,", cn, clickhouseType, cn))
			}
		}

		// add _peerdb_sign as _peerdb_record_type / 2
		projection.WriteString(fmt.Sprintf("intDiv(_peerdb_record_type, 2) AS `%s`,", signColName))
		colSelector.WriteString(fmt.Sprintf("`%s`,", signColName))

		// add _peerdb_timestamp as _peerdb_version
		projection.WriteString(fmt.Sprintf("_peerdb_timestamp AS `%s`", versionColName))
		colSelector.WriteString(versionColName)
		colSelector.WriteString(") ")

		selectQuery.WriteString(projection.String())
		selectQuery.WriteString(" FROM ")
		selectQuery.WriteString(rawTbl)
		selectQuery.WriteString(" WHERE _peerdb_batch_id > ")
		selectQuery.WriteString(strconv.FormatInt(normBatchID, 10))
		selectQuery.WriteString(" AND _peerdb_batch_id <= ")
		selectQuery.WriteString(strconv.FormatInt(req.SyncBatchID, 10))
		selectQuery.WriteString(" AND _peerdb_destination_table_name = '")
		selectQuery.WriteString(tbl)
		selectQuery.WriteString("'")

		selectQuery.WriteString(" ORDER BY _peerdb_timestamp")

		insertIntoSelectQuery := strings.Builder{}
		insertIntoSelectQuery.WriteString("INSERT INTO ")
		insertIntoSelectQuery.WriteString(tbl)
		insertIntoSelectQuery.WriteString(colSelector.String())
		insertIntoSelectQuery.WriteString(selectQuery.String())

		q := insertIntoSelectQuery.String()
		c.logger.Info("[clickhouse] insert into select query " + q)

		_, err = c.database.ExecContext(ctx, q)
		if err != nil {
			return nil, fmt.Errorf("error while inserting into normalized table: %w", err)
		}
	}

	endNormalizeBatchId := normBatchID + 1
	err = c.pgMetadata.UpdateNormalizeBatchID(ctx, req.FlowJobName, endNormalizeBatchId)
	if err != nil {
		c.logger.Error("[clickhouse] error while updating normalize batch id", "error", err)
		return nil, err
	}

	return &model.NormalizeResponse{
		Done:         true,
		StartBatchID: endNormalizeBatchId,
		EndBatchID:   req.SyncBatchID,
	}, nil
}

func (c *ClickhouseConnector) getDistinctTableNamesInBatch(
	ctx context.Context,
	flowJobName string,
	syncBatchID int64,
	normalizeBatchID int64,
) ([]string, error) {
	rawTbl := c.getRawTableName(flowJobName)

	//nolint:gosec
	q := fmt.Sprintf(
		`SELECT DISTINCT _peerdb_destination_table_name FROM %s WHERE _peerdb_batch_id > %d AND _peerdb_batch_id <= %d`,
		rawTbl, normalizeBatchID, syncBatchID)

	rows, err := c.database.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error while querying raw table for distinct table names in batch: %w", err)
	}
	defer rows.Close()
	var tableNames []string
	for rows.Next() {
		var tableName sql.NullString
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, fmt.Errorf("error while scanning table name: %w", err)
		}

		if !tableName.Valid {
			return nil, errors.New("table name is not valid")
		}

		tableNames = append(tableNames, tableName.String)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	return tableNames, nil
}

func (c *ClickhouseConnector) GetLastNormalizeBatchID(ctx context.Context, flowJobName string) (int64, error) {
	normalizeBatchID, err := c.pgMetadata.GetLastNormalizeBatchID(ctx, flowJobName)
	if err != nil {
		return 0, fmt.Errorf("error while getting last normalize batch id: %w", err)
	}

	return normalizeBatchID, nil
}
