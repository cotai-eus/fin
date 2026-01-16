# Implementation Plan: Cards Module with AES-256-GCM Encryption and Argon2id PIN Hashing

## Overview
Implement secure card management system with encrypted storage of sensitive data (card numbers, CVV) and hashed PINs following the existing 3-layer architecture pattern.

## Critical Security Requirements
- **AES-256-GCM** for card number and CVV encryption (authenticated encryption with unique nonces)
- **Argon2id** for PIN hashing (memory-hard, GPU-resistant)
- **Encryption at repository layer** (before DB write, after DB read)
- **Never expose plaintext** card numbers/CVV in logs or standard API responses
- **Nonce management**: Prepend 12-byte random nonce to each ciphertext

## Architecture Pattern (from existing codebase)
- 3-layer: Handler → Service → Repository
- SQLC for type-safe SQL queries (queries in `/back/db/queries/`, generated code in `/back/internal/shared/database/sqlc/`)
- Transaction pattern with pessimistic locking (FOR UPDATE)
- Shared response package for consistent API responses
- Dependency injection via constructors

## Files to Create

### 1. Crypto Package (`/back/internal/shared/crypto/`)

**`aes.go`** - AES-256-GCM encryption/decryption
- `Encrypt(plaintext []byte, key []byte) ([]byte, error)` - Returns [nonce][ciphertext+tag]
- `Decrypt(ciphertext []byte, key []byte) ([]byte, error)` - Expects [nonce][ciphertext+tag]
- `EncryptString(plaintext string, key []byte) ([]byte, error)`
- `DecryptString(ciphertext []byte, key []byte) (string, error)`
- **Nonce strategy**: Generate 12 random bytes using crypto/rand, prepend to ciphertext
- **Security**: GCM provides both encryption and authentication (tamper-proof)

**`hash.go`** - Argon2id PIN hashing
- `HashPIN(pin string) (string, error)` - Returns encoded: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
- `VerifyPIN(pin string, hash string) (bool, error)` - Constant-time comparison
- **Parameters**: time=1, memory=64MB, threads=4, keyLen=32
- **Salt**: 16 random bytes per hash using crypto/rand

**`crypto_test.go`** - Security tests
- Encryption uniqueness (same plaintext → different ciphertexts)
- Decryption correctness
- Tamper detection (modified ciphertext fails decryption)
- Hash uniqueness (same PIN → different hashes)
- PIN verification correctness

### 2. Cards Module (`/back/internal/modules/cards/`)

**`types.go`** - Domain models and DTOs
- `Card` - Internal model with decrypted data (CardNumber/CVV have `json:"-"` tags)
- `CardSummary` - For list responses (no sensitive data)
- `CardDetails` - For GET /{id} (masked card number only)
- `CreateCardRequest`, `UpdateLimitsRequest`, `SecuritySettingsRequest`, `SetPINRequest`, `CancelCardRequest`
- **Note**: CreateCardRequest.CardNumber is optional - if empty, auto-generate based on brand
- **Note**: RevealCardResponse deferred to later phase

**`errors.go`** - Module-specific errors
- `ErrCardNotFound`, `ErrCardNotActive`, `ErrCardExpired`
- `ErrInvalidCardNumber`, `ErrInvalidCVV`, `ErrInvalidPIN`, `ErrWeakPIN`
- `ErrDailyLimitExceeded`, `ErrMonthlyLimitExceeded`
- `ErrPINIncorrect`, `ErrInternationalBlocked`, `ErrOnlineBlocked`

**`validation.go`** - Input validation
- `ValidateCardNumber(cardNumber string) error` - Luhn algorithm
- `ValidateCVV(cvv string) error` - 3-4 digits
- `ValidatePIN(pin string) error` - 4-6 digits, no weak patterns (0000, 1234)
- `ValidateExpiryDate(month, year int) error`
- `isWeakPIN(pin string) bool` - Detect repeating/sequential patterns
- `GenerateCardNumber(brand string) (string, error)` - Generate valid card number with Luhn checksum
- `calculateLuhnChecksum(partial string) int` - Helper for Luhn algorithm

**`mapper.go`** - DB ↔ Domain conversion
- `dbCardToCard(dbCard, cardNumber, cvv) *Card` - With decrypted data
- `dbCardToCardSummary(dbCard) *CardSummary` - No sensitive data
- `cardToCardDetails(card) *CardDetails` - With masked number
- `maskCardNumber(cardNumber) string` - Returns "**** **** **** 1234"

