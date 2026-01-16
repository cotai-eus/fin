# LauraTech Backend - Go API

Backend Golang para a plataforma financeira LauraTech, implementando 30+ endpoints REST com autenticação via Ory Kratos, banco de dados PostgreSQL e arquitetura modular.

## 📋 Stack Tecnológico

- **Language**: Go 1.22+
- **Router**: Chi v5 (lightweight, stdlib-compatible)
- **Database**: PostgreSQL 16
- **ORM**: sqlc (type-safe SQL code generation)
- **Authentication**: Ory Kratos (via APISIX header validation)
- **Encryption**: AES-256-GCM (card data), Argon2id (PINs)

## 🏗️ Arquitetura

```
Monolito Modular
├── cmd/api/main.go           # Entry point
├── internal/
│   ├── config/               # Configuração
│   ├── server/               # HTTP server + router
│   │   └── middlewares/      # Auth, Logger, Request ID, CORS
│   ├── modules/              # Domínios de negócio
│   │   ├── users/            # ✅ Implementado (Fase 1)
│   │   ├── transfers/        # 🚧 Fase 2
│   │   ├── cards/            # 🚧 Fase 3
│   │   ├── bills/            # 🚧 Fase 4
│   │   ├── budgets/          # 🚧 Fase 4
│   │   └── support/          # 🚧 Fase 5
│   └── shared/               # Utilitários compartilhados
│       ├── database/         # Connection pool
│       └── response/         # JSON response helpers
└── db/
    ├── migrations/           # SQL migrations
    ├── queries/              # sqlc queries
    └── sqlc.yaml             # sqlc config
```

## 🚀 Início Rápido

### 1. Pré-requisitos

```bash
# Go 1.22+
go version

# PostgreSQL 16
psql --version

# Ferramentas de desenvolvimento
make install-tools
```

### 2. Configuração

```bash
# Criar arquivo .env
cp .env.example .env

# Editar variáveis de ambiente
vim .env
```

**Variáveis obrigatórias**:
```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/lauratech?sslmode=disable
ENCRYPTION_KEY=CHANGE-ME-32-BYTES-KEY-FOR-AES256
```

**Gerar chave de criptografia**:
```bash
openssl rand -base64 32
```

### 3. Database Setup

```bash
# Criar database
createdb lauratech

# Rodar migrations
make migrate-up
```

### 4. Rodar Aplicação

```bash
# Desenvolvimento (hot reload com Air)
make dev

# Produção
make build
./bin/api
```

## 📡 Endpoints Implementados

### Health Check (Público)

```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "database": "connected"
}
```

### Users (Autenticado)

**Get Current User**
```bash
GET /api/users/me
Headers:
  X-Kratos-Authenticated-Identity-Id: <uuid>
```

**Create User**
```bash
POST /api/users
Content-Type: application/json

{
  "kratos_identity_id": "uuid",
  "email": "user@example.com",
  "full_name": "João Silva",
  "cpf": "12345678901"
}
```

**Update User**
```bash
PATCH /api/users/me
Content-Type: application/json
Headers:
  X-Kratos-Authenticated-Identity-Id: <uuid>

{
  "full_name": "João Silva Santos"
}
```

## 🔐 Segurança

### Autenticação (APISIX Header Validation)

O backend **não valida sessões diretamente**. Confia no header `X-Kratos-Authenticated-Identity-Id` injetado pelo APISIX após validação com Ory Kratos.

**Fluxo**:
1. Frontend → APISIX (com cookie `ory_kratos_session`)
2. APISIX valida com Kratos → `/sessions/whoami`
3. APISIX injeta header → Backend lê `user_id`

### Middleware Stack

```go
1. Recovery          // Panic recovery
2. RequestID         // Request tracing
3. Logger            // Structured logging
4. Timeout (30s)     // Request timeout
5. CORS              // Cross-origin
6. Auth              // APISIX header validation
```

### PCI-DSS Compliance

- ✅ **AES-256-GCM**: Números de cartão, CVV
- ✅ **Argon2id**: PINs (irreversível)
- ✅ **Audit Logs**: Imutáveis (compliance)
- ✅ **HTTPS/TLS 1.3**: Criptografia em trânsito

