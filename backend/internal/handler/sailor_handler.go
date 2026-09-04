package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type SailorHandler struct {
	svc *service.SailorService
}

func NewSailorHandler(svc *service.SailorService) *SailorHandler {
	return &SailorHandler{svc: svc}
}

// Create creates a new sailor
// @Summary      Create sailor
// @Description  Registers a new sailor (participant) in a regatta. The combination (regatta_id, boat_class, sail_number) must be unique.
// @Tags         sailors
// @Accept       json
// @Produce      json
// @Param        sailor  body  model.Sailor  true  "Sailor details"
// @Success      201  {object}  model.Sailor
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sailors [post]
func (h *SailorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s model.Sailor
	if err := decodeJSON(r, &s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Create(r.Context(), &s); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

// GetByID returns a sailor by ID
// @Summary      Get sailor by ID
// @Description  Returns a single sailor by UUID.
// @Tags         sailors
// @Produce      json
// @Param        id  path  string  true  "Sailor UUID"
// @Success      200  {object}  model.Sailor
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "sailor not found"
// @Security     BearerAuth
// @Router       /sailors/{id} [get]
func (h *SailorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	s, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "sailor not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// ListByRegatta returns all sailors in a regatta
// @Summary      List sailors by regatta
// @Description  Returns all sailors registered in a given regatta.
// @Tags         sailors
// @Produce      json
// @Param        regatta_id  path  string  true  "Regatta UUID"
// @Success      200  {array}  model.Sailor
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /regattas/{regatta_id}/sailors [get]
func (h *SailorHandler) ListByRegatta(w http.ResponseWriter, r *http.Request) {
	regattaID, err := uuid.Parse(r.PathValue("regatta_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid regatta_id")
		return
	}
	sailors, err := h.svc.ListByRegatta(r.Context(), regattaID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sailors == nil {
		sailors = []model.Sailor{}
	}
	writeJSON(w, http.StatusOK, sailors)
}

// Update updates a sailor
// @Summary      Update sailor
// @Description  Updates an existing sailor's details. The ID in the path overrides any ID in the body.
// @Tags         sailors
// @Accept       json
// @Produce      json
// @Param        id      path  string        true  "Sailor UUID"
// @Param        sailor  body  model.Sailor  true  "Updated sailor fields"
// @Success      200  {object}  model.Sailor
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sailors/{id} [put]
func (h *SailorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var s model.Sailor
	if err := decodeJSON(r, &s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	s.ID = id
	if err := h.svc.Update(r.Context(), &s); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// Delete deletes a sailor
// @Summary      Delete sailor
// @Description  Deletes a sailor and all associated documents and race links.
// @Tags         sailors
// @Param        id  path  string  true  "Sailor UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sailors/{id} [delete]
func (h *SailorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// FindByClassNumber finds a sailor by boat class and sail number
// @Summary      Find sailor by class and sail number
// @Description  Finds a sailor by their boat class and sail number within a specific regatta. Used by judges to link race measurements to sailors.
// @Tags         sailors
// @Produce      json
// @Param        regatta_id   path  string  true  "Regatta UUID"
// @Param        boat_class   query  string  true  "Boat class (e.g. Optimist, Laser)"
// @Param        sail_number  query  string  true  "Sail number on the boat"
// @Success      200  {object}  model.Sailor
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "sailor not found"
// @Security     BearerAuth
// @Router       /regattas/{regatta_id}/sailors/find [get]
func (h *SailorHandler) FindByClassNumber(w http.ResponseWriter, r *http.Request) {
	regattaID, err := uuid.Parse(r.PathValue("regatta_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid regatta_id")
		return
	}
	boatClass := r.URL.Query().Get("boat_class")
	sailNumber := r.URL.Query().Get("sail_number")
	if boatClass == "" || sailNumber == "" {
		writeError(w, http.StatusBadRequest, "boat_class and sail_number are required")
		return
	}
	s, err := h.svc.FindByClassAndNumber(r.Context(), regattaID, boatClass, sailNumber)
	if err != nil {
		writeError(w, http.StatusNotFound, "sailor not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}