package gateway

import (
	"api_gateway/internal/apperror"
	"api_gateway/internal/gateway/routes/user_service"
	"api_gateway/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler struct {
	Logger         logging.Logger
	GatewayService Service
}

// Register adds handler functions to endpoints
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, user_service.CreateByPhoneNumberURL, apperror.Middleware(h.CreateUserByPhoneNumber))
	router.HandlerFunc(http.MethodPost, user_service.CreateByVkURL, apperror.Middleware(h.CreateUserByVk))
	router.HandlerFunc(http.MethodGet, user_service.SendVerificationURL, apperror.Middleware(h.SendVerificationCode))
	router.HandlerFunc(http.MethodGet, user_service.GetByIdURL, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, user_service.GetByPhoneNumberURL, apperror.Middleware(h.GetUserByPhoneNumber))
	router.HandlerFunc(http.MethodGet, user_service.GetUsersURL, apperror.Middleware(h.GetUsers))
	router.HandlerFunc(http.MethodPut, user_service.UpdateByIdURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, user_service.DeleteByIdURL, apperror.Middleware(h.DeleteUser))

}

func (h *Handler) CreateUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE USER BY PHONE NUMBER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) CreateUserByVk(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE USER BY VK")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) SendVerificationCode(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("SEND CODE")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USER BY ID")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) GetUserByPhoneNumber(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USER BY PHONE NUMBER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("GET USERS")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("UPDATE USER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("DELETE USER")
	request, err := h.GatewayService.CopyRequest(r, "user_service", user_service.Port)
	if err != nil {
		return err
	}
	bytes, err := h.GatewayService.DoRequest(request)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}
