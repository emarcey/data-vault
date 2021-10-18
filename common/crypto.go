package common

import (
	"fmt"

	"crypto/sha256"
)

func HashSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("sha256:%s", string(h.Sum(nil)))
}
