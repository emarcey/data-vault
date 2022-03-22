package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func HashSha256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:]))
}

func GenUuid() string {
	return uuid.New().String()
}

func GenRandBytes(numBytes int) ([]byte, error) {
	if numBytes <= 0 {
		return nil, fmt.Errorf("Expected positve numBytes. Got %d", numBytes)
	}

	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func EncryptSecret(id, value string, keySize int) (string, *EncryptedSecret, error) {
	key, err := GenRandBytes(keySize)
	if err != nil {
		return "", nil, NewInternalServerErrorFromError("EncryptSecret", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, NewInternalServerErrorFromError("EncryptSecret", err)
	}

	aesGCM, err := cipher.NewGCMWithNonceSize(block, keySize)
	if err != nil {
		return "", nil, NewInternalServerErrorFromError("EncryptSecret", err)
	}

	iv, err := GenRandBytes(aesGCM.NonceSize())
	if err != nil {
		return "", nil, NewInternalServerErrorFromError("EncryptSecret", err)
	}

	ciphertext := aesGCM.Seal(nil, iv, []byte(value), nil)
	return hex.EncodeToString(ciphertext), &EncryptedSecret{
		Id:  id,
		Key: hex.EncodeToString(key),
		Iv:  hex.EncodeToString(iv),
	}, nil
}

func DecryptSecret(ciphertext string, secret *EncryptedSecret) (string, error) {
	if secret == nil {
		return "", NewInternalServerError("EncryptSecret", "Received nil *EncryptedSecret")
	}

	iv, err := hex.DecodeString(secret.Iv)
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}

	key, err := hex.DecodeString(secret.Key)
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}

	value, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}

	aesGCM, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}

	plaintext, err := aesGCM.Open(nil, iv, value, nil)
	if err != nil {
		return "", NewInternalServerErrorFromError("DecryptSecret", err)
	}
	return string(plaintext), nil
}
