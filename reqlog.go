package reqlog

import (
	"log"
	"net/http"
	"time"
)

type reqlogHandler struct {
	handler http.Handler
	logger  *log.Logger
}

func NewHandler(handler http.Handler, logger *log.Logger) http.Handler {
	return &reqlogHandler{
		handler: handler,
		logger: logger,
	}
}

func (h *reqlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.handler.ServeHTTP(w, r)
	end := time.Now()
	h.logger.Printf("[request] %s %s [%s]", r.Method, r.URL, end.Sub(start))
}
