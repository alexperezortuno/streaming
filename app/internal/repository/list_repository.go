package repository

import (
	"context"
	"fmt"

	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListRepository struct {
	pool *pgxpool.Pool
}

func NewListRepository(pool *pgxpool.Pool) *ListRepository {
	return &ListRepository{pool: pool}
}

func (r *ListRepository) Create(ctx context.Context, name string) (*model.List, error) {
	l := &model.List{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO lists (name) VALUES ($1) 
		 RETURNING id, name, created_at, updated_at`,
		name,
	).Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create list: %w", err)
	}
	return l, nil
}

func (r *ListRepository) FindAll(ctx context.Context) ([]model.List, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT l.id, l.name, l.created_at, l.updated_at,
		        COALESCE((SELECT COUNT(*) FROM videos v WHERE v.list_id = l.id), 0) AS video_count
		 FROM lists l
		 ORDER BY l.created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("find all lists: %w", err)
	}
	defer rows.Close()

	var lists []model.List
	for rows.Next() {
		var l model.List
		if err := rows.Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt, &l.VideoCount); err != nil {
			return nil, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, l)
	}
	return lists, nil
}

func (r *ListRepository) FindByID(ctx context.Context, id string) (*model.List, error) {
	l := &model.List{}
	err := r.pool.QueryRow(ctx,
		`SELECT l.id, l.name, l.created_at, l.updated_at,
		        COALESCE((SELECT COUNT(*) FROM videos v WHERE v.list_id = l.id), 0) AS video_count
		 FROM lists l WHERE l.id = $1`,
		id,
	).Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt, &l.VideoCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find list by id: %w", err)
	}
	return l, nil
}

func (r *ListRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM lists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete list: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
