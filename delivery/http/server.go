package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Config struct {
}

func RunHttpServer(config Config) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	//app.Use(csrf.New(
	//	csrf.Config{
	//		CookieHTTPOnly: true,
	//		CookieSecure:   true,
	//	}))

	// setup routes
	setupRouter(app, config)

	return app
}
