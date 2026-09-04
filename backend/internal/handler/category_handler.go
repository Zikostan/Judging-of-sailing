package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List returns all categories
// @Summary      List categories
// @Description  Returns all scoring categories/divisions.
// @Tags         categories
// @Produce      json
// @Success      200  {array}  model.Category
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if cats == nil {
		cats = []model.Category{}
	}
	writeJSON(w, http.StatusOK, cats)
}

// Create creates a new category
// @Summary      Create category
// @Description  Creates a new scoring category (age range, gender, boat class).
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body  model.Category  true  "Category details (id is ignored)"
// @Success      201  {object}  model.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c model.Category
	if err := decodeJSON(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Create(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

// GetByID returns a category by ID
// @Summary      Get category by ID
// @Description  Returns a single category by its UUID.
// @Tags         categories
// @Produce      json
// @Param        id  path  string  true  "Category UUID"
// @Success      200  {object}  model.Category
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "category not found"
// @Security     BearerAuth
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// Update updates a category
// @Summary      Update category
// @Description  Updates an existing category. The ID in the path overrides any ID in the body.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path  string          true  "Category UUID"
// @Param        category  body  model.Category  true  "Updated category fields"
// @Success      200  {object}  model.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var c model.Category
	if err := decodeJSON(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c.ID = id
	if err := h.svc.Update(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// Delete deletes a category
// @Summary      Delete category
// @Description  Deletes a category by ID.
// @Tags         categories
// @Param        id  path  string  true  "Category UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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