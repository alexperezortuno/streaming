package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/alexperezortuno/streaming/internal/config"
	"github.com/alexperezortuno/streaming/internal/media"
	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/alexperezortuno/streaming/internal/repository"
	"github.com/google/uuid"
)

type VideoService struct {
	repo       repository.VideoRepositoryInterface
	cfg        *config.Config
	transcoder *media.Transcoder
}

func NewVideoService(repo repository.VideoRepositoryInterface, cfg *config.Config, transcoder *media.Transcoder) *VideoService {
	return &VideoService{
		repo:       repo,
		cfg:        cfg,
		transcoder: transcoder,
	}
}

func (s *VideoService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, req model.CreateVideoRequest) (*model.Video, error) {
	videoID := uuid.New().String()
	videoDir := filepath.Join(s.cfg.MediaPath, videoID)
	hlsDir := filepath.Join(videoDir, "hls")

	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		return nil, fmt.Errorf("create video directory: %w", err)
	}

	inputPath := filepath.Join(hlsDir, fmt.Sprintf("input-%s%s", videoID[:8], filepath.Ext(header.Filename)))

	dst, err := os.Create(inputPath)
	if err != nil {
		os.RemoveAll(videoDir)
		return nil, fmt.Errorf("create input file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		os.RemoveAll(videoDir)
		return nil, fmt.Errorf("write input file: %w", err)
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	video := &model.Video{
		ID:       videoID,
		Name:     req.Name,
		ListID:   req.ListID,
		FilePath: videoDir,
		MIMEType: mimeType,
		Status:   model.VideoStatusUploading,
		Size:     written,
	}

	if err := s.repo.Create(ctx, video); err != nil {
		os.RemoveAll(videoDir)
		return nil, fmt.Errorf("save video record: %w", err)
	}

	s.transcoder.Enqueue(media.TranscodeJob{
		VideoID:   videoID,
		InputPath: inputPath,
		OutputDir: hlsDir,
		OnSuccess: func() {
			if err := s.repo.UpdateStatus(context.Background(), videoID, model.VideoStatusReady); err != nil {
				slog.Error("update video status to ready", "video_id", videoID, "error", err)
			}
		},
		OnError: func(err error) {
			slog.Error("transcode failed", "video_id", videoID, "error", err)
			if err := s.repo.UpdateStatus(context.Background(), videoID, model.VideoStatusError); err != nil {
				slog.Error("update video status to error", "video_id", videoID, "error", err)
			}
		},
	})
	video.Status = model.VideoStatusTranscoding
	if err := s.repo.UpdateStatus(ctx, videoID, model.VideoStatusTranscoding); err != nil {
		slog.Error("update video status to transcoding", "video_id", videoID, "error", err)
	}

	return video, nil
}

func (s *VideoService) FindByID(ctx context.Context, id string) (*model.Video, error) {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, model.ErrNotFound
	}
	return v, nil
}

func (s *VideoService) FindAll(ctx context.Context, filter model.VideoFilter) (*model.PaginatedResponse, error) {
	videos, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	if videos == nil {
		videos = []model.Video{}
	}

	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse{
		Items:      videos,
		Total:      total,
		Page:       filter.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *VideoService) Update(ctx context.Context, id string, req model.UpdateVideoRequest) (*model.Video, error) {
	v, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, model.ErrNotFound
	}
	return v, nil
}

func (s *VideoService) Delete(ctx context.Context, id string) error {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if v == nil {
		return model.ErrNotFound
	}

	if err := os.RemoveAll(v.FilePath); err != nil {
		slog.Error("remove video files", "video_id", id, "path", v.FilePath, "error", err)
	}

	return s.repo.Delete(ctx, id)
}
