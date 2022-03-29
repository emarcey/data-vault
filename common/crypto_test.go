package common

import (
	"fmt"
	"testing"
)

func TestHashSha256(t *testing.T) {
	given := "my sha256 test"
	expected := "sha256:b32cd99d1cbe255c0cf0f5e7a139b76fe54e38f281bfb9869a7c8871b497b22d"
	if HashSha256(given) != expected {
		t.Errorf("Given, %v, != expected, %v", given, expected)
	}
}

func TestGenRandBytesErrorCases(t *testing.T) {
	tests := []int{0, -1, -12345}

	for _, given := range tests {
		t.Run(fmt.Sprintf("GenRandBytes - Error - %v", given), func(t *testing.T) {
			out, err := GenRandBytes(given)
			if err == nil || out != nil {
				t.Errorf("Expected empty response and error. Got: %v %v", out, err)
			}
		})
	}
}

func TestGenRandBytesSuccessCases(t *testing.T) {
	tests := []int{1, 11, 123456}

	for _, given := range tests {
		t.Run(fmt.Sprintf("GenRandBytes - Success - %v", given), func(t *testing.T) {
			out, err := GenRandBytes(given)
			if err != nil || len(out) != given {
				t.Errorf("Expected empty error and response of len %d. Got: %v %v", given, out, err)
			}
		})
	}
}

func TestEncryptSecretErrorCases(t *testing.T) {
	keySizes := []int{-1, 0, 1, 11, 123456}

	for _, given := range keySizes {
		t.Run(fmt.Sprintf("EncryptSecret - Error - %v", given), func(t *testing.T) {
			ciphertext, kv, err := EncryptSecret("id1", "hello", given)
			if err == nil || ciphertext != "" || kv != nil {
				t.Errorf("Expected empty response and error. Got: %v %v %v", ciphertext, kv, err)
			}
		})
	}
}

func TestDecryptSecretErrorCases(t *testing.T) {
	keySizes := []*EncryptedSecret{nil}

	for _, given := range keySizes {
		t.Run(fmt.Sprintf("DecryptSecret - Error - %v", given), func(t *testing.T) {
			plaintext, err := DecryptSecret("id1", given)
			if err == nil || plaintext != "" {
				t.Errorf("Expected empty response and error. Got: %v %v", plaintext, err)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	givenValue := "hello there"
	givenId := "id1"

	ciphertext, secret, err := EncryptSecret(givenId, givenValue, KEY_SIZE)
	if err != nil {
		t.Errorf("Unexpected error in EncryptSecret: %v", err)
		return
	}
	if ciphertext == givenValue {
		t.Error("Ciphertext was not transformed")
		return
	}
	if secret == nil {
		t.Error("Encrypted secret is nil")
		return
	}

	plaintext, err := DecryptSecret(ciphertext, secret)
	if err != nil {
		t.Errorf("Unexpected error in DecryptSecret: %v", err)
		return
	}

	if plaintext != givenValue {
		t.Errorf("Plaintext, %v, does not equal given, %v", plaintext, givenValue)
	}
}
