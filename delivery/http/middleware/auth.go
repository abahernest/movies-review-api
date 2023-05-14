package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"movies-review-api/domain"
	"reflect"
	"strings"
)

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": true, "msg": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"error": true, "msg": "Invalid or expired JWT", "data": nil})
	}
}

func Protected(userRepo domain.UserRepository) func(*fiber.Ctx) error {
	return NewJwtHandler(domain.Config{
		SigningKey:        []byte(viper.GetString("JWT_SECRET_KEY")),
		ErrorHandler:      jwtError,
		ValidatorFunction: userRepo,
	})
}

func NewJwtHandler(config ...domain.Config) fiber.Handler {

	// Init config
	var cfg domain.Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed JWT")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
			}
		}
	}
	if cfg.SigningKey == nil && len(cfg.SigningKeys) == 0 {
		panic("Fiber: JWT middleware requires signing key")
	}
	if cfg.SigningMethod == "" {
		cfg.SigningMethod = "HS256"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "token"
	}
	if cfg.Claims == nil {
		cfg.Claims = jwt.MapClaims{}
	}
	if cfg.TokenLookup == "" {
		cfg.TokenLookup = "header:" + fiber.HeaderAuthorization
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = "Bearer"
	}

	cfg.KeyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != cfg.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		if len(cfg.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := cfg.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("unexpected jwt key id=%v", t.Header["kid"])
		}
		return cfg.SigningKey, nil
	}
	// Initialize
	extractors := make([]func(c *fiber.Ctx) (string, error), 0)
	rootParts := strings.Split(cfg.TokenLookup, ",")
	for _, rootPart := range rootParts {
		parts := strings.Split(strings.TrimSpace(rootPart), ":")

		switch parts[0] {
		case "header":
			extractors = append(extractors, jwtFromHeader(parts[1], cfg.AuthScheme))
		case "query":
			extractors = append(extractors, jwtFromQuery(parts[1]))
		case "param":
			extractors = append(extractors, jwtFromParam(parts[1]))
		case "cookie":
			extractors = append(extractors, jwtFromCookie(parts[1]))
		}
	}
	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		var auth string
		var err error

		for _, extractor := range extractors {
			auth, err = extractor(c)
			if auth != "" && err == nil {
				break
			}
		}

		if err != nil {
			return cfg.ErrorHandler(c, err)
		}
		token := new(jwt.Token)

		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.KeyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.KeyFunc)
		}
		if err == nil && token.Valid {
			// Store user information from token into context.
			c.Locals(cfg.ContextKey, token)

			// parse the claims type

			if token.Claims.(jwt.MapClaims)["user"] != nil {
				if token.Claims.(jwt.MapClaims)["user"].(map[string]interface{})["_id"] != nil {
					userId := token.Claims.(jwt.MapClaims)["user"].(map[string]interface{})["_id"].(string)
					_, err := cfg.ValidatorFunction.GetById(context.Background(), userId)
					if err != nil {
						return cfg.ErrorHandler(c, err)
					}
					c.Locals("user_id", userId)
				}

				if token.Claims.(jwt.MapClaims)["user"].(map[string]interface{})["email"] != nil {
					c.Locals("email", token.Claims.(jwt.MapClaims)["user"].(map[string]interface{})["email"])
				}
			}

			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}

// jwtFromHeader returns a function that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or malformed JWT")
	}
}

// jwtFromQuery returns a function that extracts token from the query string.
func jwtFromQuery(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", errors.New("missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromParam returns a function that extracts token from the url param string.
func jwtFromParam(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", errors.New("missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromCookie returns a function that extracts token from the named cookie.
func jwtFromCookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Cookies(name)
		if token == "" {
			return "", errors.New("missing or malformed JWT")
		}
		return token, nil
	}
}
