package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"workspace-onboarding-service/internal/models"
	"workspace-onboarding-service/internal/response"
	"workspace-onboarding-service/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	u, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, u)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	u, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, u)
}

// List handles GET /users and GET /users?organization_id=1
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var orgID *int64
	if v := q.Get("organization_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid organization_id")
			return
		}
		orgID = &id
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	users, err := h.svc.List(r.Context(), orgID, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}


func (h *UserHandler) ListByOrganization(w http.ResponseWriter, r *http.Request) {
	orgID, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid organization id")
		return
	}
	users, err := h.svc.List(r.Context(), &orgID, 100, 0)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	u, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, u)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, name), 10, 64)
}

// writeServiceError is the single place that maps domain errors to HTTP
// status codes for the user module — reused by every handler method
// instead of duplicating this switch five times.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		response.Error(w, http.StatusNotFound, "user not found")
	case errors.Is(err, models.ErrDuplicate):
		response.Error(w, http.StatusConflict, "a user with this email already exists")
	case errors.Is(err, models.ErrOrgNotFound):
		response.Error(w, http.StatusBadRequest, "organization does not exist")
	case errors.Is(err, models.ErrValidation):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}