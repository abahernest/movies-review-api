package mongodb

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/mongo"
	"movies-review-api/domain"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoCommentRepository struct {
	Logger *zap.Logger
	Coll   *mgm.Collection
}

func (m *mongoCommentRepository) FetchPaginatedFilmComments(ctx context.Context, filmId string, page, limit int64) (*domain.PaginatedComment, error) {

	var comment []domain.Comment

	collection := mgm.Coll(&domain.Comment{}).Collection

	filter := bson.M{"film_id": filmId}

	paginatedData, err := mongopagination.New(collection).
		Context(ctx).
		Limit(limit).
		Page(page).
		Sort("created_at", 1).
		Filter(filter).
		Decode(&comment).
		Find()

	if err != nil {
		return nil, err
	}

	return &domain.PaginatedComment{
		Data:       comment,
		Pagination: paginatedData,
	}, nil
}

func (m mongoCommentRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {

	mgm.TransactionWithCtx(ctx, func(session mongo.Session, sc mongo.SessionContext) error {

		cmnt, err := func(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {

			err := m.Coll.CreateWithCtx(ctx, comment)

			if err != nil {
				m.Logger.Error(err.Error(), zap.Error(err))
				return nil, err
			}

			return comment, nil
		}(ctx, comment)

		if err != nil {
			return err
		}

		comment = cmnt;
		return session.CommitTransaction(ctx)
	})

	return comment, nil
}

func NewCommentRepository(logger *zap.Logger) domain.CommentRepository {
	return &mongoCommentRepository{
		Logger: logger,
		Coll:   mgm.Coll(&domain.Comment{}),
	}
}
