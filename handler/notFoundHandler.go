package handler

import (
	"net/http"
)

func NewNotFoundHandler(base baseHandler) *NotFoundHandler {
	return &NotFoundHandler{
		baseHandler: base,
	}
}

type NotFoundHandler struct {
	baseHandler
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.writeError(w, http.StatusNotFound, CodeSuccess, "")
}
