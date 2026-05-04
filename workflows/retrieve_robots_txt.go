package workflows

import (
	"net/url"
	"temporal-playground/activities"
	"temporal-playground/specs"
	"time"

	"go.temporal.io/sdk/workflow"
)

func RetrieveRobotsTxt(ctx workflow.Context, spiderUrl string) ([]specs.RobotRules, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 30,
	})

	urlData, err := url.Parse(spiderUrl)
	if err != nil {
		return nil, err
	}

	var robotsContents string
	allActivities := activities.Activities{}
	err = workflow.ExecuteActivity(
		ctx,
		allActivities.RetrieveRobots,
		urlData.Scheme,
		urlData.Host,
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
