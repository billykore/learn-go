package learning

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscall and doing clean up operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscall that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

type db struct{}

func (db *db) shutdown(ctx context.Context) error {
	return nil
}

func initDB() (*db, error) {
	log.Println("initDB")
	return &db{}, nil
}

type usecase struct{}

func (uc *usecase) shutdown(ctx context.Context) error {
	return nil
}

func initUsecase() (*usecase, error) {
	log.Println("initUsecase")
	return &usecase{}, nil
}

func TestGraceful(t *testing.T) {
	db, err := initDB()
	assert.NoError(t, err)

	uc, err := initUsecase()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wait := gracefulShutdown(ctx, 5*time.Second, map[string]operation{
		"db": func(ctx context.Context) error {
			return db.shutdown(ctx)
		},
		"usecase": func(ctx context.Context) error {
			return uc.shutdown(ctx)
		},
	})

	<-wait
}
