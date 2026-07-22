package handler

import (
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alexperezortuno/streaming/internal/service"
	"github.com/go-chi/chi/v5"
)

type StreamHandler struct {
	svc *service.StreamService
}

func NewStreamHandler(svc *service.StreamService) *StreamHandler {
	return &StreamHandler{svc: svc}
}

func (h *StreamHandler) Serve(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	seg := chi.URLParam(r, "seg")

	if seg == "" || seg == "index.m3u8" {
		playlistPath, err := h.svc.GetHLSPlaylist(videoID)
		if err != nil {
			slog.Warn("playlist not found", "video_id", videoID, "error", err)
			http.Error(w, "playlist not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-mpegURL")
		http.ServeFile(w, r, playlistPath)
		return
	}

	if strings.HasSuffix(seg, ".ts") {
		segPath, err := h.svc.GetHLSSegment(videoID, seg)
		if err != nil {
			slog.Warn("segment not found", "video_id", videoID, "segment", seg, "error", err)
			http.Error(w, "segment not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "video/MP2T")
		http.ServeFile(w, r, segPath)
		return
	}

	http.Error(w, "invalid segment", http.StatusBadRequest)
}

func (h *StreamHandler) ServeStaticMedia(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	filePath := chi.URLParam(r, "*")

	fullPath := filepath.Join(h.svc.GetVideoMediaPath(videoID), filePath)

	w.Header().Set("Content-Type", "video/MP2T")
	http.ServeFile(w, r, fullPath)
}
