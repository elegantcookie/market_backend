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
)

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createByPhoneNumberURL, apperror.Middleware(h.CreateUserByPhoneNumber))
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

	h.Logger.Debug("decode create user dto")
	var crUser CreateByPhoneDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	userUUID, err := h.UserService.CreateByPhone(r.Context(), crUser)
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
