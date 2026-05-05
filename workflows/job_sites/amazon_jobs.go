package job_sites

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"temporal-playground/activities"
	"temporal-playground/networking"
	"temporal-playground/persistence"
	"temporal-playground/specs"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const BASE_DOMAIN = "https://amazon.jobs"
const BASE_URL = "https://amazon.jobs/api/jobs/search?is_als=true"

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

	err = workflow.ExecuteActivity(
		ctx,
		jobActivities.SaveAmazonJobs,
		amazonJobs,
	).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return amazonJobs, nil
}

func (a *JobSiteActivities) RetrieveAmazonJobs(ctx context.Context) ([]AmazonJob, error) {
	search, err := json.Marshal(AmazonJobSearch{
		AccessLevel:         "EXTERNAL",
		ContentFilterFacets: []AmazonContentFilterFacet{{"primarySearchLabel", 9999}},
		ExcludeFacets: []AmazonExcludeFacet{
			{"isConfidential", []AmazonFacetValue{{"1"}}},
			{"businessCategory", []AmazonFacetValue{{"a-confidential-job"}}},
		},
		FilterFacets: []AmazonFilterFacet{{"optionalSearchLabels", 9999, []AmazonFacetValue{{"devices-services.kuiper-prod-ops"}}}},
		Size:         200,
		Start:        0,
		Sort:         AmazonSort{SortType: "SCORE", SortOrder: "DESCENDING"},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(networking.ProxifiedUrl(BASE_URL), "application/json", bytes.NewReader(search))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, temporal.NewApplicationErrorWithCause("Network Error", "NETWORK_ERROR", errors.New("Network Error (0): "+resp.Status))
	}

	bodyStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := AmazonJobSearchResult{}
	json.Unmarshal(bodyStr, &results)

	jobs := []AmazonJob{}
	for _, hit := range results.SearchHits {
		updatedAt, err := strconv.Atoi(hit.Fields.UpdatedDate[0])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error converting updated date (%s): %v", hit.Fields.UrlNextStep[0], err))
		}

		createdAt, err := strconv.Atoi(hit.Fields.CreatedDate[0])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error converting created date (%s): %v", hit.Fields.UrlNextStep[0], err))
		}

		jobs = append(jobs, AmazonJob{
			Country:                 hit.Fields.Country[0],
			PreferredQualifications: hit.Fields.PreferredQualifications[0],
			JobCode:                 hit.Fields.JobCode[0],
			Description:             hit.Fields.Description[0],
			BasicQualifications:     hit.Fields.BasicQualifications[0],
			ShortDescription:        hit.Fields.ShortDescription[0],
			UpdatedDate:             time.Unix(int64(updatedAt), 0),
			Title:                   hit.Fields.Title[0],
			NormalizedLocation:      hit.Fields.NormalizedLocation[0],
			CreatedDate:             time.Unix(int64(createdAt), 0),
			TeamCategory:            hit.Fields.TeamCategory[0],
			UrlNextStep:             hit.Fields.UrlNextStep[0],
			BusinessCategory:        hit.Fields.BusinessCategory[0],
			IsManager:               hit.Fields.IsManager[0] == 1,
			Location:                hit.Fields.Location[0],
			Category:                hit.Fields.Category[0],
			IcimsJobId:              hit.Fields.IcimsJobId[0],
		})
	}

	return jobs, nil
}

func (a *JobSiteActivities) SaveAmazonJobs(ctx context.Context, amazonJobs []AmazonJob) error {
	for _, job := range amazonJobs {
		err := persistence.SaveJob("amazon_jobs", job.IcimsJobId, job)
		if err != nil {
			return err
		}
	}

	return nil
}

type AmazonJob struct {
	Country                 string
	PreferredQualifications string
	JobCode                 string
	Description             string
	BasicQualifications     string
	ShortDescription        string
	UpdatedDate             time.Time
	Title                   string
	NormalizedLocation      string
	CreatedDate             time.Time
	TeamCategory            string
	UrlNextStep             string
	BusinessCategory        string
	IsManager               bool
	Location                string
	Category                string
	IcimsJobId              string
}

type AmazonJobSearch struct {
	AccessLevel         string                     `json:"accessLevel"`
	ContentFilterFacets []AmazonContentFilterFacet `json:"contentFilterFacets"`
	ExcludeFacets       []AmazonExcludeFacet       `json:"excludeFacets"`
	FilterFacets        []AmazonFilterFacet        `json:"filterFacets"`
	Query               string                     `json:"query"`
	Size                int                        `json:"size"`
	Start               int                        `json:"start"`
	Sort                AmazonSort                 `json:"sort"`
}

type AmazonContentFilterFacet struct {
	Name                string `json:"name"`
	RequestedFacetCount int    `json:"requestedFacetCount"`
}

type AmazonFilterFacet struct {
	Name                string             `json:"name"`
	RequestedFacetCount int                `json:"requestedFacetCount"`
	Values              []AmazonFacetValue `json:"values"`
}

type AmazonExcludeFacet struct {
	Name   string             `json:"name"`
	Values []AmazonFacetValue `json:"values"`
}

type AmazonFacetValue struct {
	Name string `json:"name"`
}

type AmazonSort struct {
	SortOrder string `json:"sortOrder"`
	SortType  string `json:"sortType"`
}

type AmazonJobSearchResult struct {
	Found      int               `json:"found"`
	SearchHits []AmazonSearchHit `json:"searchHits"`
}

type AmazonSearchHit struct {
	Fields AmazonSearchHitField `json:"fields"`
}

type AmazonSearchHitField struct {
	Country                 []string `json:"country"`
	PreferredQualifications []string `json:"preferredQualifications"`
	JobCode                 []string `json:"jobCode"`
	Description             []string `json:"description"`
	BasicQualifications     []string `json:"basicQualifications"`
	ShortDescription        []string `json:"shortDescription"`
	UpdatedDate             []string `json:"updatedDate"`
	Title                   []string `json:"title"`
	NormalizedLocation      []string `json:"normalizedLocation"`
	CreatedDate             []string `json:"createdDate"`
	TeamCategory            []string `json:"teamCategory"`
	UrlNextStep             []string `json:"urlNextStep"`
	BusinessCategory        []string `json:"businessCategory"`
	IsManager               []int    `json:"isManager"`
	Location                []string `json:"location"`
	Category                []string `json:"category"`
	IcimsJobId              []string `json:"icimsJobId"`
}
