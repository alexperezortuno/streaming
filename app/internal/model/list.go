package model

import "time"

type List struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	VideoCount int       `json:"videoCount"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateListRequest struct {
	Name string `json:"name"`
}
