package dto

type AuthResponse struct {
	Token string       `json:"id"`
	User  UserResponse `json:"user"`
}
