package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

var bgCtx = context.Background()

func runRedis(ctx context.Context, waitGroup *errgroup.Group) {
	waitGroup.Go(func() error {
		// Redis Start
		log.Println("Redis Server Start")

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		// Redis.shutdown()
		log.Println("Redis Server Shutdown")

		return nil
	})
}

func runConsumer(ctx context.Context, waitGroup *errgroup.Group) {
	waitGroup.Go(func() error {
		// Redis Start
		log.Println("Consumer Server Start")

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		// Consumer.shutdown()
		log.Println("Consumer Server Shutdown")

		return nil
	})
}

func runAnyServer(ctx context.Context, waitGroup *errgroup.Group) {
	waitGroup.Go(func() error {
		// Redis Start
		log.Println("Any Server Start")

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		// Any Server.shutdown()
		log.Println("Any Server Shutdown")

		return nil
	})
}

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, cancelCtx := signal.NotifyContext(bgCtx, interruptSignals...)

	defer cancelCtx()

	waitGroup, gtx := errgroup.WithContext(ctx)

	runRedis(gtx, waitGroup)
	runConsumer(gtx, waitGroup)
	runAnyServer(gtx, waitGroup)

	go func() error {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Error Listen", err)
		}

		return nil
	}()

	log.Printf("Server Listen on %s", ":8000")

	err := waitGroup.Wait()
	if err != nil {
		log.Fatal("Error wait group", err)
	}

	// App Server exit()
	server.Shutdown(bgCtx)
	log.Println("App Server Exit")

	log.Println("Server Exited Properly")
}