**`repository.go`** - Data access with encryption/decryption
- `NewRepository(database *sql.DB, encryptionKey string) *Repository`
- `Create(ctx, params) (*db.Card, error)` - Encrypts card_number and CVV before insert
- `GetByID(ctx, cardID) (*Card, error)` - Decrypts after select
- `ListUserCards(ctx, userID) ([]Card, error)` - Decrypts all cards
- `UpdateStatus(ctx, cardID, status) error`
- `UpdateLimits(ctx, cardID, params) error`
- `UpdateSecuritySettings(ctx, cardID, params) error`
- `UpdatePIN(ctx, cardID, pin) error` - Hashes PIN before update
- `GetForUpdate(ctx, tx, cardID) (*db.Card, error)` - Pessimistic locking
- `UpdateSpentAmounts(ctx, cardID, daily, monthly) error`
- **Encryption flow**: `crypto.EncryptString(cardNumber, r.encryptionKey)` before SQLC insert
- **Decryption flow**: `crypto.DecryptString(dbCard.CardNumberEncrypted, r.encryptionKey)` after SQLC select
- **PIN flow**: `crypto.HashPIN(pin)` before update, `crypto.VerifyPIN(inputPIN, storedHash)` for validation

**`service.go`** - Business logic
- `NewService(repo *Repository, database *sql.DB) *Service`
- `CreateCard(ctx, userID, req) (*Card, error)` - Validates, sets defaults, calls repo
- `GetCardByID(ctx, userID, cardID) (*CardDetails, error)` - Returns masked number
- `ListUserCards(ctx, userID) ([]CardSummary, error)` - No sensitive data
- `BlockCard(ctx, userID, cardID) error` - Status transition
- `UnblockCard(ctx, userID, cardID) error`
- `UpdateLimits(ctx, userID, cardID, req) error`
- `UpdateSecuritySettings(ctx, userID, cardID, req) error`
- `SetPIN(ctx, userID, cardID, req) error` - Validates format, verifies current PIN if changing
- `VerifyPIN(ctx, cardID, pin) (bool, error)` - For transaction authorization
- `CancelCard(ctx, userID, cardID, reason) error` - Soft delete (status = cancelled)
- `ProcessCardTransaction(ctx, cardID, amountCents) error` - Limit checks with locking
- **Transaction pattern**: Uses `executeInTransaction()` for atomic operations
- **Default limits**: Daily R$ 5,000 (500,000 cents), Monthly R$ 50,000 (5,000,000 cents)

**`handler.go`** - HTTP endpoints
- `NewHandler(service *Service) *Handler`
- `POST /api/cards` - CreateCard
- `GET /api/cards` - ListCards (returns CardSummary array)
- `GET /api/cards/{id}` - GetCardDetails (returns CardDetails with masked number)
- `POST /api/cards/{id}/block` - BlockCard
- `POST /api/cards/{id}/unblock` - UnblockCard
- `PATCH /api/cards/{id}/limits` - UpdateLimits
- `PATCH /api/cards/{id}/security` - UpdateSecuritySettings
- `POST /api/cards/{id}/pin` - SetPIN
- `DELETE /api/cards/{id}` - CancelCard (soft delete)
- **Pattern**: Extract user_id from context, decode request, call service, handle errors, return response
- **Note**: Reveal endpoint (POST /cards/{id}/reveal) deferred to later phase with proper 2FA authentication

### 3. SQLC Queries (`/back/db/queries/cards.sql`)

**Queries to implement**:
- `CreateCard :one` - Insert with encrypted fields
- `GetCardByID :one` - Select by ID
- `GetCardForUpdate :one` - Select with FOR UPDATE lock
- `ListUserCards :many` - Select all user cards ordered by created_at DESC
- `ListActiveUserCards :many` - Only active cards
- `UpdateCardStatus :exec` - Update status + blocked_at timestamp
- `UpdateCardLimits :exec` - Update daily/monthly limits
- `UpdateCardSecuritySettings :exec` - Update security flags
- `UpdateCardPIN :exec` - Update pin_hash
- `UpdateCardSpentAmounts :exec` - Update current_daily_spent_cents and current_monthly_spent_cents
- `ResetDailySpent :exec` - Reset to 0 (for cron jobs)
- `ResetMonthlySpent :exec` - Reset to 0 (for cron jobs)
- `DeleteCard :exec` - Soft delete (status = cancelled)
- `CountUserCards :one` - Count total cards
- `CountUserActiveCards :one` - Count active cards

**After creating**: Run `cd /home/user/fin/back && sqlc generate`

## Files to Modify

**`/back/internal/server/server.go`**
- Add `cardsHandler *cards.Handler` to Server struct
- Initialize in `New()`:
  ```go
  cardsRepo := cards.NewRepository(db, cfg.EncryptionKey)
  cardsService := cards.NewService(cardsRepo, db)
  cardsHandler := cards.NewHandler(cardsService)
  ```

**`/back/internal/server/router.go`**
- Add cards routes under `/api/cards`

## Dependencies to Add

Add to `/back/go.mod` (if not present):
```bash
go get golang.org/x/crypto/argon2
```

## Implementation Sequence

### Phase 1: Crypto Foundation (Critical - Test thoroughly)
1. Create `/back/internal/shared/crypto/aes.go`
2. Create `/back/internal/shared/crypto/hash.go`
3. Create `/back/internal/shared/crypto/crypto_test.go`
4. Run: `go test ./internal/shared/crypto/... -v`
5. **Validation**: All crypto tests must pass before proceeding

