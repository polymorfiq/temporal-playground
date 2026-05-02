package main

import (
	"log"
	"sync"
	"temporal-playground/activities"
	"temporal-playground/constants"
	"temporal-playground/workflows"

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
	worker.RegisterWorkflow(workflows.FirstWorkflow)
	worker.RegisterWorkflow(workflows.LongRunning)
	worker.RegisterActivity(allActivities)

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
