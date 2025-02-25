package connbigquery

import (
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"

	"github.com/PeerDB-io/peer-flow/connectors/utils"
	"github.com/PeerDB-io/peer-flow/generated/protos"
	"github.com/PeerDB-io/peer-flow/model/qvalue"
)

type mergeStmtGenerator struct {
	// dataset + raw table
	rawDatasetTable datasetTable
	// destination table name, used to retrieve records from raw table
	dstTableName string
	// dataset + destination table
	dstDatasetTable datasetTable
	// last synced batchID.
	syncBatchID int64
	// last normalized batchID.
	normalizeBatchID int64
	// the schema of the table to merge into
	normalizedTableSchema *protos.TableSchema
	// _PEERDB_IS_DELETED and _SYNCED_AT columns
	peerdbCols *protos.PeerDBColumns
	// map for shorter columns
	shortColumn map[string]string
}

// generateFlattenedCTE generates a flattened CTE.
func (m *mergeStmtGenerator) generateFlattenedCTE() string {
	// for each column in the normalized table, generate CAST + JSON_EXTRACT_SCALAR
	// statement.
	flattenedProjs := make([]string, 0, len(m.normalizedTableSchema.Columns)+3)

	for _, column := range m.normalizedTableSchema.Columns {
		colType := column.Type
		bqType := qValueKindToBigQueryType(colType)
		// CAST doesn't work for FLOAT, so rewrite it to FLOAT64.
		if bqType == bigquery.FloatFieldType {
			bqType = "FLOAT64"
		}
		var castStmt string
		shortCol := m.shortColumn[column.Name]
		switch qvalue.QValueKind(colType) {
		case qvalue.QValueKindJSON, qvalue.QValueKindHStore:
			// if the type is JSON, then just extract JSON
			castStmt = fmt.Sprintf("CAST(PARSE_JSON(JSON_VALUE(_peerdb_data, '$.%s'),wide_number_mode=>'round') AS %s) AS `%s`",
				column.Name, bqType, shortCol)
		// expecting data in BASE64 format
		case qvalue.QValueKindBytes, qvalue.QValueKindBit:
			castStmt = fmt.Sprintf("FROM_BASE64(JSON_VALUE(_peerdb_data,'$.%s')) AS `%s`",
				column.Name, shortCol)
		case qvalue.QValueKindArrayFloat32, qvalue.QValueKindArrayFloat64, qvalue.QValueKindArrayInt16,
			qvalue.QValueKindArrayInt32, qvalue.QValueKindArrayInt64, qvalue.QValueKindArrayString,
			qvalue.QValueKindArrayBoolean, qvalue.QValueKindArrayTimestamp, qvalue.QValueKindArrayTimestampTZ,
			qvalue.QValueKindArrayDate:
			castStmt = fmt.Sprintf("ARRAY(SELECT CAST(element AS %s) FROM "+
				"UNNEST(CAST(JSON_VALUE_ARRAY(_peerdb_data, '$.%s') AS ARRAY<STRING>)) AS element WHERE element IS NOT null) AS `%s`",
				bqType, column.Name, shortCol)
		case qvalue.QValueKindGeography, qvalue.QValueKindGeometry, qvalue.QValueKindPoint:
			castStmt = fmt.Sprintf("CAST(ST_GEOGFROMTEXT(JSON_VALUE(_peerdb_data, '$.%s')) AS %s) AS `%s`",
				column.Name, bqType, shortCol)
		// MAKE_INTERVAL(years INT64, months INT64, days INT64, hours INT64, minutes INT64, seconds INT64)
		// Expecting interval to be in the format of {"Microseconds":2000000,"Days":0,"Months":0,"Valid":true}
		// json.Marshal in SyncRecords for Postgres already does this - once new data-stores are added,
		// this needs to be handled again
		// TODO add interval types again
		// case model.ColumnTypeInterval:
		// castStmt = fmt.Sprintf("MAKE_INTERVAL(0,CAST(JSON_EXTRACT_SCALAR(_peerdb_data, '$.%s.Months') AS INT64),"+
		// 	"CAST(JSON_EXTRACT_SCALAR(_peerdb_data, '$.%s.Days') AS INT64),0,0,"+
		// 	"CAST(CAST(JSON_EXTRACT_SCALAR(_peerdb_data, '$.%s.Microseconds') AS INT64)/1000000 AS  INT64)) AS %s",
		// 	column.Name, column.Name, column.Name, column.Name)
		// TODO add proper granularity for time types, then restore this
		// case model.ColumnTypeTime:
		// 	castStmt = fmt.Sprintf("time(timestamp_micros(CAST(JSON_EXTRACT(_peerdb_data, '$.%s.Microseconds')"+
		// 		" AS int64))) AS %s",
		// 		column.Name, column.Name)
		default:
			castStmt = fmt.Sprintf("CAST(JSON_VALUE(_peerdb_data, '$.%s') AS %s) AS `%s`",
				column.Name, bqType, shortCol)
		}
		flattenedProjs = append(flattenedProjs, castStmt)
	}
	flattenedProjs = append(
		flattenedProjs,
		"_peerdb_timestamp",
		"_peerdb_record_type AS _rt",
		"_peerdb_unchanged_toast_columns AS _ut",
	)

	// normalize anything between last normalized batch id to last sync batchid
	return fmt.Sprintf("WITH _f AS "+
		"(SELECT %s FROM `%s` WHERE _peerdb_batch_id>%d AND _peerdb_batch_id<=%d AND "+
		"_peerdb_destination_table_name='%s')",
		strings.Join(flattenedProjs, ","), m.rawDatasetTable.string(), m.normalizeBatchID,
		m.syncBatchID, m.dstTableName)
}

