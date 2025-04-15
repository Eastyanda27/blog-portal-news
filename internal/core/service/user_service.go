package service

import (
	"blognewsportal/internal/adapter/repository"
	"blognewsportal/internal/core/domain/entity"
	"blognewsportal/lib/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	EditPassword(ctx context.Context, newPass string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) EditPassword(ctx context.Context, newPass string, id int64) error {
	password, err := conv.HashPassword(newPass)
	if err != nil {
		code = "[SERVICE] EditPassword - 1"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepository.EditPassword(ctx, password, id)
	if err != nil {
		code = "[SERVICE] EditPassword - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepository.GetUserByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}
