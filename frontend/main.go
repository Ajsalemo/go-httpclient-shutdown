package main

import (
	"os"
	"os/signal"
	"syscall"

	"go-k8seenvsc-shutdown-frontend/config"
	"go-k8seenvsc-shutdown-frontend/controllers"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	// Set the backend URL
	config.SetBackendURL()
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
			app.Shutdown()
		case syscall.SIGTERM:
			zap.L().Warn("SIGTERM recieved.., shutting down the application..")
			app.Shutdown()
		case syscall.SIGQUIT:
			zap.L().Warn("SIGQUIT recieved.., shutting down the application..")
			app.Shutdown()
		case syscall.SIGINT:
			zap.L().Warn("SIGINT recieved.., shutting down the application..")
			app.Shutdown()
		}
	}()

	app.Listen(":8080")
}