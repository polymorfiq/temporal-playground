package activities

import (
	"context"
	"fmt"
)

func (a *Activities) SayHello(ctx context.Context, message string) (string, error) {
	return fmt.Sprintf("Hello %s!", message), nil
}
