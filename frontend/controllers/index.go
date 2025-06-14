package controllers

import (
	"encoding/json"
	config "go-k8seenvsc-shutdown-frontend/config"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Response struct {
	Msg string `json:"msg"`
}

func Index(cxt *fiber.Ctx) error {
	backendURL := config.GetBackendURL()
	zap.L().Info("Making a request to backend", zap.String("backendURL", backendURL))
	request := fiber.Get(backendURL)
	_, data, err := request.Bytes()
	if err != nil {
		zap.L().Error("Error making request", zap.Any("error", err))
		return cxt.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	zap.L().Info("Request successful", zap.String("data", string(data)))

	var res Response
	jErr := json.Unmarshal(data, &res)
	if jErr != nil {
		zap.L().Error("Error unmarshalling response", zap.Any("error", jErr))
		return cxt.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": jErr})
	}

	return cxt.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": res.Msg,
	})
}
