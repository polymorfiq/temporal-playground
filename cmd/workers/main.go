package main

import (
	"log"
	"sync"
	"temporal-playground/activities"
	"temporal-playground/constants"
	"temporal-playground/workflows"
	"temporal-playground/workflows/job_sites"

	"go.temporal.io/sdk/client"
	temporalWorker "go.temporal.io/sdk/worker"
)

func main() {
	// Connect to local Temporal server
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Start worker
	worker := temporalWorker.New(c, constants.MainTaskQueue, temporalWorker.Options{})
	allActivities := &activities.Activities{}
	jobSiteActivities := &job_sites.JobSiteActivities{}
	worker.RegisterWorkflow(workflows.FirstWorkflow)
	worker.RegisterWorkflow(workflows.LongRunning)
	worker.RegisterWorkflow(workflows.SpiderWebsite)
	worker.RegisterWorkflow(job_sites.AmazonJobs)
	worker.RegisterActivity(allActivities)
	worker.RegisterActivity(jobSiteActivities)

	err = worker.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer worker.Stop()

	// Loop forever
	for {
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	}
}
