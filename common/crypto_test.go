package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			require.NotNil(t, err, "Expected non-nil error")
			require.Nil(t, out, "Expected out to be nil. Got: %v", out)
		})
	}
}

func TestGenRandBytesSuccessCases(t *testing.T) {
	tests := []int{1, 11, 123456}

	for _, given := range tests {
		t.Run(fmt.Sprintf("GenRandBytes - Success - %v", given), func(t *testing.T) {
			out, err := GenRandBytes(given)
			require.Nil(t, err, "Expected err to be nil. Got: %v", err)
			require.Equal(t, len(out), given, "Expected len(out) to be equal to given. Got %d and %d", len(out), given)
		})
	}
}

func TestEncryptSecretErrorCases(t *testing.T) {
	keySizes := []int{-1, 0, 1, 11, 123456}

	for _, given := range keySizes {
		t.Run(fmt.Sprintf("EncryptSecret - Error - %v", given), func(t *testing.T) {
			ciphertext, kv, err := EncryptSecret("id1", "hello", given)
			require.NotNil(t, err, "Expected non-nil error")
			require.Empty(t, ciphertext, "Expected ciphertext to be empty. Got: %v", ciphertext)
			require.Nil(t, kv, "Expected kv to be nil. Got: %v", kv)
		})
	}
}

func TestDecryptSecretErrorCases(t *testing.T) {
	keySizes := []*EncryptedSecret{nil}

	for _, given := range keySizes {
		t.Run(fmt.Sprintf("DecryptSecret - Error - %v", given), func(t *testing.T) {
			plaintext, err := DecryptSecret("id1", given)
			require.NotNil(t, err, "Expected non-nil error")
			require.Empty(t, plaintext, "Expected plaintext to be empty. Got: %v", plaintext)
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	givenValue := "hello there"
	givenId := "id1"

	ciphertext, secret, err := EncryptSecret(givenId, givenValue, KEY_SIZE)
	require.Nil(t, err, "Expected nil error at EncryptSecret. Got: %v", err)
	require.NotEqual(t, ciphertext, givenValue, "Ciphertext was not transformed")
	require.NotNil(t, secret, "Encrypted secret is nil")

	plaintext, err := DecryptSecret(ciphertext, secret)
	require.Nil(t, err, "Expected nil error at DecryptSecret. Got: %v", err)

	require.Equal(t, plaintext, givenValue, "Plaintext, %v, does not equal given, %v", plaintext, givenValue)
}
