package user

import (
	"context"
	"fmt"
	"user_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	CreateByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error)
	CreateByVk(ctx context.Context, dto CreateByVkDTO) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, uuid string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (s service) CreateByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error) {
	user := NewUserByPhone(dto)
	err := user.GeneratePhoneTokenHash()
	if err != nil {
		return "", fmt.Errorf("failed to generate phone token hash due to: %v", err)
	}
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	return userID, nil

}
func (s service) CreateByVk(ctx context.Context, dto CreateByVkDTO) (string, error) {
	return "", nil
}
func (s service) GetAll(ctx context.Context) ([]User, error) {
	return nil, nil
}
func (s service) GetById(ctx context.Context, uuid string) (User, error) {
	return User{}, nil
}
func (s service) GetByUsername(ctx context.Context, username string) (User, error) {
	return User{}, nil
}
func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	return nil
}
func (s service) Delete(ctx context.Context, uuid string) error {
	return nil
}
