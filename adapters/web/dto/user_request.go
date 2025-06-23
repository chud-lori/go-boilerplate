package dto

type UserRequest struct {
	Email    string `validate:"required,max=200,min=1" json:"email"`
	Password string `validate:"max=8,min=1" json:"password"`
}
