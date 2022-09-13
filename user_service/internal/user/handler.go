package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"user_service/internal/apperror"
	"user_service/pkg/logging"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

const (
	createByPhoneNumberURL = "/api/v1/users/phone"
	sendVerificationURL    = "/api/v1/users/send_verification/:phone_number/"
)

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createByPhoneNumberURL, apperror.Middleware(h.CreateUserByPhoneNumber))
	router.HandlerFunc(http.MethodGet, sendVerificationURL, apperror.Middleware(h.SendVerificationCode))
}

// CreateUserByPhoneNumber adds new user to db
// @Summary Create user by phone number and token
// @Accept json
// @Produce json
// @Param data body CreateByPhoneDTO true "structure holds data for user creation by phone"
// @Tags Users
// @Success 201
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/users/phone [post]
func (h *Handler) CreateUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	var crUser CreateByPhoneDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	userUUID, err := h.UserService.SignInByPhone(r.Context(), crUser)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(userUUID)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	w.Write(bytes)
	return nil
}

// SendVerificationCode adds a request to send verification code
// @Summary Sends verification code by phone number
// @Accept json
// @Produce json
// @Param phone_number path string true "phone number where the code will be sent"
// @Tags Users
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/users/send_verification/{phone_number} [get]
func (h *Handler) SendVerificationCode(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("SEND CODE")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	phoneNumber := params.ByName("phone_number")
	err := h.UserService.SendCode(r.Context(), phoneNumber)
	if err != nil {
		h.Logger.Errorf("%v", err)
		return err
	}
	h.Logger.Info("OK")
	w.WriteHeader(http.StatusOK)
	return nil
}
