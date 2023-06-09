package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
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
	var user domain.User

	err := m.Coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (m mongoUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

	err := m.Coll.CreateWithCtx(ctx, user)

	if err != nil {
		m.Logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (m mongoUserRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := m.Coll.FindByIDWithCtx(ctx, id, &user)

	if err != nil {
		m.Logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &user, nil
}
func NewUserRepository(logger *zap.Logger) domain.UserRepository {
	return &mongoUserRepository{
		Logger: logger,
		Coll:   mgm.Coll(&domain.User{}),
	}
}
