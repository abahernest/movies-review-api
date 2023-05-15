package film

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"movies-review-api/domain"
	"movies-review-api/pkg/logger"
)

type filmUsecase struct {
	filmRepo domain.FilmRepository
	logger   *zap.Logger
}

func (u filmUsecase) FetchFilmsFromAllSources(ctx context.Context, page, limit int64) (*domain.PaginatedFilm, error) {

	var externalSource string
	// check external source for new films
	err := domain.UpdateFilmFromSource(u.logger, externalSource)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error occured while updating film from source %s", zap.Error(err)))
	}

	data, err := u.filmRepo.FetchPaginatedFilms(ctx, page, limit)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func New(u domain.FilmRepository) domain.FilmUsecase {
	l, _ := logger.InitLogger()

	return &filmUsecase{
		filmRepo: u,
		logger:   l,
	}
}
