package http

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"movies-review-api/domain"
)

func ping(c *fiber.Ctx) error {
	var requestBody domain.PingRequest

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return domain.HandleError(c, err)
	}

	fmt.Println(requestBody)

	return c.JSON(domain.Ping{
		Error: false,
		Msg:   "pong",
	})
}
