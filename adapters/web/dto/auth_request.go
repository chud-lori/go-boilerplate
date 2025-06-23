package dto

type AuthSignInRequest struct {
	Email    string `json:"email" validate:"required,email,max=200,min=1"`
	Password string `json:"password" validate:"required,max=30,min=8"`
}

type AuthSignUpRequest struct {
	Email           string `json:"email" validate:"required,email,max=200,min=1"`
	Password        string `json:"password" validate:"required,max=30,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
