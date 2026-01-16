package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	// ErrInvalidKeySize is returned when the encryption key is not 32 bytes
	ErrInvalidKeySize = errors.New("encryption key must be exactly 32 bytes for AES-256")

	// ErrInvalidCiphertext is returned when the ciphertext is too short
	ErrInvalidCiphertext = errors.New("ciphertext too short or invalid")

	// ErrDecryptionFailed is returned when decryption fails (tampered data or wrong key)
	ErrDecryptionFailed = errors.New("decryption failed: data may be tampered or key is incorrect")
)

const (
	// AES-256 requires a 32-byte key
	keySize = 32

	// GCM standard nonce size is 12 bytes
	nonceSize = 12
)

// Encrypt encrypts plaintext using AES-256-GCM with the provided key.
// The returned ciphertext has the format: [nonce (12 bytes)][ciphertext + auth tag]
//
// AES-GCM provides both confidentiality and authenticity (AEAD - Authenticated Encryption with Associated Data).
// Each encryption uses a unique random nonce to ensure the same plaintext produces different ciphertexts.
//
// Security guarantees:
// - Confidentiality: Plaintext is encrypted
// - Authenticity: Any tampering with the ciphertext will be detected during decryption
// - Nonce uniqueness: crypto/rand ensures cryptographically secure random nonces
//
// Parameters:
//   - plaintext: The data to encrypt
//   - key: A 32-byte encryption key (AES-256)
//
// Returns:
//   - Encrypted data with prepended nonce
//   - Error if key size is invalid or encryption fails
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	// Validate key size (AES-256 requires 32 bytes)
	if len(key) != keySize {
		return nil, ErrInvalidKeySize
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode wrapper
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce (12 bytes for GCM)
	// CRITICAL: Never reuse nonces with the same key
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and authenticate
	// gcm.Seal appends the ciphertext and auth tag to the nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts AES-256-GCM encrypted data.
// Expects ciphertext in the format: [nonce (12 bytes)][ciphertext + auth tag]
//
// The authentication tag is verified automatically by GCM.
// If the ciphertext has been tampered with, decryption will fail with ErrDecryptionFailed.
//
// Parameters:
//   - ciphertext: Encrypted data with prepended nonce
//   - key: The same 32-byte encryption key used for encryption
//
// Returns:
//   - Decrypted plaintext
//   - Error if key size is invalid, ciphertext is too short, or authentication fails
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	// Validate key size
	if len(key) != keySize {
		return nil, ErrInvalidKeySize
	}

	// Validate ciphertext length (must have at least nonce + auth tag)
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode wrapper
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce from the beginning of ciphertext
	nonce := ciphertext[:nonceSize]

	// Extract the actual ciphertext (everything after the nonce)
	encryptedData := ciphertext[nonceSize:]

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		// Decryption failed - either tampered data or wrong key
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptString is a convenience wrapper for encrypting strings.
// It converts the string to bytes, encrypts it, and returns the ciphertext.
//
// Use this when encrypting string data like card numbers or CVV codes.
func EncryptString(plaintext string, key []byte) ([]byte, error) {
	return Encrypt([]byte(plaintext), key)
}

// DecryptString is a convenience wrapper for decrypting to strings.
// It decrypts the ciphertext and converts the result to a string.
//
// Use this when decrypting string data like card numbers or CVV codes.
func DecryptString(ciphertext []byte, key []byte) (string, error) {
	plaintext, err := Decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
