package mock

import (
	"context"

	"github.com/alexperezortuno/streaming/internal/model"
)

type VideoRepository struct {
	CreateFunc         func(ctx context.Context, video *model.Video) error
	FindByIDFunc       func(ctx context.Context, id string) (*model.Video, error)
	FindAllFunc        func(ctx context.Context, filter model.VideoFilter) ([]model.Video, int, error)
	UpdateFunc         func(ctx context.Context, id string, req model.UpdateVideoRequest) (*model.Video, error)
	UpdateStatusFunc   func(ctx context.Context, id string, status model.VideoStatus) error
	UpdateDurationFunc func(ctx context.Context, id string, duration float64) error
	DeleteFunc         func(ctx context.Context, id string) error
}

func (m *VideoRepository) Create(ctx context.Context, video *model.Video) error {
	return m.CreateFunc(ctx, video)
}
func (m *VideoRepository) FindByID(ctx context.Context, id string) (*model.Video, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *VideoRepository) FindAll(ctx context.Context, filter model.VideoFilter) ([]model.Video, int, error) {
	return m.FindAllFunc(ctx, filter)
}
func (m *VideoRepository) Update(ctx context.Context, id string, req model.UpdateVideoRequest) (*model.Video, error) {
	return m.UpdateFunc(ctx, id, req)
}
func (m *VideoRepository) UpdateStatus(ctx context.Context, id string, status model.VideoStatus) error {
	return m.UpdateStatusFunc(ctx, id, status)
}
func (m *VideoRepository) UpdateDuration(ctx context.Context, id string, duration float64) error {
	return m.UpdateDurationFunc(ctx, id, duration)
}
func (m *VideoRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
