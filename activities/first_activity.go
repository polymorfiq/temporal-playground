package activities

import (
	"context"
	"fmt"
)

func (a *Activities) FirstActivity(ctx context.Context, message string) (string, error) {
	return fmt.Sprintf("1st - Hello %s!", message), nil
}
