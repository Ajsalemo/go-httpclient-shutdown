package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Index(cxt *fiber.Ctx) error {
	time := time.Now().Format(time.ANSIC)
	zap.L().Info("Received an incoming request at index endpoint")
	return cxt.JSON(fiber.Map{"msg": "The current time is " + time})
}