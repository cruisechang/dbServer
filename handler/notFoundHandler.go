package handler

import (
	"net/http"
)

//NewNotFoundHandler returns NotFoundHandler structure
func NewNotFoundHandler(base baseHandler) *NotFoundHandler {
	return &NotFoundHandler{
		baseHandler: base,
	}
}

//NotFoundHandler returns http status NotFound
type NotFoundHandler struct {
	baseHandler
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.writeError(w, http.StatusNotFound, CodeSuccess, "")
}
