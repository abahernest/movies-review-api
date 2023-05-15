package user

import (
	"context"
	"errors"
	"movies-review-api/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func (u userUsecase) Login(ctx context.Context, data *domain.LoginRequest) (*domain.User, error) {
	existingUser, err := u.userRepo.GetByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
	}

	// check password hash
	if isCorrect := domain.CheckPasswordHash(data.Password, existingUser.Password); !isCorrect {
		return nil, errors.New("invalid login credentials")
	}
	return existingUser, nil
}

func (u userUsecase) Signup(ctx context.Context, data *domain.SignupRequest) (*domain.User, error) {
	existingUser, err := u.userRepo.GetByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user := domain.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
	}

	// hash password
	user.Password, err = domain.HashPassword(data.Password)

	newUser, err := u.userRepo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func New(u domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: u,
	}
}
