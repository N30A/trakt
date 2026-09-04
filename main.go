package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/N30A/trakt/api"
	"github.com/N30A/trakt/protocols/osmand"
	"github.com/N30A/trakt/repos"
	"github.com/N30A/trakt/server"
)

const timeout = time.Second * 5

var (
	wg       sync.WaitGroup
	exitCode = 0
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	poolCtx, poolCancel := context.WithCancel(context.Background())
	pool, err := connectToDB(poolCtx, cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	deviceRepo := repos.NewDeviceRepo(pool)
	positionRepo := repos.NewPositionRepo(pool)

	servers := []server.Server{
		api.New(deviceRepo, positionRepo),
		osmand.New(deviceRepo, positionRepo),
	}

	errChan := make(chan error, len(servers))

	for _, srv := range servers {
		wg.Add(1)
		go func(server server.Server) {
			defer wg.Done()
			if err := server.Start(); err != nil {
				errChan <- fmt.Errorf("%s: %w", server.Name(), err)
			}
		}(srv)
	}

	select {
	case err := <-errChan:
		log.Printf("startup failed: %v\n", err)
		exitCode = 1
	case <-signals:
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for _, server := range servers {
		if err := server.Stop(ctx); err != nil {
			log.Printf("failed to stop %s: %v\n", server.Name(), err)
		}
	}

	wg.Wait()

	log.Println("closing database connection")
	pool.Close()
	poolCancel()
	os.Exit(exitCode)
}
