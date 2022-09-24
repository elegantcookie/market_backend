package user

import (
	"auth_service/internal/apperror"
	"auth_service/pkg/logging"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const signInURL = "/api/auth/sign-in"

type Handler struct {
	Logger      logging.Logger
	AuthService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, signInURL, apperror.Middleware(h.SignIn))
}

// Sign in
// @Summary Sign in/sign up endpoint
// @Accept json
// @Produce json
// @Tags Auth
// @Success 201
// @Failure 400
// @Router /api/auth/sign-in [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	// HEADER, BODY etc.
	h.Logger.Println("POST SIGN IN")
	w.Header().Set("Content-Type", "application/json")

	token, err := h.AuthService.SignIn(r.Context(), r.Body)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s", signInURL))
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusCreated)
	return nil
}