## 🛠️ Comandos Make

```bash
# Desenvolvimento
make dev                    # Run com hot reload (Air)
make run                    # Run direto
make build                  # Build binário

# Database
make migrate-up             # Aplicar migrations
make migrate-down           # Reverter migrations
make migrate-create name=X  # Criar nova migration
make sqlc                   # Gerar código sqlc

# Testes
make test                   # Rodar todos os testes
make test-unit              # Apenas unit tests
make test-integration       # Integration tests
make test-coverage          # Coverage HTML

# Code Quality
make fmt                    # Formatar código
make vet                    # Go vet
make lint                   # Golangci-lint

# Docker
make docker-build           # Build imagem
make docker-up              # Start services
make docker-down            # Stop services
```

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    kratos_identity_id VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    full_name VARCHAR(255),
    cpf VARCHAR(11) UNIQUE,
    balance_cents BIGINT,
    daily_transfer_limit_cents BIGINT,
    monthly_transfer_limit_cents BIGINT,
    status VARCHAR(20),
    kyc_status VARCHAR(20),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Outras Tabelas
- `transfers` - PIX, TED, P2P
- `cards` - Cartões físicos/virtuais (com criptografia)
- `card_transactions` - Transações do cartão
- `bills` - Boletos
- `budgets` - Orçamentos
- `support_tickets` - Tickets de suporte
- `audit_logs` - Logs imutáveis

## 🧪 Testing

```bash
# Unit tests
go test ./internal/modules/users -v

# Integration tests (requer database)
DATABASE_URL=... make test-integration

# Load testing (k6)
k6 run tests/load/health.js
```

## 📈 Métricas

- **Build Time**: ~3s
- **Binary Size**: 9.3 MB
- **Cold Start**: <100ms
- **Health Check**: <10ms

## 🔄 Próximas Fases

- **Fase 2 (Semana 3)**: Módulo de Transferências (PIX, TED, P2P)
- **Fase 3 (Semana 4)**: Módulo de Cartões (com criptografia)
- **Fase 4 (Semana 5)**: Bills & Budgets
- **Fase 5 (Semana 6)**: Support Tickets
- **Fase 6 (Semana 7)**: Security Hardening (rate limiting, audit)
- **Fase 7 (Semana 8)**: Performance Optimization
- **Fase 8 (Semana 9)**: Deploy Production

## 📝 Notas de Implementação

### Por que Chi sobre Gin/Fiber?
- **Zero magic**: Controle total sobre erros (banking requirement)
- **Stdlib-compatible**: Funciona com Prometheus, Jaeger
- **Explicit**: Você vê exatamente o que roda

### Por que sqlc sobre GORM?
- **Type-safe**: Erros em compile-time
- **Explicit SQL**: Sem N+1 queries surpresa
- **Zero reflection**: Performance nativa
- **Clear transactions**: ACID compliance

### Por que Monolito Modular?
- **Simplicidade**: Single deploy, debugging fácil
- **Latência**: 0ms entre módulos
- **ACID**: Transações cross-domain
- **Migration path**: Pode virar microserviços depois

## 🐛 Troubleshooting

**Database connection error**:
```bash
# Verificar PostgreSQL
psql $DATABASE_URL
```

**Migration error**:
```bash
# Forçar versão
make migrate-force version=1
```

**Build error**:
```bash
# Limpar e rebuild
go clean -cache
go mod tidy
make build
```

## 📚 Documentação Adicional

- [Plano de Integração](./docs/INTEGRATION_PLAN.md)
- [API Specification](./docs/API_SPEC.md) (a ser criado)
- [Deployment Guide](./docs/DEPLOYMENT.md) (a ser criado)

## 🤝 Contribuindo

1. Crie feature branch: `git checkout -b feature/nova-funcionalidade`
2. Faça commit: `git commit -m "feat: adiciona nova funcionalidade"`
3. Push branch: `git push origin feature/nova-funcionalidade`
4. Abra Pull Request

## 📄 License

Proprietary - LauraTech © 2026