// generateDeDupedCTE generates a de-duped CTE.
func (m *mergeStmtGenerator) generateDeDupedCTE() string {
	const cte = `_dd AS (
		SELECT _peerdb_ranked.* FROM(
				SELECT RANK() OVER(
					PARTITION BY %s ORDER BY _peerdb_timestamp DESC
				) AS _peerdb_rank,* FROM _f
			) _peerdb_ranked
			WHERE _peerdb_rank=1
	) SELECT * FROM _dd`

	shortPkeys := make([]string, 0, len(m.normalizedTableSchema.PrimaryKeyColumns))
	for _, pkeyCol := range m.normalizedTableSchema.PrimaryKeyColumns {
		shortPkeys = append(shortPkeys, m.shortColumn[pkeyCol])
	}

	pkeyColsStr := fmt.Sprintf("(CONCAT(%s))", strings.Join(shortPkeys,
		", '_peerdb_concat_', "))
	return fmt.Sprintf(cte, pkeyColsStr)
}

// generateMergeStmt generates a merge statement.
func (m *mergeStmtGenerator) generateMergeStmt(unchangedToastColumns []string) string {
	// comma separated list of column names
	columnCount := len(m.normalizedTableSchema.Columns)
	backtickColNames := make([]string, 0, columnCount)
	shortBacktickColNames := make([]string, 0, columnCount)
	pureColNames := make([]string, 0, columnCount)
	for i, col := range m.normalizedTableSchema.Columns {
		shortCol := fmt.Sprintf("_c%d", i)
		m.shortColumn[col.Name] = shortCol
		backtickColNames = append(backtickColNames, fmt.Sprintf("`%s`", col.Name))
		shortBacktickColNames = append(shortBacktickColNames, fmt.Sprintf("`%s`", shortCol))
		pureColNames = append(pureColNames, col.Name)
	}
	csep := strings.Join(backtickColNames, ", ")
	shortCsep := strings.Join(shortBacktickColNames, ", ")
	insertColumnsSQL := csep + fmt.Sprintf(", `%s`", m.peerdbCols.SyncedAtColName)
	insertValuesSQL := shortCsep + ",CURRENT_TIMESTAMP"

	updateStatementsforToastCols := m.generateUpdateStatements(pureColNames, unchangedToastColumns)
	if m.peerdbCols.SoftDelete {
		softDeleteInsertColumnsSQL := insertColumnsSQL + fmt.Sprintf(",`%s`", m.peerdbCols.SoftDeleteColName)
		softDeleteInsertValuesSQL := insertValuesSQL + ",TRUE"

		updateStatementsforToastCols = append(updateStatementsforToastCols,
			fmt.Sprintf("WHEN NOT MATCHED AND _d._rt=2 THEN INSERT (%s) VALUES(%s)",
				softDeleteInsertColumnsSQL, softDeleteInsertValuesSQL))
	}
	updateStringToastCols := strings.Join(updateStatementsforToastCols, " ")

	pkeySelectSQLArray := make([]string, 0, len(m.normalizedTableSchema.PrimaryKeyColumns))
	for _, pkeyColName := range m.normalizedTableSchema.PrimaryKeyColumns {
		pkeySelectSQLArray = append(pkeySelectSQLArray, fmt.Sprintf("_t.%s=_d.%s",
			pkeyColName, m.shortColumn[pkeyColName]))
	}
	// t.<pkey1> = d.<pkey1> AND t.<pkey2> = d.<pkey2> ...
	pkeySelectSQL := strings.Join(pkeySelectSQLArray, " AND ")

	deletePart := "DELETE"
	if m.peerdbCols.SoftDelete {
		colName := m.peerdbCols.SoftDeleteColName
		deletePart = fmt.Sprintf("UPDATE SET %s=TRUE", colName)
		if m.peerdbCols.SyncedAtColName != "" {
			deletePart = fmt.Sprintf("%s,%s=CURRENT_TIMESTAMP",
				deletePart, m.peerdbCols.SyncedAtColName)
		}
	}

	return fmt.Sprintf("MERGE `%s` _t USING(%s,%s) _d"+
		" ON %s WHEN NOT MATCHED AND _d._rt!=2 THEN "+
		"INSERT (%s) VALUES(%s) "+
		"%s WHEN MATCHED AND _d._rt=2 THEN %s;",
		m.dstDatasetTable.table, m.generateFlattenedCTE(), m.generateDeDupedCTE(),
		pkeySelectSQL, insertColumnsSQL, insertValuesSQL, updateStringToastCols, deletePart)
}

