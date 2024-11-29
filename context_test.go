package learning

import (
	"context"
	"errors"
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

func handleCtxCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("context a cancelled")
		return ctx.Err()
	default:
		return nil
	}
}

func a(ctx context.Context) error {
	err := handleCtxCancellation(ctx)
	if err != nil {
		return err
	}
	fmt.Println("a finished")
	return nil
}

func b(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("b finished")
		return errors.New("b error")
	case <-ctx.Done():
		fmt.Println("context b cancelled")
		return ctx.Err()
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := b(ctx)
	if err != nil {
		cancel()
	}

	err = a(ctx)
	if err != nil {
		cancel()
	}

	t.Log("process finished")
}
