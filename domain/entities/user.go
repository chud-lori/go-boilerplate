package entities

import "time"

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Passcode  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
