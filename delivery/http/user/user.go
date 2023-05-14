package user

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"

	"movies-review-api/domain"
	"movies-review-api/pkg/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	validate = validator.New()
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
	Logger      *zap.Logger
}

func New(userRouter fiber.Router, u domain.UserUsecase, auth fiber.Router) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	l, _ := logger.InitLogger()

	handler.Logger = l

	auth.Post("/login", handler.Login)
	userRouter.Post("/signup", handler.SignUp)
}

func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	var data domain.SignupRequest

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return domain.HandleError(c, err)
	}

	if err := validate.Struct(data); err != nil {
		return domain.HandleValidationError(c, err)
	}

	_, err := h.UserUsecase.Signup(context.TODO(), &data)

	if err != nil {
		h.Logger.Error(err.Error(), zap.Error(err))
		return domain.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  nil,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	// get user data
	var data domain.LoginRequest

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
	}

	if err := validate.Struct(data); err != nil {
		return domain.HandleValidationError(c, err)
	}

	data.Email = strings.ToLower(data.Email)

	existingUser, err := h.UserUsecase.Login(context.TODO(), &data)

	if err != nil {
		return domain.HandleError(c, err)
	}

	token, err := domain.GenerateToken(*existingUser)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"error": true,
				"msg":   "Invalid login credentials",
			})
	}

	existingUser.Password = ""

	return c.JSON(fiber.Map{
		"error": false,
		"data": fiber.Map{
			"user":  existingUser,
			"token": token,
		},
	})
}
