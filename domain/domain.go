package domain

import (
	"github.com/gofiber/fiber/v2"
)

type Ping struct {
	Error bool   `json:"error" bson:"error"`
	Msg   string `json:"msg" bson:"msg"`
}

type PingRequest struct {
	Message string `json:"message"`
}

func HandleError(c *fiber.Ctx, err error) error {
	return c.Status(400).JSON(
		fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
}
