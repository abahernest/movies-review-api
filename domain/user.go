package domain

// take a look at
// https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
// for struct validation

import (
	"context"
	"time"

	"github.com/Kamva/mgm/v2"
	//"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Firstname        string     `validate:"required" json:"firstname" bson:"firstname"`
	Lastname         string     `validate:"required" json:"lastname" bson:"lastname"`
	Email            string     `validate:"required" json:"email" bson:"email"`
	Password         string     `validate:"required" json:"password,omitempty" bson:"password"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
}

func (u *User) Default() interface{} {
	return &User{
		DeletedAt: nil,
	}
}

type SignupRequest struct {
	Firstname string `validate:"required" json:"firstname" bson:"firstname"`
	Lastname  string `validate:"required" json:"lastname" bson:"lastname"`
	Email     string `validate:"required" json:"email" bson:"email"`
	Password  string `validate:"required" json:"password" bson:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required" json:"email" bson:"email"`
	Password string `validate:"required" json:"password" bson:"password"`
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, userId string) (*User, error)
}

type UserUsecase interface {
	Login(ctx context.Context, reqBody *LoginRequest) (*User, error)
	Signup(ctx context.Context, reqBody *SignupRequest) (*User, error)
}
