package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"movies-review-api/domain"
)

type Config struct {
	UserRepo domain.UserRepository
	FilmRepo domain.FilmRepository
}

func RunHttpServer(config Config) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	// setup routes
	setupRouter(app, config)

	return app
}
