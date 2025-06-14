package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-k8seenvsc-shutdown-frontend/config"
	"go-k8seenvsc-shutdown-frontend/controllers"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func makeRequestsOnShutdown() {
	backendURL := config.GetBackendURL()
	zap.L().Info("Making a request to backend", zap.String("backendURL", backendURL))
	request := fiber.Get(backendURL)
	_, data, err := request.Bytes()
	if err != nil {
		zap.L().Error("Error making request", zap.Any("error", err))
	}

	zap.L().Info("Request successful", zap.String("data", string(data)))
}

func initiateShutdownCountdown() {
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()
	start := time.Now()

	shutdownTimeLimit := os.Getenv("SHUTDOWN_TIME_LIMIT")
	if shutdownTimeLimit == "" {
		shutdownTimeLimit = "15"
	}
	// Parse the shutdownTimeLimit (of type string) into a time.Duration
	duration, err := time.ParseDuration(shutdownTimeLimit + "s")
	if err != nil {
		zap.L().Error("Error parsing SHUTDOWN_TIME_LIMIT", zap.Error(err))
		return
	}
	// When a signal is caught, loop over this every 2 seconds up until whatever time SHUTDOWN_TIME_LIMIT  is set to
	// in which at that point, this function should return and the app will actually exit
	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start)
			// If the elapsed time is greater than 100 seconds, return a 500
			if elapsed > duration {
				zap.L().Info("Elapsed time: " + elapsed.String())
				zap.L().Warn("SHUTDOWN_TIME_LIMIT was reached, exiting application..")
				return
			}

			makeRequestsOnShutdown()
		}
	}
}

func main() {
	app := fiber.New()

	app.Get("/", controllers.Index)
	// Notify the application of the below signals to be handled on shutdown
	s := make(chan os.Signal, 1)
	signal.Notify(s,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	// Goroutine to clean up prior to shutting down
	go func() {
		sig := <-s
		switch sig {
		case os.Interrupt:
			zap.L().Warn("CTRL+C / os.Interrupt recieved, shutting down the application..")
			initiateShutdownCountdown()
			app.Shutdown()
		case syscall.SIGTERM:
			zap.L().Warn("SIGTERM recieved.., shutting down the application..")
			initiateShutdownCountdown()
			app.Shutdown()
		case syscall.SIGQUIT:
			zap.L().Warn("SIGQUIT recieved.., shutting down the application..")
			initiateShutdownCountdown()
			app.Shutdown()
		case syscall.SIGINT:
			zap.L().Warn("SIGINT recieved.., shutting down the application..")
			initiateShutdownCountdown()
			app.Shutdown()
		}
	}()

	app.Listen(":8080")
}
