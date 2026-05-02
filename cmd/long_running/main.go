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
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                 "long-running-workflow-1",
		TaskQueue:          constants.MainTaskQueue,
		WorkflowRunTimeout: (24 * time.Hour) * 30,
	}, workflows.LongRunning)

	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                 "long-running-workflow-2",
		TaskQueue:          constants.MainTaskQueue,
		WorkflowRunTimeout: (24 * time.Hour) * 30,
	}, workflows.LongRunning)

	if err != nil {
		log.Fatal("Failed to start async workflow", err)
	}
}
