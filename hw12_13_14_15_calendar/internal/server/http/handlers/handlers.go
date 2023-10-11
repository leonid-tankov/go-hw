package handlers

import (
	"net/http"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
)

type HTTPHandler struct {
	logger repository.Logger
}

func NewHTTPHandler(logger repository.Logger) *HTTPHandler {
	return &HTTPHandler{
		logger: logger,
	}
}

func (hh *HTTPHandler) HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
