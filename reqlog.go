package reqlog

import (
	"log"
	"net/http"
	"time"

	"github.com/stathat/go"
)

type reqlogHandler struct {
	handler http.Handler
	env string
	backends []Backend
}

type Backend interface {
	RecordRequest(*http.Request, time.Duration, string)
}

type Logger struct {
	logger *log.Logger
}

type StatHat struct {
	ezkey string
}

func NewLogger(l *log.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func NewStatHat(ezkey string) *StatHat {
	return &StatHat{
		ezkey: ezkey,
	}
}

func NewHandler(handler http.Handler, env string, backends ...Backend) http.Handler {
	return &reqlogHandler{
		handler: handler,
		env: env,
		backends: backends,
	}
}

func (l *Logger) RecordRequest(r *http.Request, duration time.Duration, env string) {
	l.logger.Printf("(%s) %s %s [%s]", env, r.Method, r.URL, duration)
}

func (s *StatHat) requestCount(r *http.Request, env string) {
	stathat.PostEZCount(env+": site request count", s.ezkey, 1)
}

func (s *StatHat) requestDuration(r *http.Request, duration time.Duration, env string) {
	stathat.PostEZValue(env+": site request duration ms", s.ezkey, float64(duration*time.Millisecond))
}

func (s *StatHat) RecordRequest(r *http.Request, duration time.Duration, env string) {
	go s.requestCount(r, env)
	go s.requestDuration(r, duration, env)
}

func (h *reqlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.handler.ServeHTTP(w, r)
	dur := time.Now().Sub(start)
	for _, backend := range h.backends {
		go backend.RecordRequest(r, dur, h.env)
	}
}
