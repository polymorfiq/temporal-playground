package workflows

import (
	"temporal-playground/activities"
	"time"

	"go.temporal.io/sdk/workflow"
)

func FirstWorkflow(ctx workflow.Context, input string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 30,
	})

	var result string
	allActivities := activities.Activities{}
	err := workflow.ExecuteActivity(
		ctx,
		allActivities.SayHello,
		input,
	).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
