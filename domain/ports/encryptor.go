package ports

type Encryptor interface {
	HashPassword(password string) (string, error)
	CompareHash(hash, password string) error
}
