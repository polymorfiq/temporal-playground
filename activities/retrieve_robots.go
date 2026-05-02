package activities

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.temporal.io/sdk/temporal"
)

func (a *Activities) RetrieveRobots(ctx context.Context, proto string, host string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s://%s/robots.txt", proto, host))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == 0:
		return "", temporal.NewApplicationErrorWithCause("Network Error", "NETWORK_ERROR", errors.New("Network Error (0): "+resp.Status))

	case resp.StatusCode >= 100 && resp.StatusCode < 200:
		return "", temporal.NewNonRetryableApplicationError("Information response", "INFORMATION_RESP", errors.New("Information (1XX): "+resp.Status))

	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return "", temporal.NewNonRetryableApplicationError("Client error", "CLIENT_ERROR", errors.New("Client Error (4XX): "+resp.Status))

	case resp.StatusCode >= 500:
		return "", temporal.NewNonRetryableApplicationError("Server error", "SERVER_ERROR", errors.New("Server Error (5XX): "+resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
