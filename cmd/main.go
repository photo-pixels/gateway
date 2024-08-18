package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/photo-pixels/gateway/internal/app"
	"github.com/photo-pixels/gateway/internal/server"
	"github.com/photo-pixels/platform/config"
)

const defaultPort = "8083"

func main() {
	var args config.Arguments
	if _, err := flags.Parse(&args); err != nil {
		panic(err)
	}

	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.NewApp(cfgProvider)
	if err = application.Create(ctx); err != nil {
		panic(err)
	}

	logger := application.GetLogger()

	srv := server.NewGraphQLServer(
		application.GetLogger(),
		application.GetServerConfig(),
		application.GetAuth(),
		application.GetServiceClientsContainer(),
	)

	go func() {
		err = srv.Start(ctx)
		if err != nil {
			log.Fatalf("fail start app: %v", err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	logger.Infof("--- stopped application ---")
	srv.Stop()
	logger.Infof("--- stop application ---")
}
