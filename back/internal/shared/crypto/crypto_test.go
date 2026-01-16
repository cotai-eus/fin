package crypto

import (
	"bytes"
	"strings"
	"testing"
)

// TestEncryptDecrypt tests basic encryption and decryption functionality
func TestEncryptDecrypt(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	plaintext := []byte("sensitive card number: 4532123456789012")

	// Encrypt
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Verify ciphertext is different from plaintext
	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext should not equal plaintext")
	}

	// Decrypt
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Verify decrypted data matches original
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted data doesn't match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

// TestEncryptionUniqueness ensures same plaintext produces different ciphertexts
// This is critical for security - nonces must be unique
func TestEncryptionUniqueness(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("same plaintext every time")

	ciphertext1, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	ciphertext2, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Ciphertexts should be different due to unique nonces
	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("Same plaintext produced identical ciphertexts - nonces may not be unique!")
	}

	// But both should decrypt to the same plaintext
	decrypted1, _ := Decrypt(ciphertext1, key)
	decrypted2, _ := Decrypt(ciphertext2, key)

	if !bytes.Equal(decrypted1, plaintext) || !bytes.Equal(decrypted2, plaintext) {
		t.Error("Decrypted data doesn't match original plaintext")
	}
}

// TestTamperDetection verifies that tampered ciphertext is rejected
func TestTamperDetection(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("important data")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Tamper with the ciphertext (flip a bit in the encrypted data portion)
	if len(ciphertext) > nonceSize {
		ciphertext[nonceSize] ^= 0xFF // Flip all bits in first byte of encrypted data
	}

	// Attempt to decrypt tampered ciphertext
	_, err = Decrypt(ciphertext, key)
	if err == nil {
		t.Error("Decryption should fail for tampered ciphertext")
	}

	if err != ErrDecryptionFailed {
		t.Errorf("Expected ErrDecryptionFailed, got: %v", err)
	}
}

// TestInvalidKeySize tests error handling for invalid key sizes
func TestInvalidKeySize(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"too short", 16},
		{"too long", 64},
		{"empty", 0},
	}

	plaintext := []byte("test data")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)

			_, err := Encrypt(plaintext, key)
			if err != ErrInvalidKeySize {
				t.Errorf("Expected ErrInvalidKeySize for %d-byte key, got: %v", tt.keySize, err)
			}

			_, err = Decrypt(plaintext, key)
			if err != ErrInvalidKeySize {
				t.Errorf("Expected ErrInvalidKeySize for decryption, got: %v", err)
			}
		})
	}
}

// TestInvalidCiphertext tests error handling for invalid ciphertext
func TestInvalidCiphertext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")

	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{"too short", []byte("short")},
		{"empty", []byte{}},
		{"only nonce", make([]byte, nonceSize)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext, key)
			if err == nil {
				t.Error("Decryption should fail for invalid ciphertext")
			}
		})
	}
}

// TestEncryptDecryptString tests string encryption convenience functions
func TestEncryptDecryptString(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := "4532123456789012" // Card number

	// Encrypt string
	ciphertext, err := EncryptString(plaintext, key)
	if err != nil {
		t.Fatalf("String encryption failed: %v", err)
	}

	// Decrypt string
	decrypted, err := DecryptString(ciphertext, key)
	if err != nil {
		t.Fatalf("String decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted string doesn't match.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

// TestEmptyPlaintext tests encryption of empty data
func TestEmptyPlaintext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Empty plaintext encryption failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Empty ciphertext decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Empty plaintext roundtrip failed")
	}
}

// TestLargePlaintext tests encryption of larger data (10KB)
func TestLargePlaintext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := bytes.Repeat([]byte("A"), 10*1024) // 10 KB

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Large plaintext encryption failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Large ciphertext decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Large plaintext roundtrip failed")
	}
}

// --- Argon2 PIN Hashing Tests ---

// TestHashVerifyPIN tests basic PIN hashing and verification
func TestHashVerifyPIN(t *testing.T) {
	pin := "1234"

	// Hash the PIN
	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("PIN hashing failed: %v", err)
	}

	// Verify hash format
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("Hash should start with $argon2id$, got: %s", hash)
	}

	// Verify correct PIN
	match, err := VerifyPIN(pin, hash)
	if err != nil {
		t.Fatalf("PIN verification failed: %v", err)
	}
	if !match {
		t.Error("Correct PIN should match the hash")
	}

	// Verify incorrect PIN
	match, err = VerifyPIN("9999", hash)
	if err != nil {
		t.Fatalf("Incorrect PIN verification failed: %v", err)
	}
	if match {
		t.Error("Incorrect PIN should not match the hash")
	}
}

