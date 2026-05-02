package workflows

import (
	"temporal-playground/activities"
	"time"

	"go.temporal.io/sdk/workflow"
)

func LongRunning(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 30,
	})

	var result string
	allActivities := activities.Activities{}
	err := workflow.ExecuteActivity(
		ctx,
		allActivities.SayHello,
		"First Message",
	).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	workflow.Sleep(ctx, time.Hour*24)

	err = workflow.ExecuteActivity(
		ctx,
		allActivities.SayHello,
		"Second Message",
	).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	workflow.Sleep(ctx, time.Hour*24)

	err = workflow.ExecuteActivity(
		ctx,
		allActivities.SayHello,
		"Third Message",
	).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
