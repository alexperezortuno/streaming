package service

import (
	"context"

	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/alexperezortuno/streaming/internal/repository"
)

type ListService struct {
	repo repository.ListRepositoryInterface
}

func NewListService(repo repository.ListRepositoryInterface) *ListService {
	return &ListService{repo: repo}
}

func (s *ListService) Create(ctx context.Context, req model.CreateListRequest) (*model.List, error) {
	return s.repo.Create(ctx, req.Name)
}

func (s *ListService) FindAll(ctx context.Context) ([]model.List, error) {
	return s.repo.FindAll(ctx)
}

func (s *ListService) FindByID(ctx context.Context, id string) (*model.List, error) {
	l, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, model.ErrNotFound
	}
	return l, nil
}

func (s *ListService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
