package mongodb

import (
	"context"
	"errors"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"movies-review-api/domain"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoFilmRepository struct {
	Logger *zap.Logger
	Coll   *mgm.Collection
}

func (m *mongoFilmRepository) FetchPaginatedFilms(ctx context.Context, page, limit int64) (*domain.PaginatedFilm, error) {

	var films []domain.Film

	collection := mgm.Coll(&domain.Film{}).Collection

	filter := bson.D{}

	paginatedData, err := mongopagination.New(collection).
		Context(ctx).
		Limit(limit).
		Page(page).
		Sort("release_date", 1).
		Filter(filter).
		Decode(&films).
		Find()

	if err != nil {
		return nil, err
	}

	return &domain.PaginatedFilm{
		Data:       films,
		Pagination: paginatedData,
	}, nil
}

func (m *mongoFilmRepository) GetById(ctx context.Context, id string) (*domain.Film, error) {
	var film domain.Film

	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid resource id")
	}

	err = m.Coll.FindByIDWithCtx(ctx, primitiveId, &film)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("resource not found")
		}
		m.Logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &film, nil
}

func NewFilmRepository(logger *zap.Logger) domain.FilmRepository {
	return &mongoFilmRepository{
		Logger: logger,
		Coll:   mgm.Coll(&domain.Film{}),
	}
}
