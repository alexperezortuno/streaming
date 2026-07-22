package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/alexperezortuno/streaming/internal/service"
	"github.com/go-chi/chi/v5"
)

type ListHandler struct {
	svc *service.ListService
}

func NewListHandler(svc *service.ListService) *ListHandler {
	return &ListHandler{svc: svc}
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	lists, err := h.svc.FindAll(r.Context())
	if err != nil {
		respondError(w, err)
		return
	}
	if lists == nil {
		lists = []model.List{}
	}
	respondJSON(w, http.StatusOK, lists)
}

func (h *ListHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Name == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "name is required"})
		return
	}

	list, err := h.svc.Create(r.Context(), req)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, list)
}

func (h *ListHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	list, err := h.svc.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, list)
}

func (h *ListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "list deleted"})
}
