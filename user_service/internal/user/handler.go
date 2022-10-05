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

// Endpoints paths
const (
	createByPhoneNumberURL = "/api/v1/user_service/phone"
	createByVkURL          = "/api/v1/user_service/vk/:vk_token/"
	sendVerificationURL    = "/api/v1/user_service/send_verification/:phone_number/"
	getUsersURL            = "/api/v1/user_service/get/all"
	getByIdURL             = "/api/v1/user_service/get/id/:id/"
	getByPhoneNumberURL    = "/api/v1/user_service/get/num/:phone_number/"
	updateByIdURL          = "/api/v1/user_service/upd/"
	deleteByIdURL          = "/api/v1/user_service/del/id/:id/"
)

// Register adds handler functions to endpoints
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createByPhoneNumberURL, apperror.Middleware(h.CreateUserByPhoneNumber))
	router.HandlerFunc(http.MethodPost, createByVkURL, apperror.Middleware(h.CreateUserByVk))
	router.HandlerFunc(http.MethodGet, sendVerificationURL, apperror.Middleware(h.SendVerificationCode))
	router.HandlerFunc(http.MethodGet, getByIdURL, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, getByPhoneNumberURL, apperror.Middleware(h.GetUserByPhoneNumber))
	router.HandlerFunc(http.MethodGet, getUsersURL, apperror.Middleware(h.GetUsers))
	router.HandlerFunc(http.MethodPut, updateByIdURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, deleteByIdURL, apperror.Middleware(h.DeleteUser))

}

// CreateUserByPhoneNumber adds new user to db
// @Summary Create user by phone number and token
// @Accept json
// @Produce json
// @Param data body CreateByPhoneDTO true "structure holds data for user creation by phone"
// @Tags Users
// @Success 201 If new user has been signed in
// @Success 200 If existing user has been signed in
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/phone [post]
func (h *Handler) CreateUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	// TODO: refactor

	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	var crUser CreateByPhoneDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	response, err := h.UserService.SignInByPhone(r.Context(), crUser)
	if err != nil {
		return err
	}
	tmp := make(map[string]string)
	tmp["id"] = response.ID
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if response.IsNewUser {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(bytes)
	return nil
}

// CreateUserByVk adds new user to db
// @Summary Create user by phone number and token
// @Accept json
// @Produce json
// @Param vk_token path string true "VK Access token"
// @Tags Users
// @Success 201 If new user has been signed in
// @Success 200 If existing user has been signed in
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/vk/{vk_token}/ [post]
func (h *Handler) CreateUserByVk(w http.ResponseWriter, r *http.Request) error {
	// TODO: refactor

	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	vkToken := params.ByName("vk_token")

	response, err := h.UserService.SignInByVk(r.Context(), vkToken)
	if err != nil {
		return err
	}
	tmp := make(map[string]string)
	tmp["id"] = response.ID
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if response.IsNewUser {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
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
// @Router /api/v1/user_service/send_verification/{phone_number}/ [get]
func (h *Handler) SendVerificationCode(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("SEND CODE")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	phoneNumber := params.ByName("phone_number")
	dummyMessage, err := h.UserService.SendCode(r.Context(), phoneNumber)
	if err != nil {
		h.Logger.Errorf("%v", err)
		return err
	}
	h.Logger.Info("OK")
	w.WriteHeader(http.StatusOK)

	// TODO delete on release
	w.Write([]byte(dummyMessage))
	return nil
}

// GetUserById swaggo
// @Summary Returns user data by id
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Tags Users
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/get/id/{id} [get]
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("id")

	user, err := h.UserService.GetById(r.Context(), userUUID)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// GetUserByPhoneNumber swaggo
// @Summary Returns user data by phone number
// @Accept json
// @Produce json
// @Param phone_number path string true "Phone number"
// @Tags Users
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/get/id/{phone_number}/ [get]
func (h *Handler) GetUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	username := params.ByName("phone_number")

	user, err := h.UserService.GetByPhoneNumber(r.Context(), username)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// GetUsers swaggo
// @Summary Returns data of all user_service
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/v1/user_service/get/all [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserService.GetAll(r.Context())
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// UpdateUser swaggo
// @Summary Updates user data
// @Accept json
// @Produce json
// @Param data body UpdateUserDTO true "update user struct"
// @Tags Users
// @Success 204
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/upd/ [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var dto UpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.UserService.Update(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// DeleteUser swaggo
// @Summary Deletes user by id
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Tags Users
// @Success 204
// @Failure 400 {object} apperror.AppError
// @Router /api/v1/user_service/del/id/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("id")

	err := h.UserService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
