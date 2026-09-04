package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/middleware"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type DocumentHandler struct {
	svc *service.DocumentService
}

func NewDocumentHandler(svc *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

type ApproveRequest struct {
	Comment *string `json:"comment,omitempty" example:"Document looks correct"`
}

// Upload uploads a new document for a sailor
// @Summary      Upload document
// @Description  Uploads a new document (file reference) for a sailor. The file itself should be stored externally; this endpoint stores the reference URL.
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        document  body  model.Document  true  "Document details"
// @Success      201  {object}  model.Document
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /documents [post]
func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var doc model.Document
	if err := decodeJSON(r, &doc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Upload(r.Context(), &doc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

// GetByID returns a document by ID
// @Summary      Get document by ID
// @Description  Returns a single document by UUID.
// @Tags         documents
// @Produce      json
// @Param        id  path  string  true  "Document UUID"
// @Success      200  {object}  model.Document
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string  "document not found"
// @Security     BearerAuth
// @Router       /documents/{id} [get]
func (h *DocumentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	doc, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

// ListBySailor returns all documents for a sailor
// @Summary      List documents by sailor
// @Description  Returns all documents for a given sailor, ordered by creation date descending.
// @Tags         documents
// @Produce      json
// @Param        sailor_id  path  string  true  "Sailor UUID"
// @Success      200  {array}  model.Document
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sailors/{sailor_id}/documents [get]
func (h *DocumentHandler) ListBySailor(w http.ResponseWriter, r *http.Request) {
	sailorID, err := uuid.Parse(r.PathValue("sailor_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid sailor_id")
		return
	}
	docs, err := h.svc.ListBySailor(r.Context(), sailorID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if docs == nil {
		docs = []model.Document{}
	}
	writeJSON(w, http.StatusOK, docs)
}

// Approve approves a document
// @Summary      Approve document
// @Description  Sets the document status to "approved". The checked_by field is set to the authenticated user.
// @Tags         documents
// @Param        id       path  string                true  "Document UUID"
// @Param        request  body  handler.ApproveRequest  false  "Optional comment"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /documents/{id}/approve [post]
func (h *DocumentHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	if err := h.svc.Approve(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Reject rejects a document
// @Summary      Reject document
// @Description  Sets the document status to "rejected" with a required comment explaining the reason.
// @Tags         documents
// @Accept       json
// @Param        id       path  string                true  "Document UUID"
// @Param        request  body  handler.ApproveRequest  true  "Rejection reason (comment is required)"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /documents/{id}/reject [post]
func (h *DocumentHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req ApproveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	comment := ""
	if req.Comment != nil {
		comment = *req.Comment
	}
	if err := h.svc.Reject(r.Context(), id, userID, comment); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RequestRevision requests a document revision
// @Summary      Request document revision
// @Description  Sets the document status to "need_revision" with a comment explaining what needs to be changed.
// @Tags         documents
// @Accept       json
// @Param        id       path  string                  true  "Document UUID"
// @Param        request  body  handler.ApproveRequest  true  "Revision request (comment is required)"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /documents/{id}/request-revision [post]
func (h *DocumentHandler) RequestRevision(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req ApproveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	comment := ""
	if req.Comment != nil {
		comment = *req.Comment
	}
	if err := h.svc.RequestRevision(r.Context(), id, userID, comment); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete deletes a document
// @Summary      Delete document
// @Description  Deletes a document by ID.
// @Tags         documents
// @Param        id  path  string  true  "Document UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /documents/{id} [delete]
func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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