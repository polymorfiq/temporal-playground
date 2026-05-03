package main

import (
	"context"
	"log"
	"os"
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
		ID:                 "spider-website-example",
		TaskQueue:          constants.MainTaskQueue,
		WorkflowRunTimeout: (24 * time.Hour) * 30,
	}, workflows.SpiderWebsite, os.Args[1])

	if err != nil {
		log.Fatal("Failed to start spider workflow", err)
	}
}
