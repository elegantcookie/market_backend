package user

import (
	"context"
	"fmt"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"user_service/pkg/client/twilio"
	"user_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	cache   Cache
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, userCache Cache, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		cache:   userCache,
		logger:  logger,
	}, nil
}

type Service interface {
	SignInByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error)
	SendCode(ctx context.Context, phoneNumber string) error
	CreateByVk(ctx context.Context, dto CreateByVkDTO) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, uuid string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (s service) SignInByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error) {
	val, err := s.cache.Get(ctx, dto.PhoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to create verification check: %v", err)
	}
	if val != dto.VerificationCode {
		return "", fmt.Errorf("wrong code, right one: %s", val)
	}

	// TODO check if there is an account with existing number, then just return 200
	user := NewUserByPhone(dto)
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	return userID, nil

}

func (s service) SendCode(ctx context.Context, phoneNumber string) error {
	s.logger.Info("SEND CODE SERVICE")

	client := twilio.GetClient()
	s.logger.Info("SET UP CLIENT")

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(client.FromPhone)

	code := generateVerificationCode(6)
	message := fmt.Sprintf(verificationMessage, code)

	params.SetBody(message)

	resp, err := client.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		s.logger.Infof("failed to create verification: %v", err)
		return fmt.Errorf("failed to create verification: %v", err)
	}

	err = s.cache.Set(ctx, phoneNumber, code, deleteRegisterNumberTime)
	if err != nil {
		return err
	}
	s.logger.Infof("Sent notification: %+v", resp)
	return nil
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
