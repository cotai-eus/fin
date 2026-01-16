# Plano de Integração Backend-Frontend-Database
## LauraTech Financial Platform - Go + gRPC + PostgreSQL + sqlc

**Data**: Janeiro 2026
**Arquiteto**: Claude Sonnet 4.5
**Versão**: 1.0

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Decisões Arquiteturais](#decisões-arquiteturais)
3. [Estrutura do Projeto](#estrutura-do-projeto)
4. [Schema do Banco de Dados](#schema-do-banco-de-dados)
5. [APIs e Endpoints](#apis-e-endpoints)
6. [Segurança e Compliance](#segurança-e-compliance)
7. [Implementação sqlc](#implementação-sqlc)
8. [Fases de Implementação](#fases-de-implementação)
9. [Arquivos Críticos](#arquivos-críticos)
10. [Verificação e Testes](#verificação-e-testes)

---

## Visão Geral

### Estado Atual
- **Frontend**: Next.js 16.1 ✅ (completo, 6 módulos implementados)
- **Backend**: Vazio (apenas Dockerfile)
- **Database**: PostgreSQL do Kratos ✅, database da aplicação ❌
- **Auth**: Ory Kratos ✅ (self-hosted + cloud)
- **Gateway**: APISIX ✅ (configurado)

### Objetivo
Implementar backend Golang que:
- Serve 30+ endpoints REST para o frontend
- Usa gRPC para comunicação interna (futuro microserviços)
- Persiste dados em PostgreSQL com sqlc (type-safe SQL)
- Cumpre PCI-DSS para dados de cartão

### Arquitetura

```
┌─────────────────────────────────┐
│   FRONTEND (Next.js 16.1)       │
│   Server Actions → HTTP/REST    │
└──────────────┬──────────────────┘
               │ HTTP/JSON
    ┌──────────┼──────────┐
    │     APISIX Gateway   │
    │  - Forward Auth      │
    │  - Rate Limiting     │
    │  - CORS              │
    └──────────┬───────────┘
               │
┌──────────────▼──────────────────┐
│  BACKEND (Go Monolith Modular)  │
│  ├─ Chi Router                  │
│  ├─ Middleware (Auth, Audit)    │
│  ├─ Modules (6 domínios)        │
│  └─ sqlc (type-safe queries)    │
└──────────────┬──────────────────┘
               │
    ┌──────────┴──────────┐
    │   PostgreSQL 16     │
    │  - Application DB   │
    │  - Audit Log        │
    └─────────────────────┘
```

---

## Decisões Arquiteturais

### 1. Framework: Chi Router ✅

**Escolhido**: `github.com/go-chi/chi/v5`

**Por quê?**
- ✅ **Stdlib-compatible**: Zero magic, controle total
- ✅ **Explicit errors**: Sem surpresas em JSON binding (crítico em banking)
- ✅ **Context-aware**: Rastreamento de requests (audit trail)
- ✅ **Lightweight**: ~0 overhead, memory previsível
- ✅ **Production-proven**: Uber, Cloudflare em sistemas financeiros

**Rejeitado**: Gin (magic demais), Echo (context customizado), Fiber (não stdlib)

### 2. Database: sqlc ✅

**Escolhido**: `github.com/sqlc-dev/sqlc`

**Por quê?**
- ✅ **Type-safe SQL**: Erros em compile-time, não runtime
- ✅ **Zero reflection**: Performance nativa
- ✅ **Explicit SQL**: Você vê exatamente o que executa (auditável)
- ✅ **Transactions claras**: BEGIN/COMMIT explícitos (ACID compliance)
- ✅ **Fintech-approved**: Usado em Stripe, Coinbase

**Rejeitado**: GORM (magic, N+1 queries), Ent (complexo demais)

### 3. Arquitetura: Monolito Modular ✅

**Por quê?**
- ✅ **Simplicidade**: Single deploy, debugging fácil
- ✅ **Latência**: 0ms entre módulos (vs. gRPC overhead)
- ✅ **Transactions**: Atomicidade cross-domain (ex: transfer + balance update)
- ✅ **Migration path**: Pode separar em microserviços depois

**gRPC**: Usado apenas se/quando separar serviços (ex: Payment Gateway externo)

### 4. Auth: APISIX Header Validation ✅

**Por quê?**
- ✅ **Zero latency**: APISIX já validou sessão, backend apenas lê header
- ✅ **Security**: Header só aceito de APISIX (validar IP)
- ✅ **Padrão**: API Gateway pattern em fintechs

**Fluxo**:
1. Frontend → APISIX (com cookie `ory_kratos_session`)
2. APISIX valida com Kratos → `/sessions/whoami`
3. APISIX injeta header `X-Kratos-Authenticated-Identity-Id: <uuid>`
4. Backend lê header (trusted) → extrai `user_id`

---

## Estrutura do Projeto

```
/home/user/fin/back/
├── cmd/api/main.go                       # Entry point ⭐
│
├── internal/
│   ├── config/                           # Configuração
│   │   └── config.go
│   │
│   ├── server/
│   │   ├── server.go                     # HTTP server setup
│   │   ├── router.go                     # Rotas centrais ⭐
│   │   └── middleware/
│   │       ├── auth.go                   # APISIX header validation ⭐
│   │       ├── logger.go                 # Structured logging
│   │       ├── rate_limit.go             # Per-endpoint rate limit
│   │       ├── request_id.go             # Request tracing
│   │       ├── recovery.go               # Panic recovery
│   │       └── audit.go                  # Audit logging
│   │
│   ├── modules/                          # 6 domínios
│   │   ├── transfers/
│   │   │   ├── handler.go                # HTTP handlers
│   │   │   ├── service.go                # Business logic ⭐
│   │   │   ├── repository.go             # Data access (sqlc)
│   │   │   ├── types.go                  # Domain models
│   │   │   └── validation.go             # Input validation
│   │   ├── cards/
│   │   │   ├── ... (mesma estrutura)
│   │   │   └── encryption.go             # AES-256-GCM ⭐
│   │   ├── bills/
│   │   │   ├── ...
│   │   │   └── barcode.go                # Validação de código de barras
│   │   ├── budgets/
│   │   ├── support/
│   │   └── users/
│   │
│   ├── shared/                           # Utilitários compartilhados
│   │   ├── database/
│   │   │   ├── postgres.go               # Connection pool ⭐
│   │   │   └── transaction.go            # Transaction helpers
│   │   ├── errors/                       # Error codes & types
│   │   ├── response/                     # JSON response helpers
│   │   ├── crypto/
│   │   │   ├── aes.go                    # AES-256-GCM ⭐
│   │   │   └── hash.go                   # Argon2id PIN hashing ⭐
│   │   └── validator/                    # Input validation
│   │
│   └── audit/                            # Audit logging service
│
├── db/
│   ├── migrations/                       # SQL migrations
│   │   ├── 000001_init_schema.up.sql     # Schema completo ⭐
│   │   └── 000001_init_schema.down.sql
│   │
│   ├── queries/                          # sqlc queries ⭐
│   │   ├── users.sql
│   │   ├── transfers.sql
│   │   ├── cards.sql
│   │   ├── bills.sql
│   │   ├── budgets.sql
│   │   └── support.sql
│   │
│   └── sqlc.yaml                         # sqlc config ⭐
│
├── tests/
│   ├── integration/
│   └── e2e/
│
├── go.mod
├── Makefile                              # Build automation ⭐
├── Dockerfile                            # Multi-stage build
└── docker-compose.yml
```

**⭐ = Arquivos críticos para implementação**

---

## Schema do Banco de Dados

### 1. Users (Contas e Saldos)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kratos_identity_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    cpf VARCHAR(11) UNIQUE,  -- CPF brasileiro

    -- Balance em centavos (evita floating point)
    balance_cents BIGINT DEFAULT 0 CHECK (balance_cents >= 0),

    -- Limites
    daily_transfer_limit_cents BIGINT DEFAULT 100000,   -- R$ 1.000
    monthly_transfer_limit_cents BIGINT DEFAULT 500000, -- R$ 5.000

    status VARCHAR(20) DEFAULT 'active',
    kyc_status VARCHAR(20) DEFAULT 'pending',

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_kratos_id ON users(kratos_identity_id);
CREATE INDEX idx_users_cpf ON users(cpf);
```

### 2. Transfers (Transferências PIX/TED/P2P)

```sql
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    type VARCHAR(20) NOT NULL CHECK (type IN ('pix', 'ted', 'p2p', 'deposit', 'withdrawal')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    fee_cents BIGINT DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'BRL',

    -- PIX específico
    pix_key VARCHAR(255),
    pix_key_type VARCHAR(20),

    -- TED específico
    recipient_name VARCHAR(255),
    recipient_document VARCHAR(14),
    recipient_bank VARCHAR(3),
    recipient_branch VARCHAR(5),
    recipient_account VARCHAR(12),
    recipient_account_type VARCHAR(10),

    -- P2P específico
    recipient_user_id UUID REFERENCES users(id),

    -- Agendamento
    scheduled_for TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    authentication_code VARCHAR(50),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_transfers_user_id ON transfers(user_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_created_at ON transfers(created_at DESC);
CREATE INDEX idx_transfers_scheduled ON transfers(scheduled_for) WHERE scheduled_for IS NOT NULL;
```

### 3. Cards (Cartões com Criptografia)

```sql
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    type VARCHAR(20) NOT NULL CHECK (type IN ('physical', 'virtual')),
    brand VARCHAR(20) NOT NULL CHECK (brand IN ('visa', 'mastercard', 'elo')),
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Campos criptografados (BYTEA = binary data)
    card_number_encrypted BYTEA NOT NULL,
    cvv_encrypted BYTEA NOT NULL,
    pin_hash VARCHAR(255),  -- Argon2id hash (irreversível)

    -- Metadata não-criptografada
    last_four_digits VARCHAR(4) NOT NULL,
    holder_name VARCHAR(255) NOT NULL,
    expiry_month SMALLINT NOT NULL,
    expiry_year SMALLINT NOT NULL,

    -- Limites (em centavos)
    daily_limit_cents BIGINT DEFAULT 500000,
    monthly_limit_cents BIGINT DEFAULT 5000000,
    current_daily_spent_cents BIGINT DEFAULT 0,
    current_monthly_spent_cents BIGINT DEFAULT 0,

    -- Configurações de segurança
    is_contactless BOOLEAN DEFAULT TRUE,
    is_international BOOLEAN DEFAULT FALSE,
    block_international BOOLEAN DEFAULT FALSE,
    block_online BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_last_four ON cards(last_four_digits);
```

### 4. Card Transactions (Transações do Cartão)

```sql
CREATE TABLE card_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_id UUID NOT NULL REFERENCES cards(id),
    user_id UUID NOT NULL REFERENCES users(id),

    amount_cents BIGINT NOT NULL,
    merchant_name VARCHAR(255) NOT NULL,
    merchant_category VARCHAR(50),

    status VARCHAR(20) NOT NULL,
    is_international BOOLEAN DEFAULT FALSE,

    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_card_txn_card_id ON card_transactions(card_id);
CREATE INDEX idx_card_txn_date ON card_transactions(transaction_date DESC);
CREATE INDEX idx_card_txn_category ON card_transactions(merchant_category);
```

### 5. Bills (Boletos)

```sql
CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    barcode VARCHAR(50) NOT NULL UNIQUE,
    amount_cents BIGINT NOT NULL,
    fee_cents BIGINT DEFAULT 0,
    final_amount_cents BIGINT NOT NULL,

    recipient_name VARCHAR(255) NOT NULL,
    due_date DATE NOT NULL,
    payment_date TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_bills_user_id ON bills(user_id);
CREATE INDEX idx_bills_barcode ON bills(barcode);
```

### 6. Budgets (Orçamentos)

```sql
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    category VARCHAR(50) NOT NULL,
    period VARCHAR(20) NOT NULL,

    limit_cents BIGINT NOT NULL,
    current_spent_cents BIGINT DEFAULT 0,

    alert_threshold SMALLINT DEFAULT 75,
    alerts_enabled BOOLEAN DEFAULT TRUE,

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_budgets_user_id ON budgets(user_id);
```

### 7. Support Tickets

```sql
CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    category VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',

    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE ticket_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),

    message TEXT NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_ticket_msgs_ticket_id ON ticket_messages(ticket_id);
```

### 8. Audit Logs (IMUTÁVEL - Compliance)

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),

    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,

    old_values JSONB,
    new_values JSONB,

    ip_address INET,
    user_agent TEXT,
    request_id VARCHAR(50),

    status VARCHAR(20) NOT NULL,  -- success, failure

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Tornar imutável (compliance)
CREATE RULE audit_logs_no_update AS ON UPDATE TO audit_logs DO INSTEAD NOTHING;
CREATE RULE audit_logs_no_delete AS ON DELETE TO audit_logs DO INSTEAD NOTHING;

CREATE INDEX idx_audit_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_resource ON audit_logs(resource_type, resource_id);
```

---

## APIs e Endpoints

### Mapeamento Frontend → Backend

#### Transfers (10/hora cada)
- `executePIXTransfer` → `POST /api/transfers/pix`
- `executeTEDTransfer` → `POST /api/transfers/ted`
- `executeP2PTransfer` → `POST /api/transfers/p2p`
- `cancelTransfer` → `POST /api/transfers/{id}/cancel`
- `fetchUserTransfers` → `GET /api/transfers?page={}&limit={}`

#### Cards
- `fetchUserCards` → `GET /api/cards` (100/hora)
- `getCardDetails` → `GET /api/cards/{id}/details` (10/hora ⚠️)
- `createVirtualCard` → `POST /api/cards/virtual` (5/hora)
- `updateCardLimits` → `PATCH /api/cards/{id}/limits` (20/hora)
- `toggleCardStatus` → `POST /api/cards/{id}/block` (20/hora)
- `changeCardPIN` → `POST /api/cards/{id}/pin` (3/hora ⚠️)
- `fetchCardTransactions` → `GET /api/cards/{id}/transactions`

#### Bills
- `validateBarcode` → `POST /api/bills/validate` (20/hora)
- `payBill` → `POST /api/bills/pay` (10/hora)
- `fetchUserBills` → `GET /api/bills`

#### Budgets
- `createBudget` → `POST /api/budgets` (20/hora)
- `getBudgetSummary` → `GET /api/budgets/summary`
- `getCategorySpending` → `GET /api/analytics/category-spending`
- `getSpendingTrends` → `GET /api/analytics/spending-trends`

#### Support
- `createSupportTicket` → `POST /api/support/tickets` (10/hora)
- `fetchUserTickets` → `GET /api/support/tickets`
- `addTicketMessage` → `POST /api/support/tickets/{id}/messages`

### Formato de Resposta Padrão

**Sucesso:**
```json
{
  "data": { /* recurso */ },
  "meta": {
    "request_id": "uuid",
    "timestamp": "2026-01-15T14:30:00Z"
  }
}
```

**Erro:**
```json
{
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Saldo insuficiente",
    "details": {
      "required": 15000,
      "available": 10000
    }
  },
  "meta": {
    "request_id": "uuid",
    "timestamp": "2026-01-15T14:30:00Z"
  }
}
```

**Paginação:**
```json
{
  "data": [ /* items */ ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3,
    "has_more": true
  }
}
```

### Error Codes

```go
const (
    // Authentication (1xxx)
    ErrUnauthorized   = "AUTH_001"
    ErrInvalidSession = "AUTH_002"

    // Validation (2xxx)
    ErrInvalidInput   = "VAL_001"
    ErrInvalidPIXKey  = "VAL_002"
    ErrInvalidBarcode = "VAL_003"

    // Business Logic (3xxx)
    ErrInsufficientBalance   = "BUS_001"
    ErrDailyLimitExceeded    = "BUS_002"
    ErrMonthlyLimitExceeded  = "BUS_003"
    ErrCardBlocked           = "BUS_005"

    // Resource (4xxx)
    ErrNotFound         = "RES_001"
    ErrUserNotFound     = "RES_002"
    ErrCardNotFound     = "RES_003"

    // System (9xxx)
    ErrDatabaseError     = "SYS_001"
    ErrRateLimitExceeded = "SYS_003"
)
```

---

## Segurança e Compliance

### 1. Autenticação (APISIX Header Validation)

**Middleware de Autenticação:**

```go
// internal/server/middleware/auth.go
const HeaderKratosIdentityID = "X-Kratos-Authenticated-Identity-Id"

func AuthMiddleware(trustedProxyIP string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Security: Apenas aceita header de APISIX
            // Produção: validar r.RemoteAddr == trustedProxyIP

            identityID := r.Header.Get(HeaderKratosIdentityID)
            if identityID == "" || !isValidUUID(identityID) {
                http.Error(w, `{"error":{"code":"AUTH_001"}}`, 401)
                return
            }

            // Armazena no context para handlers downstream
            ctx := context.WithValue(r.Context(), "user_id", identityID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Fluxo:**
1. Frontend → APISIX (com cookie)
2. APISIX → Kratos (valida sessão)
3. APISIX → Backend (injeta header `X-Kratos-Authenticated-Identity-Id`)
4. Backend lê header (trusted)

### 2. Criptografia de Dados Sensíveis

#### AES-256-GCM (Números de Cartão, CVV)

```go
// internal/shared/crypto/aes.go
type AESEncryptor struct {
    key []byte // 32 bytes para AES-256
}

func (e *AESEncryptor) Encrypt(plaintext string) ([]byte, error) {
    block, _ := aes.NewCipher(e.key)
    gcm, _ := cipher.NewGCM(block)

    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)

    // GCM fornece criptografia + autenticação
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return ciphertext, nil
}

func (e *AESEncryptor) Decrypt(ciphertext []byte) (string, error) {
    block, _ := aes.NewCipher(e.key)
    gcm, _ := cipher.NewGCM(block)

    nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
    plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

    return string(plaintext), nil
}
```

#### Argon2id (PINs)

```go
// internal/shared/crypto/hash.go
func HashPIN(pin string) (string, error) {
    salt := make([]byte, 16)
    rand.Read(salt)

    // Parâmetros PCI-DSS compliant
    hash := argon2.IDKey(
        []byte(pin),
        salt,
        3,        // iterations
        64*1024,  // memory (64 MB)
        2,        // parallelism
        32,       // key length
    )

    return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=2$%s$%s",
        base64.RawStdEncoding.EncodeToString(salt),
        base64.RawStdEncoding.EncodeToString(hash),
    ), nil
}
```

### 3. Rate Limiting por Endpoint

```go
// Limites por endpoint
func TransferRateLimit() func(http.Handler) http.Handler {
    return RateLimitMiddleware(10, time.Hour) // 10 transfers/hora
}

func CardDetailsRateLimit() func(http.Handler) http.Handler {
    return RateLimitMiddleware(10, time.Hour) // Dados sensíveis
}

func PINChangeRateLimit() func(http.Handler) http.Handler {
    return RateLimitMiddleware(3, time.Hour) // Muito sensível
}
```

### 4. Audit Logging

Toda mutação (POST/PATCH/DELETE) é logada em `audit_logs`:

```go
type AuditEntry struct {
    UserID       string
    Action       string  // "POST /api/transfers/pix"
    ResourceType string  // "TRANSFER"
    ResourceID   string
    OldValues    map[string]interface{}
    NewValues    map[string]interface{}
    IPAddress    string
    UserAgent    string
    RequestID    string
    Status       string  // "success" ou "failure"
}
```

**Logs são imutáveis** (enforced por database rules).

### 5. PCI-DSS Compliance

✅ **Requirement 3**: Protect stored cardholder data
- Card numbers: AES-256-GCM
- CVV: AES-256-GCM
- PIN: Argon2id (irreversível)

✅ **Requirement 4**: Encrypt transmission
- HTTPS/TLS 1.3
- Sem dados sensíveis em URLs

✅ **Requirement 8**: Identify and authenticate
- APISIX + Ory Kratos

✅ **Requirement 10**: Track and monitor
- Audit logs imutáveis
- Rastreamento completo

---

## Implementação sqlc

### 1. Configuração

```yaml
# db/sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "migrations/"
    gen:
      go:
        package: "db"
        out: "../internal/shared/database/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_exact_table_names: false
        json_tags_case_style: "snake"
        overrides:
          - db_type: "pg_catalog.numeric"
            go_type: "int64"  # Valores em centavos
```

### 2. Exemplo de Query

```sql
-- db/queries/transfers.sql

-- name: CreateTransfer :one
INSERT INTO transfers (
    user_id, type, amount_cents, fee_cents,
    pix_key, pix_key_type, status
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTransferByID :one
SELECT * FROM transfers WHERE id = $1 LIMIT 1;

-- name: ListUserTransfers :many
SELECT * FROM transfers
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTransferStatus :one
UPDATE transfers
SET status = $2,
    completed_at = CASE WHEN $2 = 'completed' THEN NOW() END
WHERE id = $1
RETURNING *;

-- name: GetUserForUpdate :one
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: UpdateUserBalance :exec
UPDATE users
SET balance_cents = balance_cents + $2
WHERE id = $1;
```

**Código Gerado:**

```go
// internal/shared/database/sqlc/transfers.sql.go (auto-generated)
type Transfer struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    Type         string
    AmountCents  int64  // ✅ Type-safe de BIGINT
    FeeCents     int64
    PixKey       sql.NullString
    Status       string
    CreatedAt    time.Time
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
    // Implementação gerada automaticamente
}
```

### 3. Pattern de Transação

```go
// internal/modules/transfers/repository.go
func (r *Repository) ExecutePIXTransfer(
    ctx context.Context,
    params CreatePIXTransferParams,
) (*Transfer, error) {
    tx, _ := r.db.BeginTx(ctx, nil)
    defer tx.Rollback()

    qtx := r.queries.WithTx(tx)

    // 1. Verificar saldo (com lock FOR UPDATE)
    user, _ := qtx.GetUserForUpdate(ctx, params.UserID)
    if user.BalanceCents < params.AmountCents {
        return nil, ErrInsufficientBalance
    }

    // 2. Debitar saldo
    _ = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
        ID:     params.UserID,
        Amount: -params.AmountCents,
    })

    // 3. Criar transfer
    transfer, _ := qtx.CreateTransfer(ctx, db.CreateTransferParams{
        UserID:      params.UserID,
        Type:        "pix",
        AmountCents: params.AmountCents,
        // ...
    })

    // 4. Commit
    tx.Commit()
    return &transfer, nil
}
```

### 4. Workflow de Desenvolvimento

```bash
# 1. Criar migration
make migrate-create  # nome: add_transfers_table

# 2. Escrever SQL em migrations/
vim db/migrations/000002_add_transfers_table.up.sql

# 3. Rodar migration
make migrate-up

# 4. Escrever queries em queries/
vim db/queries/transfers.sql

# 5. Gerar código Go
make sqlc

# 6. Usar código gerado
# internal/modules/transfers/repository.go
transfer, err := r.queries.CreateTransfer(ctx, params)
```

---

## Fases de Implementação

### Fase 1: Fundação (Semanas 1-2)
**Objetivo**: Core infrastructure + autenticação

**Tarefas:**
- [ ] Setup projeto Go (go.mod, structure)
- [ ] Dockerfile multi-stage
- [ ] Migrations (000001_init_schema.sql)
- [ ] sqlc config + setup
- [ ] Core middleware (auth, logger, recovery, request_id)
- [ ] Users module (CRUD básico)
- [ ] Health check endpoint

**Deliverable:** Backend rodando em `http://localhost:8080/health`

**Arquivos:**
- `/back/cmd/api/main.go`
- `/back/internal/server/router.go`
- `/back/internal/server/middleware/auth.go`
- `/back/db/migrations/000001_init_schema.up.sql`
- `/back/db/sqlc.yaml`

### Fase 2: Transfers (Semana 3)
**Objetivo**: Funcionalidade completa de transferências

**Tarefas:**
- [ ] Repository (PIX, TED, P2P com transactions)
- [ ] Service (balance checks, limits validation)
- [ ] Handler (HTTP endpoints)
- [ ] Unit tests
- [ ] Integration tests (testcontainers)

**Deliverable:** Todos endpoints de transfer funcionando

**Arquivos:**
- `/back/internal/modules/transfers/handler.go`
- `/back/internal/modules/transfers/service.go`
- `/back/internal/modules/transfers/repository.go`
- `/back/db/queries/transfers.sql`

### Fase 3: Cards (Semana 4)
**Objetivo**: Gestão segura de cartões

**Tarefas:**
- [ ] AES-256-GCM encryption implementation
- [ ] Argon2id PIN hashing
- [ ] Repository (com encrypt/decrypt)
- [ ] Service (limits, security settings)
- [ ] Handler endpoints
- [ ] Security tests

**Deliverable:** CRUD de cartões, dados sensíveis criptografados

**Arquivos:**
- `/back/internal/shared/crypto/aes.go`
- `/back/internal/shared/crypto/hash.go`
- `/back/internal/modules/cards/encryption.go`
- `/back/internal/modules/cards/service.go`

### Fase 4: Bills & Budgets (Semana 5)
**Objetivo**: Boletos e orçamentos

**Tarefas:**
- [ ] Bills: barcode validation logic
- [ ] Bills: repository + service
- [ ] Budgets: repository + service
- [ ] Analytics queries (spending by category)
- [ ] Handler endpoints

**Deliverable:** Bill payment + budget tracking funcionais

**Arquivos:**
- `/back/internal/modules/bills/barcode.go`
- `/back/internal/modules/budgets/service.go`
- `/back/db/queries/bills.sql`
- `/back/db/queries/budgets.sql`

### Fase 5: Support (Semana 6)
**Objetivo**: Sistema de tickets

**Tarefas:**
- [ ] Support tickets repository
- [ ] Ticket messages repository
- [ ] Service layer
- [ ] Handler endpoints

**Deliverable:** Sistema de tickets completo

### Fase 6: Security Hardening (Semana 7)
**Objetivo**: Rate limiting + audit logging

**Tarefas:**
- [ ] Per-endpoint rate limits
- [ ] Audit middleware (gravar em audit_logs)
- [ ] Security testing
- [ ] Penetration testing prep

**Deliverable:** Sistema production-ready em segurança

**Arquivos:**
- `/back/internal/server/middleware/rate_limit.go`
- `/back/internal/server/middleware/audit.go`

### Fase 7: Performance (Semana 8)
**Objetivo**: Otimização

**Tarefas:**
- [ ] Query performance analysis (EXPLAIN ANALYZE + pg_stat_statements)
- [ ] Identify top 10 queries by total time e P99
- [ ] Index optimization (inclui indexes compostos e parciais)
- [ ] Vacuum/Analyze strategy + autovacuum tuning
- [ ] Load testing (k6) com cenários: transfers, cards, bills
- [ ] Concurrency test (picos 200-500 RPS)
- [ ] Connection pooling tuning (pgxpool: min/max conns, idle timeout)
- [ ] JSON payload profiling (response sizes e serialization time)
- [ ] Cache headers para GET list endpoints (ETag/Last-Modified)
- [ ] SLA dashboards (P50/P95/P99) + error rate

**Entregáveis:**
- Relatório de queries lentas + plano de ação
- Migration de índices (se necessário)
- Scripts de carga k6
- Config de pool de conexões validada

**Meta:** P99 < 200ms queries, P99 < 500ms mutations

### Fase 8: Deployment (Semana 9+)
**Objetivo**: Produção

**Tarefas:**
- [ ] Dockerfile multi-stage otimizado
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Monitoring (Prometheus/Grafana)
- [ ] Alerting

**Deliverable:** Deploy production-ready

---

## Arquivos Críticos

### 1. `/back/cmd/api/main.go`
**Razão**: Entry point, dependency injection, server startup

### 2. `/back/internal/server/router.go`
**Razão**: Definição central de rotas, middleware chain

### 3. `/back/internal/server/middleware/auth.go`
**Razão**: Validação de segurança APISIX headers

### 4. `/back/db/migrations/000001_init_schema.up.sql`
**Razão**: Schema completo do database

### 5. `/back/db/sqlc.yaml`
**Razão**: Configuração sqlc para geração de código

### 6. `/back/internal/modules/transfers/service.go`
**Razão**: Lógica de negócio mais complexa (transactions, balances)

### 7. `/back/internal/shared/crypto/aes.go`
**Razão**: Criptografia AES-256-GCM (PCI-DSS compliance)

### 8. `/back/internal/shared/database/postgres.go`
**Razão**: Connection pool, transaction helpers

### 9. `/back/Makefile`
**Razão**: Automação de build, migrations, sqlc generation

---

## Verificação e Testes

### End-to-End Testing

**1. Health Check**
```bash
curl http://localhost:8080/health
# Expected: {"status": "healthy"}
```

**2. Transfer PIX (com auth)**
```bash
curl -X POST http://localhost:9080/api/transfers/pix \
  -H "Cookie: ory_kratos_session=..." \
  -H "Content-Type: application/json" \
  -d '{
    "pix_key": "test@example.com",
    "pix_key_type": "email",
    "amount": 10000,
    "description": "Test transfer"
  }'

# Expected: {"data": {"id": "...", "status": "completed"}}
```

**3. List User Cards**
```bash
curl http://localhost:9080/api/cards \
  -H "Cookie: ory_kratos_session=..."

# Expected: {"data": [...]}
```

### Unit Tests

```bash
make test
# Runs all unit tests with coverage
```

### Integration Tests

```bash
make test-integration
# Runs integration tests with testcontainers (PostgreSQL)
```

### Load Testing

```bash
k6 run tests/load/transfers.js
# Target: P99 < 500ms
```

### Security Testing

- [ ] OWASP ZAP scan
- [ ] SQL injection tests
- [ ] Rate limit validation
- [ ] Encryption verification

---

## Justificativas Técnicas

### Por que Chi sobre Gin/Echo/Fiber?

**Contexto Banking**: Aplicações financeiras requerem controle total sobre erros, boundaries transacionais explícitos, zero comportamento oculto.

Chi fornece:
- ✅ **Zero Magic**: O que você vê é o que roda (crítico para auditorias)
- ✅ **Stdlib-Compatible**: Funciona com todas tools de observabilidade
- ✅ **Explicit Errors**: Você controla cada resposta de erro
- ✅ **Context Propagation**: Suporte nativo para request tracing (audit trail)

### Por que sqlc sobre ORMs (GORM)?

**Contexto Banking**: ORMs escondem queries, podem gerar N+1, tornam transaction boundaries confusos.

sqlc fornece:
- ✅ **Explicit SQL**: Você escreve o SQL exato executado
- ✅ **Type Safety**: Erros em compile-time
- ✅ **Zero Reflection**: Sem overhead de performance
- ✅ **Clear Transactions**: BEGIN/COMMIT explícitos (ACID compliance)

### Por que Argon2id para PINs?

**Contexto Banking**: PCI-DSS requer hashing forte resistente a GPU/ASIC attacks.

Argon2id:
- ✅ **Memory-Hard**: Resistente a brute-force GPU
- ✅ **Configurável**: Ajusta parâmetros baseado em threat model
- ✅ **Side-Channel Resistant**: Operações constant-time

### Por que AES-256-GCM para Card Data?

**Contexto Banking**: PCI-DSS requer encryption at rest. GCM mode fornece confidencialidade + autenticação.

AES-256-GCM:
- ✅ **Authentication**: Previne tampering
- ✅ **NIST-Approved**: Requerido por PCI-DSS
- ✅ **Hardware Acceleration**: Suporte AES-NI
- ✅ **Nonce-Based**: Sem gestão complexa de IV

---

## Próximos Passos

1. ✅ **Revisar este plano** e dar feedback
2. ⏳ **Setup ambiente** (Fase 1, Semana 1)
3. ⏳ **Implementação** seguindo fases
4. ⏳ **Security reviews** a cada fase

**Timeline**: 9 semanas até production-ready backend

**Principais entregas**:
- Semana 2: Health check + Users CRUD
- Semana 3: Transfers funcionais
- Semana 4: Cards com criptografia
- Semana 7: Production-ready security
- Semana 9: Deploy production
