package comment

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"movies-review-api/delivery/http/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"movies-review-api/domain"
	"movies-review-api/pkg/logger"
)

var (
	validate = validator.New()
)

type CommentHandler struct {
	CommentRepo    domain.CommentRepository
	CommentUsecase domain.CommentUsecase
	Logger         *zap.Logger
}

func New(commentRouter fiber.Router, r domain.CommentRepository, userRepo domain.UserRepository, commentUsecase domain.CommentUsecase) {
	handler := &CommentHandler{
		CommentRepo:    r,
		CommentUsecase: commentUsecase,
	}

	l, _ := logger.InitLogger()

	handler.Logger = l

	commentRouter.Post("/", middleware.Protected(userRepo), handler.AddComment)
	commentRouter.Get("/:filmId", middleware.Protected(userRepo), handler.FetchPostComments)
}

func (h *CommentHandler) AddComment(c *fiber.Ctx) error {
	var data domain.NewCommentRequest

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

	data.UserId = c.Locals("user_id").(string)

	comment, err := h.CommentUsecase.AddComment(context.TODO(), &data)

	if err != nil {
		return domain.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  comment,
	})
}

func (h *CommentHandler) FetchPostComments(c *fiber.Ctx) error {

	filmId := c.Params("filmId")

	page := c.Query("page", "1")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return domain.HandleError(c, err)
	}

	limit := c.Query("limit", "20")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return domain.HandleError(c, err)
	}

	data, err := h.CommentRepo.FetchPaginatedFilmComments(context.TODO(), filmId, int64(pageInt), int64(limitInt))

	if err != nil {
		return domain.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  data,
	})
}
