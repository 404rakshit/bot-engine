package encryption

import (
	"bot-engine/config"
	"bot-engine/utils"
)

type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type encryptionService struct {
	key []byte
}

// NewEncryptionService initializes the service with your environment secret key.
func NewEncryptionService(botCfg *config.BotSecretsConfig) *encryptionService {

	return &encryptionService{
		key: []byte(botCfg.AESEncryptionKey),
	}
}

// Encrypt delegates the technical cryptographic work to the utils package.
func (s *encryptionService) Encrypt(plaintext string) (string, error) {
	return utils.EncryptAES(plaintext, s.key)
}

func (s *encryptionService) Decrypt(ciphertext string) (string, error) {
	return utils.DecryptAES(ciphertext, s.key)
}