// TestHashUniqueness ensures same PIN produces different hashes
// This is critical - salts must be unique
func TestHashUniqueness(t *testing.T) {
	pin := "1234"

	hash1, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("First hash failed: %v", err)
	}

	hash2, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("Second hash failed: %v", err)
	}

	// Hashes should be different due to unique salts
	if hash1 == hash2 {
		t.Error("Same PIN produced identical hashes - salts may not be unique!")
	}

	// But both should verify with the original PIN
	match1, _ := VerifyPIN(pin, hash1)
	match2, _ := VerifyPIN(pin, hash2)

	if !match1 || !match2 {
		t.Error("Both hashes should verify with the correct PIN")
	}
}

// TestInvalidHashFormat tests error handling for invalid hash formats
func TestInvalidHashFormat(t *testing.T) {
	pin := "1234"

	tests := []struct {
		name string
		hash string
	}{
		{"empty hash", ""},
		{"wrong format", "invalid-hash-format"},
		{"missing parts", "$argon2id$v=19"},
		{"wrong algorithm", "$bcrypt$v=19$m=65536,t=1,p=4$salt$hash"},
		{"malformed params", "$argon2id$v=19$invalid$salt$hash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyPIN(pin, tt.hash)
			if err == nil {
				t.Error("VerifyPIN should fail for invalid hash format")
			}
		})
	}
}

// TestEmptyPIN tests hashing and verification of empty PIN
func TestEmptyPIN(t *testing.T) {
	pin := ""

	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("Empty PIN hashing failed: %v", err)
	}

	match, err := VerifyPIN(pin, hash)
	if err != nil {
		t.Fatalf("Empty PIN verification failed: %v", err)
	}
	if !match {
		t.Error("Empty PIN should match its hash")
	}

	// Different empty input should not match
	match, err = VerifyPIN("0", hash)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}
	if match {
		t.Error("Non-empty PIN should not match empty PIN hash")
	}
}

// TestLongPIN tests hashing of longer PINs
func TestLongPIN(t *testing.T) {
	pin := "123456" // 6-digit PIN

	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("Long PIN hashing failed: %v", err)
	}

	match, err := VerifyPIN(pin, hash)
	if err != nil {
		t.Fatalf("Long PIN verification failed: %v", err)
	}
	if !match {
		t.Error("Long PIN should match its hash")
	}
}

// TestArgon2Parameters verifies the default parameters are as expected
func TestArgon2Parameters(t *testing.T) {
	params := DefaultArgon2Params()

	expectedMemory := uint32(64 * 1024) // 64 MiB
	expectedIterations := uint32(1)
	expectedParallelism := uint8(4)
	expectedKeyLength := uint32(32)
	expectedSaltLength := uint32(16)

	if params.Memory != expectedMemory {
		t.Errorf("Expected memory %d, got %d", expectedMemory, params.Memory)
	}
	if params.Iterations != expectedIterations {
		t.Errorf("Expected iterations %d, got %d", expectedIterations, params.Iterations)
	}
	if params.Parallelism != expectedParallelism {
		t.Errorf("Expected parallelism %d, got %d", expectedParallelism, params.Parallelism)
	}
	if params.KeyLength != expectedKeyLength {
		t.Errorf("Expected key length %d, got %d", expectedKeyLength, params.KeyLength)
	}
	if params.SaltLength != expectedSaltLength {
		t.Errorf("Expected salt length %d, got %d", expectedSaltLength, params.SaltLength)
	}
}

// BenchmarkEncrypt benchmarks encryption performance
func BenchmarkEncrypt(b *testing.B) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("4532123456789012")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(plaintext, key)
	}
}

// BenchmarkDecrypt benchmarks decryption performance
func BenchmarkDecrypt(b *testing.B) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("4532123456789012")
	ciphertext, _ := Encrypt(plaintext, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(ciphertext, key)
	}
}

// BenchmarkHashPIN benchmarks PIN hashing performance
func BenchmarkHashPIN(b *testing.B) {
	pin := "1234"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPIN(pin)
	}
}

// BenchmarkVerifyPIN benchmarks PIN verification performance
func BenchmarkVerifyPIN(b *testing.B) {
	pin := "1234"
	hash, _ := HashPIN(pin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyPIN(pin, hash)
	}
}
