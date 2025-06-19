package ports

type TokenManager interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenStr string) (string, error) // return userId
}
