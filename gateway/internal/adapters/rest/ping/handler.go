package ping

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Hubcher/project-management/gateway/internal/core"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

// NewPingHandler godoc
// @Summary Check service availability
// @Description Calls downstream services and returns their availability status from the gateway perspective.
// @Tags system
// @Produce json
// @Success 200 {object} PingResponse
// @Router /api/ping [get]
func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := PingResponse{Replies: make(map[string]string)}

		for name, pinger := range pingers {
			if err := pinger.Ping(r.Context()); err != nil {
				reply.Replies[name] = "unavailable"
				log.Error("one of services is not available", "service", name)
				continue
			}
			reply.Replies[name] = "ok"
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}