/*
This function takes an array of unique unchanged toast column groups and an array of all column names,
and returns suitable UPDATE statements as part of a MERGE operation.

Algorithm:
1. Iterate over each unique unchanged toast column group.
2. Split the group into individual column names.
3. Calculate the other columns by finding the set difference between all column names
and the unchanged columns.
4. Generate an update statement for the current group by setting the appropriate conditions
and updating the other columns (not the unchanged toast columns)
5. Append the update statement to the list of generated statements.
6. Repeat steps 1-5 for each unique unchanged toast column group.
7. Return the list of generated update statements.
*/
func (m *mergeStmtGenerator) generateUpdateStatements(allCols []string, unchangedToastColumns []string) []string {
	handleSoftDelete := m.peerdbCols.SoftDelete && (m.peerdbCols.SoftDeleteColName != "")
	// weird way of doing it but avoids prealloc lint
	updateStmts := make([]string, 0, func() int {
		if handleSoftDelete {
			return 2 * len(unchangedToastColumns)
		}
		return len(unchangedToastColumns)
	}())

	for _, cols := range unchangedToastColumns {
		unchangedColsArray := strings.Split(cols, ",")
		otherCols := utils.ArrayMinus(allCols, unchangedColsArray)
		tmpArray := make([]string, 0, len(otherCols))
		for _, colName := range otherCols {
			tmpArray = append(tmpArray, fmt.Sprintf("`%s`=_d.%s", colName, m.shortColumn[colName]))
		}

		// set the synced at column to the current timestamp
		if m.peerdbCols.SyncedAtColName != "" {
			tmpArray = append(tmpArray, fmt.Sprintf("`%s`=CURRENT_TIMESTAMP",
				m.peerdbCols.SyncedAtColName))
		}
		// set soft-deleted to false, tackles insert after soft-delete
		if handleSoftDelete {
			tmpArray = append(tmpArray, fmt.Sprintf("`%s`=FALSE",
				m.peerdbCols.SoftDeleteColName))
		}

		ssep := strings.Join(tmpArray, ",")
		updateStmt := fmt.Sprintf(`WHEN MATCHED AND
		_rt!=2 AND _ut='%s'
		THEN UPDATE SET %s`, cols, ssep)
		updateStmts = append(updateStmts, updateStmt)

		// generates update statements for the case where updates and deletes happen in the same branch
		// the backfill has happened from the pull side already, so treat the DeleteRecord as an update
		// and then set soft-delete to true.
		if handleSoftDelete {
			tmpArray = append(tmpArray[:len(tmpArray)-1],
				fmt.Sprintf("`%s`=TRUE", m.peerdbCols.SoftDeleteColName))
			ssep := strings.Join(tmpArray, ",")
			updateStmt := fmt.Sprintf(`WHEN MATCHED AND
			_rt=2 AND _ut='%s'
			THEN UPDATE SET %s`, cols, ssep)
			updateStmts = append(updateStmts, updateStmt)
		}
	}
	return updateStmts
}
