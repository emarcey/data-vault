package common

import (
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
