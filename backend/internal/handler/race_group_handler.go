package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type RaceGroupHandler struct {
	svc *service.RaceGroupService
}

func NewRaceGroupHandler(svc *service.RaceGroupService) *RaceGroupHandler {
	return &RaceGroupHandler{svc: svc}
}

// Create creates a new race group (heat)
// @Summary      Create race group
// @Description  Creates a new race group (a single start session with multiple boats). Groups race measurements together.
// @Tags         race-groups
// @Accept       json
// @Produce      json
// @Param        group  body  model.RaceGroup  true  "Race group details"
// @Success      201  {object}  model.RaceGroup
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /race-groups [post]
func (h *RaceGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var rg model.RaceGroup
	if err := decodeJSON(r, &rg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Create(r.Context(), &rg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rg)
}

// GetByID returns a race group by ID
// @Summary      Get race group by ID
// @Description  Returns a single race group by UUID.
// @Tags         race-groups
// @Produce      json
// @Param        id  path  string  true  "Race group UUID"
// @Success      200  {object}  model.RaceGroup
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "race group not found"
// @Security     BearerAuth
// @Router       /race-groups/{id} [get]
func (h *RaceGroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	rg, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "race group not found")
		return
	}
	writeJSON(w, http.StatusOK, rg)
}

// ListByRegatta returns all race groups for a regatta
// @Summary      List race groups by regatta
// @Description  Returns all race groups (heats) for a given regatta, ordered by group number.
// @Tags         race-groups
// @Produce      json
// @Param        regatta_id  path  string  true  "Regatta UUID"
// @Success      200  {array}  model.RaceGroup
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas/{regatta_id}/race-groups [get]
func (h *RaceGroupHandler) ListByRegatta(w http.ResponseWriter, r *http.Request) {
	regattaID, err := uuid.Parse(r.PathValue("regatta_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid regatta_id")
		return
	}
	groups, err := h.svc.ListByRegatta(r.Context(), regattaID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if groups == nil {
		groups = []model.RaceGroup{}
	}
	writeJSON(w, http.StatusOK, groups)
}

// Start sets a race group status to "running"
// @Summary      Start race group
// @Description  Changes the race group status to "running", indicating the heat has started.
// @Tags         race-groups
// @Param        id  path  string  true  "Race group UUID"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /race-groups/{id}/start [post]
func (h *RaceGroupHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.StartRaceGroup(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Finish sets a race group status to "finished"
// @Summary      Finish race group
// @Description  Changes the race group status to "finished", indicating the heat has ended and no more measurements can be added.
// @Tags         race-groups
// @Param        id  path  string  true  "Race group UUID"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /race-groups/{id}/finish [post]
func (h *RaceGroupHandler) Finish(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.FinishRaceGroup(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete deletes a race group
// @Summary      Delete race group
// @Description  Deletes a race group and all associated race measurements.
// @Tags         race-groups
// @Param        id  path  string  true  "Race group UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /race-groups/{id} [delete]
func (h *RaceGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
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