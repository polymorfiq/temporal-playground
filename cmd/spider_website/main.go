package main

import (
	"context"
	"log"
	"temporal-playground/constants"
	"temporal-playground/workflows/job_sites"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                 "amazon-jobs-example",
		TaskQueue:          constants.MainTaskQueue,
		WorkflowRunTimeout: (24 * time.Hour) * 30,
	}, job_sites.AmazonJobs)

	if err != nil {
		log.Fatal("Failed to start spider workflow", err)
	}
}
