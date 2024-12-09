package entities

import "time"


type User struct {
    Id          string `json:"id"`
    Email       string `json:"email"`
    Passcode    string `json:"-"`
    Created_at  time.Time `json:"created_at"`
}

