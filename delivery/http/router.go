package http

import (
	"github.com/gofiber/fiber/v2"
	"movies-review-api/delivery/http/film"
	"movies-review-api/delivery/http/user"

	filmU "movies-review-api/application/film"
	userU "movies-review-api/application/user"
)

func setupRouter(app *fiber.App, config Config) {
	app.Post("/api/v1/ping", ping)

	// route group
	v1 := app.Group("/api/v1/")
	authRouter := v1.Group("/auth")
	userRouter := v1.Group("/user")
	filmRouter := v1.Group("/films")

	userUseCase := userU.New(config.UserRepo)
	user.New(userRouter, userUseCase, config.UserRepo, authRouter)

	filmUsecase := filmU.New(config.FilmRepo)
	film.New(filmRouter, config.FilmRepo, config.UserRepo, filmUsecase)
}
