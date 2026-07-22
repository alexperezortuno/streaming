package mock

import (
	"context"

	"github.com/alexperezortuno/streaming/internal/model"
)

type ListRepository struct {
	CreateFunc   func(ctx context.Context, name string) (*model.List, error)
	FindAllFunc  func(ctx context.Context) ([]model.List, error)
	FindByIDFunc func(ctx context.Context, id string) (*model.List, error)
	DeleteFunc   func(ctx context.Context, id string) error
}

func (m *ListRepository) Create(ctx context.Context, name string) (*model.List, error) {
	return m.CreateFunc(ctx, name)
}
func (m *ListRepository) FindAll(ctx context.Context) ([]model.List, error) {
	return m.FindAllFunc(ctx)
}
func (m *ListRepository) FindByID(ctx context.Context, id string) (*model.List, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *ListRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