### Phase 2: Database Layer
6. Create `/back/db/queries/cards.sql`
7. Run: `cd /back && sqlc generate`
8. Create `/back/internal/modules/cards/types.go`
9. Create `/back/internal/modules/cards/errors.go`
10. Create `/back/internal/modules/cards/mapper.go`
11. Create `/back/internal/modules/cards/repository.go`

### Phase 3: Business Logic
12. Create `/back/internal/modules/cards/validation.go`
13. Create `/back/internal/modules/cards/service.go`

### Phase 4: API Layer
14. Create `/back/internal/modules/cards/handler.go`
15. Modify `/back/internal/server/server.go`
16. Modify `/back/internal/server/router.go`

### Phase 5: Testing & Validation
17. Create `/back/internal/modules/cards/repository_test.go` - Repository encryption/decryption tests
18. Create `/back/internal/modules/cards/service_test.go` - Service business logic tests
19. Create `/back/internal/modules/cards/validation_test.go` - Validation and card generation tests
20. Integration tests (create → retrieve → update → delete)
21. Security validation (no plaintext in logs/DB)
22. End-to-end API tests with curl or Postman

## Security Validation Checklist

- [ ] AES encryption produces unique ciphertexts for same plaintext
- [ ] Decryption recovers original plaintext correctly
- [ ] Tampered ciphertext fails decryption
- [ ] Argon2id produces unique hashes for same PIN
- [ ] PIN verification works (correct PIN returns true, incorrect returns false)
- [ ] Card numbers encrypted in database (verify BYTEA field contains non-readable data)
- [ ] CVV encrypted in database
- [ ] PIN stored as hash only (not plaintext)
- [ ] GET /cards returns only masked card numbers
- [ ] GET /cards/{id} returns only masked card number
- [ ] No plaintext card numbers in application logs
- [ ] No plaintext PINs in application logs
- [ ] Daily limit enforcement works
- [ ] Monthly limit enforcement works
- [ ] Luhn validation rejects invalid card numbers
- [ ] Weak PIN patterns rejected (0000, 1234, etc.)

## Key Security Trade-offs

### Implemented in Phase 3:
- ✅ AES-256-GCM encryption (authenticated encryption)
- ✅ Argon2id PIN hashing (memory-hard, GPU-resistant)
- ✅ Nonce prepended to ciphertext (no separate storage)
- ✅ Masked card numbers in standard responses
- ✅ Automatic card number generation with Luhn validation
- ✅ Comprehensive test suite (crypto + integration + security)

### Deferred (recommend for production):
- ⏭️ Reveal endpoint with 2FA authentication (POST /cards/{id}/reveal)
- ⏭️ Rate limiting on PIN attempts (3 failures → 15 min lockout)
- ⏭️ Audit logging for sensitive operations (card creation, PIN changes, reveals)
- ⏭️ Encryption key rotation mechanism
- ⏭️ PCI DSS compliance validation

## Critical Files Summary

1. **`/back/internal/shared/crypto/aes.go`** - AES-256-GCM implementation; must encrypt with unique nonces
2. **`/back/internal/shared/crypto/hash.go`** - Argon2id implementation; must use constant-time comparison
3. **`/back/internal/modules/cards/repository.go`** - Encrypts before DB write, decrypts after read
4. **`/back/internal/modules/cards/service.go`** - Business logic, validation, limit enforcement
5. **`/back/db/queries/cards.sql`** - SQLC queries for type-safe database access

## Testing Strategy

**Unit Tests**:
- Crypto functions (encryption, decryption, hashing, verification)
- Validation functions (Luhn, CVV, PIN format)
- Mapper functions

**Integration Tests**:
- Repository encryption/decryption flow
- Service business logic (limits, status transitions)
- Full CRUD operations

**Security Tests**:
- Verify encrypted data in database
- Verify no plaintext in logs
- Verify tamper detection
- Verify PIN hash security
- Verify generated card numbers pass Luhn validation
- Verify unique card numbers generated per request

## Verification Plan

After implementation:
1. Create a card via API without card_number: Verify auto-generation with valid Luhn checksum
2. Create a card via API with card_number: Verify provided number is validated and used
3. Inspect database: Verify card_number_encrypted and cvv_encrypted are BYTEA (binary, unreadable)
4. Inspect database: Verify pin_hash starts with `$argon2id$`
5. Retrieve card via GET /cards/{id}: Verify only masked number returned
6. Attempt invalid card number: Verify Luhn validation rejects it
7. Set weak PIN (0000): Verify rejection
8. Exceed daily limit: Verify transaction blocked with ErrDailyLimitExceeded
9. Check logs: Verify no plaintext card numbers or PINs logged
10. Run all tests: `go test ./internal/shared/crypto/... ./internal/modules/cards/... -v`
