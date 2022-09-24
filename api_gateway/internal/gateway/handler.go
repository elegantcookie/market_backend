package gateway

import (
	"api_gateway/internal/apperror"
	"api_gateway/internal/gateway/routes/auth_service"
	"api_gateway/internal/gateway/routes/user_service"
	"api_gateway/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler struct {
	Logger         logging.Logger
	GatewayService Service
}

// Register adds handler functions to endpoints
func (h *Handler) Register(router *httprouter.Router) {
	/* =================================================================================================================
		User service routes
	================================================================================================================= */
	router.HandlerFunc(http.MethodPost, user_service.CreateByPhoneNumberURL, apperror.Middleware(h.CreateUserByPhoneNumber))
	router.HandlerFunc(http.MethodPost, user_service.CreateByVkURL, apperror.Middleware(h.CreateUserByVk))
	router.HandlerFunc(http.MethodGet, user_service.SendVerificationURL, apperror.Middleware(h.SendVerificationCode))
	router.HandlerFunc(http.MethodGet, user_service.GetByIdURL, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, user_service.GetByPhoneNumberURL, apperror.Middleware(h.GetUserByPhoneNumber))
	router.HandlerFunc(http.MethodGet, user_service.GetUsersURL, apperror.AuthMiddleware(h.GetUsers))
	router.HandlerFunc(http.MethodPut, user_service.UpdateByIdURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, user_service.DeleteByIdURL, apperror.Middleware(h.DeleteUser))

	/* =================================================================================================================
		Auth service routes
	================================================================================================================= */

	router.HandlerFunc(http.MethodGet, auth_service.CreateTokenURL, apperror.Middleware(h.AuthCreateToken))
	router.HandlerFunc(http.MethodPost, auth_service.CheckTokenURL, apperror.Middleware(h.AuthCheckToken))
}

/* =====================================================================================================================
	User service handlers
===================================================================================================================== */

func (h *Handler) CreateUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE USER BY PHONE NUMBER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	if data.StatusCode != http.StatusCreated && data.StatusCode != http.StatusOK {
		w.WriteHeader(data.StatusCode)
		w.Write(data.Body)
		return nil
	}

	// Get id from user_service response
	var dto CreateUserResponse
	err = json.Unmarshal(data.Body, &dto)
	if err != nil {
		return err
	}

	// Create JWT using user id
	JWTDTO, err := h.GatewayService.CreateJWT(r, dto.ID)
	if err != nil {
		return err
	}

	// Sets JWT access token as response header
	// TODO add refresh token
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", JWTDTO.AccessToken))
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) CreateUserByVk(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE USER BY VK")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	if data.StatusCode != http.StatusCreated && data.StatusCode != http.StatusOK {
		w.WriteHeader(data.StatusCode)
		w.Write(data.Body)
		return nil
	}

	// Get id from user_service response
	var dto CreateUserResponse
	err = json.Unmarshal(data.Body, &dto)
	if err != nil {
		return err
	}

	// Create JWT using user id
	JWTDTO, err := h.GatewayService.CreateJWT(r, dto.ID)
	if err != nil {
		return err
	}

	// Sets JWT access token as response header
	// TODO add refresh token
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", JWTDTO.AccessToken))
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) SendVerificationCode(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("SEND CODE")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USER BY ID")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) GetUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USER BY PHONE NUMBER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USERS")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("UPDATE USER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("DELETE USER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}

/* =====================================================================================================================
	Auth service handlers
===================================================================================================================== */

func (h *Handler) AuthCheckToken(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CHECK TOKEN")
	request, err := h.GatewayService.CopyRequest(r, auth_service.DockerName, auth_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}

func (h *Handler) AuthCreateToken(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE TOKEN")
	request, err := h.GatewayService.CopyRequest(r, auth_service.DockerName, auth_service.Port)
	if err != nil {
		return err
	}
	data, err := h.GatewayService.DoRequest(request)

	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	w.Write(data.Body)
	return nil
}
