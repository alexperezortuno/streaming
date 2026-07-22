package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/alexperezortuno/streaming/internal/service"
	"github.com/go-chi/chi/v5"
)

type VideoHandler struct {
	svc *service.VideoService
}

func NewVideoHandler(svc *service.VideoService) *VideoHandler {
	return &VideoHandler{svc: svc}
}

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := model.VideoFilter{
		Page:  1,
		Limit: 20,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			filter.Page = p
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}
	if listID := r.URL.Query().Get("listId"); listID != "" {
		filter.ListID = &listID
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
	}

	resp, err := h.svc.FindAll(r.Context(), filter)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	video, err := h.svc.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, video)
}

func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "failed to parse form"})
		return
	}

	name := r.FormValue("name")
	if name == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "name is required"})
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "video file is required"})
		return
	}
	defer file.Close()

	var listID *string
	if lid := r.FormValue("listId"); lid != "" {
		listID = &lid
	}

	req := model.CreateVideoRequest{
		Name:   name,
		ListID: listID,
	}

	video, err := h.svc.Upload(r.Context(), file, header, req)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, video)
}

func (h *VideoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	video, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, video)
}

func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "video deleted"})
}
