package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptAES securely encrypts a string using AES-GCM.
// It expects a 32-byte secret key.
func EncryptAES(plaintext string, secretKey []byte) (string, error) {
	if len(secretKey) != 32 {
		return "", errors.New("encryption key must be exactly 32 bytes")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES reverses the AES-GCM encryption.
func DecryptAES(cryptoText string, secretKey []byte) (string, error) {
	if len(secretKey) != 32 {
		return "", errors.New("encryption key must be exactly 32 bytes")
	}

	// 1. Decode the Base64 string back into bytes
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	// 2. Initialize the cipher block
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	// 3. Initialize GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// 4. Extract the nonce and the actual encrypted message
	// Remember, during encryption, we prepended the nonce to the ciphertext
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 5. Decrypt and authenticate the data
	plaintextBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err // This will fail if the secretKey is wrong or data was tampered with
	}

	return string(plaintextBytes), nil
}
