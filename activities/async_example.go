package activities

import (
	"context"
	"fmt"
	"log"
	"temporal-playground/data"
	"time"

	"go.temporal.io/sdk/client"
)

func (a *Activities) StartAsyncProcess(ctx context.Context, workflowId string, runId string, data string) error {
	fmt.Sprintf("Sending data: %s\n", data)
	go sendMessageAsync(workflowId, runId, data)

	return nil
}

func (a *Activities) ContinueAsyncProcess(ctx context.Context, resp data.FinishedResponse) error {
	fmt.Sprintf("Received data: %s\n", resp.Data)
	return nil
}

func sendMessageAsync(workflowId string, runId string, msg string) {
	time.Sleep(30 * time.Second)
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to connect to temporal server", err)
	}

	err = temporalClient.SignalWorkflow(
		context.Background(),
		workflowId,
		runId,
		data.FinishedSignal,
		data.FinishedResponse{Data: msg},
	)
	if err != nil {
		log.Fatalf("Unable to signal workflow: %v", err)
	}
}
