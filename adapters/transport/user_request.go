package transport

type UserRequest struct {
	Email string `validate:"required,max=200,min=1" json:"email"`
}
