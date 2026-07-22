package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoRepository struct {
	pool *pgxpool.Pool
}

func NewVideoRepository(pool *pgxpool.Pool) *VideoRepository {
	return &VideoRepository{pool: pool}
}

func (r *VideoRepository) Create(ctx context.Context, video *model.Video) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO videos (id, name, list_id, file_path, mime_type, status, size) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) 
		 RETURNING created_at, updated_at`,
		video.ID, video.Name, video.ListID, video.FilePath,
		video.MIMEType, video.Status, video.Size,
	).Scan(&video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create video: %w", err)
	}
	return nil
}

func (r *VideoRepository) FindByID(ctx context.Context, id string) (*model.Video, error) {
	v := &model.Video{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, list_id, file_path, mime_type, status, duration, size, created_at, updated_at 
		 FROM videos WHERE id = $1`,
		id,
	).Scan(&v.ID, &v.Name, &v.ListID, &v.FilePath, &v.MIMEType, &v.Status,
		&v.Duration, &v.Size, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find video by id: %w", err)
	}
	return v, nil
}

func (r *VideoRepository) FindAll(ctx context.Context, filter model.VideoFilter) ([]model.Video, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.ListID != nil && *filter.ListID != "" {
		conditions = append(conditions, fmt.Sprintf("list_id = $%d", argIdx))
		args = append(args, *filter.ListID)
		argIdx++
	}
	if filter.Search != nil && *filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+*filter.Search+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM videos %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count videos: %w", err)
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit

	args = append(args, filter.Limit, offset)
	query := fmt.Sprintf(
		`SELECT id, name, list_id, file_path, mime_type, status, duration, size, created_at, updated_at 
		 FROM videos %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1,
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("find videos: %w", err)
	}
	defer rows.Close()

	var videos []model.Video
	for rows.Next() {
		var v model.Video
		if err := rows.Scan(&v.ID, &v.Name, &v.ListID, &v.FilePath, &v.MIMEType,
			&v.Status, &v.Duration, &v.Size, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan video: %w", err)
		}
		videos = append(videos, v)
	}
	return videos, total, nil
}

func (r *VideoRepository) Update(ctx context.Context, id string, req model.UpdateVideoRequest) (*model.Video, error) {
	var sets []string
	var args []interface{}
	argIdx := 1

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.ListID != nil {
		sets = append(sets, fmt.Sprintf("list_id = $%d", argIdx))
		args = append(args, *req.ListID)
		argIdx++
	}

	if len(sets) == 0 {
		return r.FindByID(ctx, id)
	}

	sets = append(sets, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE videos SET %s WHERE id = $%d 
		 RETURNING id, name, list_id, file_path, mime_type, status, duration, size, created_at, updated_at`,
		strings.Join(sets, ", "), argIdx,
	)

	v := &model.Video{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&v.ID, &v.Name, &v.ListID, &v.FilePath, &v.MIMEType, &v.Status,
		&v.Duration, &v.Size, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update video: %w", err)
	}
	return v, nil
}

func (r *VideoRepository) UpdateStatus(ctx context.Context, id string, status model.VideoStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE videos SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update video status: %w", err)
	}
	return nil
}

func (r *VideoRepository) UpdateDuration(ctx context.Context, id string, duration float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE videos SET duration = $1, updated_at = NOW() WHERE id = $2`,
		duration, id,
	)
	if err != nil {
		return fmt.Errorf("update video duration: %w", err)
	}
	return nil
}

func (r *VideoRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM videos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete video: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
