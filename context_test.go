package learning

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func anjay(ctx context.Context) {
	fmt.Println("anjay")
}

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("overslept")
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		}
	}()

	anjay(ctx)
}
