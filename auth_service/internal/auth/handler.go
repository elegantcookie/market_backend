package auth

import (
	"auth_service/internal/apperror"
	"auth_service/pkg/logging"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	checkTokenURL  = "/api/v1/auth_service/token/check"
	createTokenURL = "/api/v1/auth_service/token/create/:id"
)

type Handler struct {
	Logger      logging.Logger
	AuthService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, checkTokenURL, apperror.Middleware(h.CheckToken))
	router.HandlerFunc(http.MethodGet, createTokenURL, apperror.Middleware(h.CreateToken))
}

func (h *Handler) CheckToken(w http.ResponseWriter, r *http.Request) error {
	var dto JWTDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}
	err = h.AuthService.CheckToken(r.Context(), dto.Token)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Println("CREATE TOKEN")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	tokens, err := h.AuthService.CreateToken(r.Context(), userID)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(&tokens)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}
