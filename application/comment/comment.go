package comment

import (
	"context"
	"movies-review-api/domain"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

func (u commentUsecase) AddComment(ctx context.Context, data *domain.NewCommentRequest) (*domain.Comment, error) {

	comment := domain.Comment{
		FilmId:  data.FilmId,
		Summary: data.Summary,
		UserId:  data.UserId,
	}

	newComment, err := u.commentRepo.Create(ctx, &comment)

	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func New(u domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{
		commentRepo: u,
	}
}
