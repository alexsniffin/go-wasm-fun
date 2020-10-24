package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexsniffin/website/internal/server/models"
	"github.com/alexsniffin/website/internal/server/server"
	"github.com/alexsniffin/website/pkg/config"
	"github.com/alexsniffin/website/pkg/logger"
)

const (
	configFileName = "server"
	prefix         = "WEBSITE"
)

func main() {
	newCfg := models.Config{}
	err := config.NewConfig(configFileName, prefix, &newCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	newLogger, err := logger.NewLogger(newCfg.Logger, newCfg.Environment)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		p, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		newCfg.HTTPServer.Port = p
	}

	newLogger.Info().Msg("setting up server")
	newServer, err := server.NewServer(newCfg, newLogger)
	if err != nil {
		newLogger.Panic().Err(err).Msg("failed to init server")
	}

	go newServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	stopped := <-stop
	newLogger.Info().Msg(stopped.String() + " signal received")
	newServer.Shutdown(false)

	newLogger.Info().Msg("exiting server")
}
