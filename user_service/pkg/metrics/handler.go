package metrics

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const HeartbeatURL = "/api/v1/user_service/heartbeat"

type Handler struct {
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, HeartbeatURL, h.Heartbeat)
}

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/v1/user_service/heartbeat [get]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	// Check if the server is up
	w.WriteHeader(204)
}
