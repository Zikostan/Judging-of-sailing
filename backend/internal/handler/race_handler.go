package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/middleware"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type RaceHandler struct {
	svc *service.RaceService
}

func NewRaceHandler(svc *service.RaceService) *RaceHandler {
	return &RaceHandler{svc: svc}
}

type RecordRaceRequest struct {
	GroupID    uuid.UUID  `json:"group_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	SailorID   *uuid.UUID `json:"sailor_id,omitempty" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	BoatClass  string     `json:"boat_class" example:"Optimist"`
	SailNumber string     `json:"sail_number" example:"10"`
	FinishTime *float32   `json:"finish_time,omitempty" example:"125.5"`
	Penalty    *string    `json:"penalty,omitempty" example:"DSQ"`
}

type UpdateResultRequest struct {
	FinishTime *float32 `json:"finish_time,omitempty" example:"125.5"`
	Place      *int16   `json:"place,omitempty" example:"1"`
	Penalty    *string  `json:"penalty,omitempty" example:"DSQ"`
}

type LinkSailorRequest struct {
	SailorID uuid.UUID `json:"sailor_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

// Record records a new race measurement
// @Summary      Record race measurement
// @Description  Records a new individual race measurement for a boat. The judge_id is automatically set from the authenticated user. The sailor_id can be null if the sailor is not yet identified.
// @Tags         races
// @Accept       json
// @Produce      json
// @Param        request  body  handler.RecordRaceRequest  true  "Race measurement details"
// @Success      201  {object}  model.Race
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /races [post]
func (h *RaceHandler) Record(w http.ResponseWriter, r *http.Request) {
	var req RecordRaceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	judgeID, _ := middleware.GetUserID(r.Context())
	race := &model.Race{
		GroupID:    req.GroupID,
		SailorID:   req.SailorID,
		JudgeID:    judgeID,
		BoatClass:  req.BoatClass,
		SailNumber: req.SailNumber,
		FinishTime: req.FinishTime,
		Penalty:    req.Penalty,
	}
	if err := h.svc.Record(r.Context(), race); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, race)
}

// GetByID returns a race measurement by ID
// @Summary      Get race by ID
// @Description  Returns a single race measurement by UUID.
// @Tags         races
// @Produce      json
// @Param        id  path  string  true  "Race UUID"
// @Success      200  {object}  model.Race
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "race not found"
// @Security     BearerAuth
// @Router       /races/{id} [get]
func (h *RaceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	race, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "race not found")
		return
	}
	writeJSON(w, http.StatusOK, race)
}

// ListByGroup returns all race measurements in a race group
// @Summary      List races by group
// @Description  Returns all race measurements for a given race group, ordered by boat class and finish time.
// @Tags         races
// @Produce      json
// @Param        group_id  path  string  true  "Race group UUID"
// @Success      200  {array}  model.Race
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /race-groups/{group_id}/races [get]
func (h *RaceHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group_id")
		return
	}
	races, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if races == nil {
		races = []model.Race{}
	}
	writeJSON(w, http.StatusOK, races)
}

// UpdateResult updates the result of a race measurement
// @Summary      Update race result
// @Description  Updates the finish time, place, and penalty for a race measurement.
// @Tags         races
// @Accept       json
// @Produce      json
// @Param        id       path  string                    true  "Race UUID"
// @Param        result   body  handler.UpdateResultRequest  true  "Updated result fields"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /races/{id}/result [patch]
func (h *RaceHandler) UpdateResult(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateResultRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateResult(r.Context(), id, req.FinishTime, req.Place, req.Penalty); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// LinkSailor links a race measurement to a sailor
// @Summary      Link race to sailor
// @Description  Links a previously unlinked race measurement to a sailor. Used when the sailor was not known at measurement time.
// @Tags         races
// @Accept       json
// @Param        id       path  string                  true  "Race UUID"
// @Param        request  body  handler.LinkSailorRequest  true  "Sailor ID to link"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /races/{id}/link-sailor [post]
func (h *RaceHandler) LinkSailor(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req LinkSailorRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.LinkSailor(r.Context(), id, req.SailorID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete deletes a race measurement
// @Summary      Delete race
// @Description  Deletes a race measurement by ID.
// @Tags         races
// @Param        id  path  string  true  "Race UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /races/{id} [delete]
func (h *RaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
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