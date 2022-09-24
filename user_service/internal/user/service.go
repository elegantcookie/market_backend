package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"user_service/internal/apperror"
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
	SignInByVk(ctx context.Context, dto CreateByVkDTO) (string, error)
	SendCode(ctx context.Context, phoneNumber string) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, uuid string) (User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (s service) SignInByPhone(ctx context.Context, dto CreateByPhoneDTO) (string, error) {
	s.logger.Infof("dto: %+v", dto)
	// Gets sent verification code from redis
	val, err := s.cache.Get(ctx, dto.PhoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to create verification check: %v", err)
	}
	// Compares code is correct to one, sent by user
	if val != dto.VerificationCode {
		// TODO remove right code from error message
		return "", fmt.Errorf("wrong code, right one: %s", val)
	}

	// Returns id of the existing user if it is, instead of creating new one
	if user, err := s.storage.FindByNumber(ctx, dto.PhoneNumber); err == nil {
		return user.ID, nil
	}
	// Returns id of a new user if it is not
	user := NewUserByPhone(dto)
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	return userID, nil
}

func (s service) SignInByVk(ctx context.Context, dto CreateByVkDTO) (string, error) {

	// Returns error if token is invalid
	vkID, err := s.CheckVkToken(ctx, dto.VkToken)
	if err != nil {
		return "", err
	}

	// Returns id of the existing user if it is, instead of creating new one
	if user, err := s.storage.FindByVkID(ctx, vkID); err == nil {
		return user.ID, nil
	}
	// Returns id of a new user if it is not
	user := NewUserByVkID(dto)
	user.VkID = vkID
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	return userID, nil
}

// For dev
func (s service) SendCode(ctx context.Context, phoneNumber string) (string, error) {
	s.logger.Info("SEND CODE SERVICE")

	code := generateVerificationCode(verificationCodeLength)

	// Fills verificationMessage template with code
	message := fmt.Sprintf(verificationMessage, code)

	// Adds to cache phoneNumber as a key and right code as a value to compare at SignInByPhone
	err := s.cache.Set(ctx, phoneNumber, code, deleteRegisterNumberTime)
	if err != nil {
		return "", err
	}
	return message, nil
}

func (s service) CheckVkToken(ctx context.Context, token string) (string, error) {
	var accessKey = os.Getenv("VK_ACCESS_TOKEN")
	s.logger.Printf("VK_ACCESS_TOKEN: %s", accessKey)
	url := fmt.Sprintf(checkTokenURL, token, accessKey, checkTokenVersion)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response == nil {
		return "", fmt.Errorf("response is null")
	}

	var dto VkCheckTokenDTO
	err = json.NewDecoder(response.Body).Decode(&dto)
	if err != nil {
		return "", err
	}
	//s.logger.Printf("%+v", dto)
	// if vk api returned response.success != 1 then vk access token is invalid
	if dto.Response.Success != 1 {
		s.logger.Printf("wrong response code, success: %d", dto.Response.Success)
		return "", fmt.Errorf("wrong response code, success: %d", dto.Response.Success)
	}
	return strconv.Itoa(dto.Response.UserID), nil
}

// For prod
//func (s service) SendCode(ctx context.Context, phoneNumber string) error {
//	s.logger.Info("SEND CODE SERVICE")
//
//	// Initializes twilio client
//	client := twilio.GetClient()
//	s.logger.Info("SET UP CLIENT")
//
//	// Adds params:
//	// SetTo arg - number of user, will be sending verification code
//	// SetFrom arg - twilio number from account
//	params := &openapi.CreateMessageParams{}
//	params.SetTo(phoneNumber)
//	params.SetFrom(client.FromPhone)
//
//	code := generateVerificationCode(verificationCodeLength)
//
//	// Fills verificationMessage template with code
//	message := fmt.Sprintf(verificationMessage, code)
//
//	// Sets message text to be sent
//	params.SetBody(message)
//
//	// Creates message via twilio api
//	resp, err := client.TwilioClient.Api.CreateMessage(params)
//	if err != nil {
//		s.logger.Infof("failed to create verification: %v", err)
//		return fmt.Errorf("failed to create verification: %v", err)
//	}
//
//	// Adds to cache phoneNumber as a key and right code as a value to compare at SignInByPhone
//	err = s.cache.Set(ctx, phoneNumber, code, deleteRegisterNumberTime)
//	if err != nil {
//		return err
//	}
//	s.logger.Infof("Sent notification: %+v", resp)
//	return nil
//}

func (s service) GetAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to find user_service. error: %v", err)
	}
	return users, nil
}
func (s service) GetById(ctx context.Context, uuid string) (User, error) {
	user, err := s.storage.FindById(ctx, uuid)
	if err != nil {
		// Returns 404 if an entity not found
		if errors.Is(err, apperror.ErrNotFound) {
			return user, err
		}
		return user, fmt.Errorf("failed to find user by uuid: %w", err)
	}
	return user, nil
}
func (s service) GetByPhoneNumber(ctx context.Context, phoneNumber string) (u User, err error) {
	s.logger.Infof("phone number: %s", phoneNumber)
	u, err = s.storage.FindByNumber(ctx, phoneNumber)
	if err != nil {
		s.logger.Infof("error: %+v", err)
		if errors.Is(err, apperror.ErrNotFound) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find user by phone number: %v", err)
	}
	s.logger.Infof("user: %+v", u)
	return u, nil
}
func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	user := NewUserFromUpdateDTO(dto)
	err := s.storage.Update(ctx, user)
	if err != nil {
		// Returns 404 if an entity not found
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil
}
func (s service) Delete(ctx context.Context, uuid string) error {
	err := s.storage.Delete(ctx, uuid)

	if err != nil {
		// Returns 404 if an entity not found
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}
