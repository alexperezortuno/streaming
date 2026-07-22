package mock

import (
	"context"

	"github.com/alexperezortuno/streaming/internal/model"
)

type UserRepository struct {
	CreateFunc         func(ctx context.Context, req model.RegisterRequest, hashedPassword string) (*model.User, error)
	FindByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
	FindByIDFunc       func(ctx context.Context, id string) (*model.User, error)
}

func (m *UserRepository) Create(ctx context.Context, req model.RegisterRequest, hashedPassword string) (*model.User, error) {
	return m.CreateFunc(ctx, req, hashedPassword)
}

func (m *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return m.FindByUsernameFunc(ctx, username)
}

func (m *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	return m.FindByIDFunc(ctx, id)
}
