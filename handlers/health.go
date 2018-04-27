package handlers

import (
	"net/http"
)

// HealthHandler.
type HealthHandler struct {
}

// NewHealthHandler описывает конструктор со структурой HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP перехватывает http запрос /health/*
func (h *HealthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(`{"state":"success","sessionId":null,"value":{},"health":0}`))
}
