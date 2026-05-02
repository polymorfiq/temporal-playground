package workflows

import (
	"temporal-playground/activities"
	"temporal-playground/specs"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SpiderWebsite(ctx workflow.Context, proto string, host string) ([]specs.RobotRules, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 30,
	})

	var robotsContents string
	allActivities := activities.Activities{}
	err := workflow.ExecuteActivity(
		ctx,
		allActivities.RetrieveRobots,
		proto,
		host,
	).Get(ctx, &robotsContents)
	if err != nil {
		return nil, err
	}

	var robotsTxt []specs.RobotRules
	err = workflow.ExecuteActivity(
		ctx,
		allActivities.ParseRobots,
		robotsContents,
	).Get(ctx, &robotsTxt)
	if err != nil {
		return nil, err
	}

	return robotsTxt, nil
}
