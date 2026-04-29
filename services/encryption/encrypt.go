package encryption

type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
}

type encryptionService struct{}

func NewEncryptionService() *encryptionService {
	return &encryptionService{}
}

func (*encryptionService) Encrypt(plaintext string) (string, error) {
	return "", nil
}
