package mongodb

import (
	"context"
	"movies-review-api/domain"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoUserRepository struct {
	Logger *zap.Logger
	Coll   *mgm.Collection
}

func (m mongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var users []domain.User

	err := m.Coll.SimpleFindWithCtx(ctx, &users, bson.D{{Key: "email", Value: email}})

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func (m mongoUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

	err := m.Coll.CreateWithCtx(ctx, user)

	if err != nil {
		m.Logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func NewUserRepository(logger *zap.Logger) domain.UserRepository {
	return &mongoUserRepository{
		Logger: logger,
		Coll:   mgm.Coll(&domain.User{}),
	}
}
