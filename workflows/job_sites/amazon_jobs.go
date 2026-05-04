package job_sites

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"temporal-playground/activities"
	"temporal-playground/specs"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ericchiang/css"
	"github.com/go-shiori/dom"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/net/html"
)

const BASE_DOMAIN = "https://amazon.jobs"
const BASE_URL = "https://amazon.jobs/content/en/teams/devices-services/leo/product-operations"

func AmazonJobs(ctx workflow.Context) ([]AmazonJob, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 30,
	})

	urlData, err := url.Parse(BASE_URL)
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

	err = workflow.ExecuteActivity(
		ctx,
		allActivities.CheckRobotsAllow,
		robotsTxt,
		BASE_URL,
	).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	jobActivities := JobSiteActivities{}

	var amazonJobs []AmazonJob
	err = workflow.ExecuteActivity(
		ctx,
		jobActivities.RetrieveAmazonJobs,
	).Get(ctx, &amazonJobs)
	if err != nil {
		return nil, err
	}

	return amazonJobs, nil
}

func (a *JobSiteActivities) RetrieveAmazonJobs(ctx context.Context) ([]AmazonJob, error) {
	chromeCtx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var respBody string
	resp, err := chromedp.RunResponse(chromeCtx,
		chromedp.Navigate(BASE_URL),
		chromedp.WaitReady(BASE_URL),
		chromedp.Sleep(time.Second*3),
		chromedp.OuterHTML("html", &respBody),
	)
	if err != nil {
		return nil, err
	}

	if resp.Status != http.StatusOK {
		return nil, temporal.NewApplicationErrorWithCause("Network Error", "NETWORK_ERROR", errors.New("Network Error (0): "+resp.StatusText))
	}

	fmt.Println("BODY", respBody)
	sel, err := css.Parse("a")
	if err != nil {
		panic(err)
	}

	bodyStr := respBody
	node, err := html.Parse(strings.NewReader(bodyStr))
	if err != nil {
		panic(err)
	}

	var amazonJobs []AmazonJob
	for _, ele := range sel.Select(node) {
		var linkHref string
		var linkText string
		hrefAttr := dom.GetAttribute(ele, "href")
		if strings.HasPrefix(hrefAttr, "/jobs/") {
			linkHref = BASE_DOMAIN + hrefAttr
			linkText = dom.InnerText(ele)
		}

		if linkHref != "" {
			amazonJobs = append(amazonJobs, AmazonJob{Url: linkHref, Name: linkText})
		}
	}

	fmt.Printf("JOBS: %s\n", amazonJobs)
	return amazonJobs, nil
}

type AmazonJob struct {
	Url  string
	Name string
}
