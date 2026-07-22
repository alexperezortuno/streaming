package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexperezortuno/streaming/internal/config"
)

type StreamService struct {
	cfg *config.Config
}

func NewStreamService(cfg *config.Config) *StreamService {
	return &StreamService{cfg: cfg}
}

func (s *StreamService) GetHLSPlaylist(videoID string) (string, error) {
	playlistPath := filepath.Join(s.cfg.MediaPath, videoID, "hls", "index.m3u8")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		return "", fmt.Errorf("playlist not found for video %s", videoID)
	}
	return playlistPath, nil
}

func (s *StreamService) GetHLSSegment(videoID, segmentName string) (string, error) {
	segPath := filepath.Join(s.cfg.MediaPath, videoID, "hls", segmentName)
	if _, err := os.Stat(segPath); os.IsNotExist(err) {
		return "", fmt.Errorf("segment not found: %s/%s", videoID, segmentName)
	}
	return segPath, nil
}

func (s *StreamService) GetVideoMediaPath(videoID string) string {
	return filepath.Join(s.cfg.MediaPath, videoID)
}
