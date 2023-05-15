package domain

import (
	"context"
	"github.com/Kamva/mgm/v2"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/bson"
)

type Comment struct {
	mgm.DefaultModel `bson:",inline"`
	FilmId           string `json:"film_id" bson:"film_id"`
	UserId           string `json:"user_id" bson:"user_id"`
	Summary          string `json:"summary" bson:"summary"`
}

type PaginatedComment struct {
	Pagination *mongopagination.PaginatedData `json:"pagination" bson:"pagination"`
	Data       []Comment                      `json:"data" bson:"data"`
}

type NewCommentRequest struct {
	FilmId  string `validate:"required" json:"film_id" bson:"film_id"`
	UserId  string `json:"user_id" bson:"user_id"`
	Summary string `validate:"required,max=500" json:"summary" bson:"summary"`
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) (*Comment, error)
	FetchPaginatedFilmComments(ctx context.Context, filmId string, page, limit int64) (*PaginatedComment, error)
}

type CommentUsecase interface {
	AddComment(ctx context.Context, reqBody *NewCommentRequest) (*Comment, error)
}

// After Create Hook. Inherited from mgm.CreateWithCtx
func (m Comment) Created() error {
	filmId, err := primitive.ObjectIDFromHex(m.FilmId)
	if err != nil {
		return err
	}

	_, err = mgm.Coll(&Film{}).UpdateOne(
		context.Background(),
		bson.M{"_id": filmId},
		bson.M{"$inc": bson.M{"comment_count": 1}})

	if err != nil {
		return err
	}

	return nil
}
