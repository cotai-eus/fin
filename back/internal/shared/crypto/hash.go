package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash is returned when the hash format is invalid
	ErrInvalidHash = errors.New("invalid hash format")

	// ErrIncompatibleVersion is returned when the Argon2 version doesn't match
	ErrIncompatibleVersion = errors.New("incompatible argon2 version")
)

// Argon2Params holds the parameters for Argon2id hashing
type Argon2Params struct {
	// Memory cost in KiB (65536 KiB = 64 MiB)
	Memory uint32

	// Number of iterations (time cost)
	Iterations uint32

	// Number of parallel threads
	Parallelism uint8

	// Length of the derived key in bytes
	KeyLength uint32

	// Length of the random salt in bytes
	SaltLength uint32
}

// DefaultArgon2Params returns secure default parameters for Argon2id
//
// Parameters chosen for balance between security and performance:
// - Memory: 64 MiB (prevents GPU/ASIC attacks)
// - Iterations: 1 (Argon2id is already slow enough with memory-hard function)
// - Parallelism: 4 threads (standard parallelism)
// - KeyLength: 32 bytes (256-bit output, same as AES-256 key)
// - SaltLength: 16 bytes (128-bit salt)
//
// These parameters provide strong protection against:
// - Brute force attacks (memory-hard function)
// - GPU/ASIC attacks (high memory requirement)
// - Rainbow table attacks (unique salt per hash)
// - Timing attacks (constant-time comparison in VerifyPIN)
func DefaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:      64 * 1024, // 64 MiB
		Iterations:  1,
		Parallelism: 4,
		KeyLength:   32,
		SaltLength:  16,
	}
}

// HashPIN hashes a PIN using Argon2id with default parameters.
//
// The returned hash is encoded in the PHC string format:
// $argon2id$v=19$m=65536,t=1,p=4$<base64-salt>$<base64-hash>
//
// This format includes:
// - Algorithm identifier (argon2id)
// - Version (19 = Argon2 version 1.3)
// - Parameters (m=memory, t=iterations, p=parallelism)
// - Base64-encoded salt
// - Base64-encoded hash
//
// Security features:
// - Unique random salt per hash (prevents rainbow tables)
// - Memory-hard function (prevents GPU/ASIC attacks)
// - PHC format allows parameter upgrades without breaking existing hashes
//
// Parameters:
//   - pin: The PIN to hash (typically 4-6 digits)
//
// Returns:
//   - Encoded hash string
//   - Error if salt generation fails
func HashPIN(pin string) (string, error) {
	params := DefaultArgon2Params()

	// Generate a cryptographically secure random salt
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the PIN using Argon2id
	hash := argon2.IDKey(
		[]byte(pin),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Encode in PHC string format
	// Format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil
}

// VerifyPIN verifies a PIN against an Argon2id hash.
//
// This function uses constant-time comparison to prevent timing attacks.
// An attacker cannot learn anything about the correct PIN by measuring
// how long the verification takes.
//
// Parameters:
//   - pin: The PIN to verify
//   - encodedHash: The encoded hash string from HashPIN
//
// Returns:
//   - true if the PIN matches the hash
//   - false if the PIN doesn't match or an error occurs
//   - error if the hash format is invalid
func VerifyPIN(pin string, encodedHash string) (bool, error) {
	// Parse the encoded hash to extract parameters, salt, and hash
	params, salt, hash, err := parseArgon2Hash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the input PIN with the same salt and parameters
	inputHash := argon2.IDKey(
		[]byte(pin),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Use constant-time comparison to prevent timing attacks
	// subtle.ConstantTimeCompare returns 1 if equal, 0 otherwise
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}

	return false, nil
}

// parseArgon2Hash parses an encoded Argon2 hash string and extracts the parameters, salt, and hash.
//
// Expected format: $argon2id$v=19$m=65536,t=1,p=4$<base64-salt>$<base64-hash>
//
// Returns:
//   - Argon2 parameters
//   - Salt (decoded from base64)
//   - Hash (decoded from base64)
//   - Error if parsing fails or format is invalid
func parseArgon2Hash(encodedHash string) (*Argon2Params, []byte, []byte, error) {
	// Split the encoded hash by '$'
	// Expected: ["", "argon2id", "v=19", "m=65536,t=1,p=4", "salt", "hash"]
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Verify algorithm identifier
	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	// Parse version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, ErrInvalidHash
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Parse parameters (m=memory, t=iterations, p=parallelism)
	params := &Argon2Params{}
	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	); err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	// Decode salt from base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	// Decode hash from base64
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	// Set the key length based on the decoded hash length
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
