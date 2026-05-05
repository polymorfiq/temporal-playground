package workflows

import (
	"temporal-playground/activities"
	"temporal-playground/data"
	"time"

	"go.temporal.io/sdk/workflow"
)

func LongRunning(ctx workflow.Context) (string, error) {
	// Set a maximum time limit for the workflow
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24 * 7,
	})

	// Send the message to Kafka
	execInfo := workflow.GetInfo(ctx).WorkflowExecution
	allActivities := activities.Activities{}
	err := workflow.ExecuteActivity(
		ctx,
		allActivities.StartAsyncProcess,
		execInfo.ID,
		execInfo.RunID,
		"Some Message",
	).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	// Sit idle until receiving external Signal
	var asyncResp data.FinishedResponse
	workflow.GetSignalChannel(
		ctx,
		data.FinishedSignal,
	).Receive(ctx, &asyncResp)

	// Use signal response in future Activities
	var result string
	err = workflow.ExecuteActivity(
		ctx,
		allActivities.ContinueAsyncProcess,
		asyncResp,
	).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
