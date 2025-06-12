package dto

import "time"

type UserResponse struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
