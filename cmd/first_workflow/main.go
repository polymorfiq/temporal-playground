package main

import (
	"context"
	"log"
	"temporal-playground/constants"
	"temporal-playground/workflows"
	"time"

	"go.temporal.io/sdk/client"
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

	// Run workflow
	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                 "first-workflow-id",
		TaskQueue:          constants.MainTaskQueue,
		WorkflowRunTimeout: 5 * time.Minute,
	}, workflows.FirstWorkflow, "Test Input")

	if err != nil {
		log.Fatal("Failed to start async workflow", err)
	}
}
