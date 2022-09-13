package user

import (
	"context"
	"fmt"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
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
	CreateByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error)
	CreateByVk(ctx context.Context, dto CreateByVkDTO) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, uuid string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (s service) CreateByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error) {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(dto.PhoneNumber)
	params.SetCode(dto.VerificationCode)

	client := twilio.GetClient()
	resp, err := client.TwilioClient.VerifyV2.CreateVerificationCheck(client.ServiceSID, params)
	if err != nil {
		return "", fmt.Errorf("failed to create verification check: %v", err)
	}
	if *resp.Status != twilioStatusApproved {
		return "", fmt.Errorf("wrong code")
	}
	user := NewUserByPhone(dto)
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	return userID, nil

}

func (s service) SendCode(ctx context.Context, phoneNumber string) error {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")
	client := twilio.GetClient()
	resp, err := client.TwilioClient.VerifyV2.CreateVerification(client.ServiceSID, params)
	if err != nil {
		return fmt.Errorf("failed to create verification: %v", err)
	}
	s.logger.Infof("Sent notification: %v", resp)
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
