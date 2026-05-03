package workflows

import (
	"errors"
	"net/url"
	"temporal-playground/activities"
	"temporal-playground/specs"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func SpiderWebsite(ctx workflow.Context, spiderUrl string) ([]specs.RobotRules, error) {
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

	var relevantRobotsTxt *specs.RobotRules
	for _, robotsTxt := range robotsTxt {
		if robotsTxt.UserAgent == "*" {
			relevantRobotsTxt = &robotsTxt
		}
	}

	allowed := true
	if relevantRobotsTxt != nil {
		for _, disallow := range relevantRobotsTxt.Disallow {
			if checkDisallow(disallow, urlData.Path) {
				allowed = false
			}
		}
	}

	if !allowed {
		return nil, temporal.NewNonRetryableApplicationError(
			"Disallowed by robots.txt",
			"Disallowed",
			errors.New("Disallowed by robots.txt"),
		)
	}

	return robotsTxt, nil
}

func checkDisallow(disallow string, path string) bool {
	return true
}
