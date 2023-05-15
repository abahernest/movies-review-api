package film

import (
	"context"
	"movies-review-api/delivery/http/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"movies-review-api/domain"
	"movies-review-api/pkg/logger"
)

type FilmHandler struct {
	FilmRepo    domain.FilmRepository
	FilmUsecase domain.FilmUsecase
	Logger      *zap.Logger
}

func New(filmRouter fiber.Router, r domain.FilmRepository, userRepo domain.UserRepository, filmUsecase domain.FilmUsecase) {
	handler := &FilmHandler{
		FilmRepo:    r,
		FilmUsecase: filmUsecase,
	}

	l, _ := logger.InitLogger()

	handler.Logger = l

	filmRouter.Get("/", middleware.Protected(userRepo), handler.FetchPaginatedFilms)
	filmRouter.Get("/:id", middleware.Protected(userRepo), handler.FetchSingleFilm)
}

func (h *FilmHandler) FetchPaginatedFilms(c *fiber.Ctx) error {
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

	data, err := h.FilmUsecase.FetchFilmsFromAllSources(context.TODO(), int64(pageInt), int64(limitInt))

	if err != nil {
		return domain.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  data,
	})
}

func (h *FilmHandler) FetchSingleFilm(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := h.FilmRepo.GetById(context.TODO(), id)

	if err != nil {
		return domain.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  data,
	})
}
