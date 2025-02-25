// This file corresponds to xmin based replication.
package peerflow

import (
	"fmt"
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/PeerDB-io/peer-flow/generated/protos"
	"github.com/PeerDB-io/peer-flow/model"
	"github.com/PeerDB-io/peer-flow/shared"
)

func XminFlowWorkflow(
	ctx workflow.Context,
	config *protos.QRepConfig,
	state *protos.QRepFlowState,
) error {
	originalRunID := workflow.GetInfo(ctx).OriginalRunID
	ctx = workflow.WithValue(ctx, shared.FlowNameKey, config.FlowJobName)
	// Support a Query for the current state of the xmin flow.
	err := setWorkflowQueries(ctx, state)
	if err != nil {
		return err
	}

	q := NewQRepFlowExecution(ctx, config, originalRunID)
	logger := q.logger

	err = q.SetupWatermarkTableOnDestination(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup watermark table: %w", err)
	}

	err = q.SetupMetadataTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup metadata tables: %w", err)
	}
	logger.Info("metadata tables setup for peer flow - ", config.FlowJobName)

	err = q.handleTableCreationForResync(ctx, state)
	if err != nil {
		return err
	}

	var lastPartition int64
	replicateXminPartitionCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 24 * 5 * time.Hour,
		HeartbeatTimeout:    time.Minute,
	})
	err = workflow.ExecuteActivity(
		replicateXminPartitionCtx,
		flowable.ReplicateXminPartition,
		q.config,
		state.LastPartition,
		q.runUUID,
	).Get(ctx, &lastPartition)
	if err != nil {
		return fmt.Errorf("xmin replication failed: %w", err)
	}

	if err := q.consolidatePartitions(ctx); err != nil {
		return err
	}

	if config.InitialCopyOnly {
		logger.Info("initial copy completed for peer flow - ", config.FlowJobName)
		return nil
	}

	err = q.handleTableRenameForResync(ctx, state)
	if err != nil {
		return err
	}

	state.LastPartition = &protos.QRepPartition{
		PartitionId: q.runUUID,
		Range:       &protos.PartitionRange{Range: &protos.PartitionRange_IntRange{IntRange: &protos.IntPartitionRange{Start: lastPartition}}},
	}

	logger.Info("Continuing as new workflow",
		"Last Partition", state.LastPartition,
		"Number of Partitions Processed", state.NumPartitionsProcessed)

	// here, we handle signals after the end of the flow because a new workflow does not inherit the signals
	// and the chance of missing a signal is much higher if the check is before the time consuming parts run
	signalChan := model.FlowSignal.GetSignalChannel(ctx)
	q.receiveAndHandleSignalAsync(signalChan)
	if q.activeSignal == model.PauseSignal {
		startTime := workflow.Now(ctx)
		state.CurrentFlowStatus = protos.FlowStatus_STATUS_PAUSED

		for q.activeSignal == model.PauseSignal {
			logger.Info("mirror has been paused", slog.Any("duration", time.Since(startTime)))
			// only place we block on receive, so signal processing is immediate
			val, ok, _ := signalChan.ReceiveWithTimeout(ctx, 1*time.Minute)
			if ok {
				q.activeSignal = model.FlowSignalHandler(q.activeSignal, val, logger)
			} else if err := ctx.Err(); err != nil {
				return err
			}
		}
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	// Continue the workflow with new state
	return workflow.NewContinueAsNewError(ctx, XminFlowWorkflow, config, state)
}
