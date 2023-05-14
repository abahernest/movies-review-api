package http

import (
	"github.com/gofiber/fiber/v2"
	"movies-review-api/delivery/http/user"

	userU "movies-review-api/application/user"
)

func setupRouter(app *fiber.App, config Config) {
	app.Post("/api/v1/ping", ping)

	// route group
	v1 := app.Group("/api/v1/")
	authRouter := v1.Group("/auth")
	userRouter := v1.Group("/user")

	userUseCase := userU.New(config.UserRepo)
	user.New(userRouter, userUseCase, authRouter)
}
