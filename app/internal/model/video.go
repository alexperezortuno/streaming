package model

import "time"

type VideoStatus string

const (
	VideoStatusUploading   VideoStatus = "uploading"
	VideoStatusTranscoding VideoStatus = "transcoding"
	VideoStatusReady       VideoStatus = "ready"
	VideoStatusError       VideoStatus = "error"
)

type Video struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	ListID    *string     `json:"listId,omitempty"`
	FilePath  string      `json:"-"`
	MIMEType  string      `json:"mimeType"`
	Status    VideoStatus `json:"status"`
	Duration  *float64    `json:"duration,omitempty"`
	Size      int64       `json:"size"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type CreateVideoRequest struct {
	Name   string  `json:"name"`
	ListID *string `json:"listId,omitempty"`
}

type UpdateVideoRequest struct {
	Name   *string `json:"name,omitempty"`
	ListID *string `json:"listId,omitempty"`
}

type VideoFilter struct {
	ListID *string
	Search *string
	Page   int
	Limit  int
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}
