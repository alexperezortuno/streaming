package repository

import (
	"context"

	"github.com/alexperezortuno/streaming/internal/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, req model.RegisterRequest, hashedPassword string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

type VideoRepositoryInterface interface {
	Create(ctx context.Context, video *model.Video) error
	FindByID(ctx context.Context, id string) (*model.Video, error)
	FindAll(ctx context.Context, filter model.VideoFilter) ([]model.Video, int, error)
	Update(ctx context.Context, id string, req model.UpdateVideoRequest) (*model.Video, error)
	UpdateStatus(ctx context.Context, id string, status model.VideoStatus) error
	UpdateDuration(ctx context.Context, id string, duration float64) error
	Delete(ctx context.Context, id string) error
}

type ListRepositoryInterface interface {
	Create(ctx context.Context, name string) (*model.List, error)
	FindAll(ctx context.Context) ([]model.List, error)
	FindByID(ctx context.Context, id string) (*model.List, error)
	Delete(ctx context.Context, id string) error
}
