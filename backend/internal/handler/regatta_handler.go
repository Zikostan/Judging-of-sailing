package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type RegattaHandler struct {
	svc *service.RegattaService
}

func NewRegattaHandler(svc *service.RegattaService) *RegattaHandler {
	return &RegattaHandler{svc: svc}
}

// List returns all regattas
// @Summary      List regattas
// @Description  Returns all regattas ordered by start date descending.
// @Tags         regattas
// @Produce      json
// @Success      200  {array}  model.Regatta
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas [get]
func (h *RegattaHandler) List(w http.ResponseWriter, r *http.Request) {
	regs, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if regs == nil {
		regs = []model.Regatta{}
	}
	writeJSON(w, http.StatusOK, regs)
}

// Create creates a new regatta
// @Summary      Create regatta
// @Description  Creates a new regatta with title, dates and optional location.
// @Tags         regattas
// @Accept       json
// @Produce      json
// @Param        regatta  body  model.Regatta  true  "Regatta details (id, created_at, updated_at are ignored)"
// @Success      201  {object}  model.Regatta
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas [post]
func (h *RegattaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reg model.Regatta
	if err := decodeJSON(r, &reg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Create(r.Context(), &reg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, reg)
}

// GetByID returns a regatta by ID
// @Summary      Get regatta by ID
// @Description  Returns a single regatta by its UUID.
// @Tags         regattas
// @Produce      json
// @Param        id  path  string  true  "Regatta UUID"
// @Success      200  {object}  model.Regatta
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "regatta not found"
// @Security     BearerAuth
// @Router       /regattas/{id} [get]
func (h *RegattaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	reg, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "regatta not found")
		return
	}
	writeJSON(w, http.StatusOK, reg)
}

// Update updates a regatta
// @Summary      Update regatta
// @Description  Updates an existing regatta. The ID in the path overrides any ID provided in the body.
// @Tags         regattas
// @Accept       json
// @Produce      json
// @Param        id       path  string        true  "Regatta UUID"
// @Param        regatta  body  model.Regatta  true  "Updated regatta fields (id from path is used)"
// @Success      200  {object}  model.Regatta
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas/{id} [put]
func (h *RegattaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var reg model.Regatta
	if err := decodeJSON(r, &reg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	reg.ID = id
	if err := h.svc.Update(r.Context(), &reg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reg)
}

// Delete deletes a regatta
// @Summary      Delete regatta
// @Description  Deletes a regatta and all associated data (race groups, races, sailors, documents).
// @Tags         regattas
// @Param        id  path  string  true  "Regatta UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas/{id} [delete]
func (h *RegattaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}