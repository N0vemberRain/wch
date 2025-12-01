package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"wch/users/internal/controller"
)

// Handler defines a users handler
type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// GetOrgDetails handles GET/users requests
func (h *Handler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	details, err := h.ctrl.Get(r.Context(), id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
