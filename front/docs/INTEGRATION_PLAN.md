# Plano de Integração Backend-Frontend-Database
## LauraTech Financial Platform

**Data**: Janeiro 2026  
**Status**: Plano de Implementação  
**Versão**: 1.0

---

## 📋 Índice

1. [Visão Geral Arquitetural](#visão-geral-arquitetural)
2. [Arquitetura em Camadas](#arquitetura-em-camadas)
3. [Especificação de APIs REST](#especificação-de-apis-rest)
4. [Modelo de Dados](#modelo-de-dados)
5. [Fluxos de Integração](#fluxos-de-integração)
6. [Padrões de Comunicação](#padrões-de-comunicação)
7. [Persistência de Dados](#persistência-de-dados)
8. [Autenticação e Segurança](#autenticação-e-segurança)
9. [Tratamento de Erros](#tratamento-de-erros)
10. [Caching e Performance](#caching-e-performance)
11. [Monitoramento e Observabilidade](#monitoramento-e-observabilidade)
12. [Roadmap de Implementação](#roadmap-de-implementação)

---

## Visão Geral Arquitetural

### Stack Tecnológico

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND (Next.js 16.1)              │
│  ├─ React Server Components (RSC)                       │
│  ├─ Client Components (Interativo)                      │
│  └─ Server Actions (Chamadas API)                       │
└────────────────────┬────────────────────────────────────┘
                     │
         ┌───────────┴────────────┐
         │  HTTP/REST + WebSocket │
         │  HTTPS (TLS 1.3)       │
         │  Content-Type: JSON    │
         │  gzip Compression      │
         │
┌────────▼────────────────────────────────────────────────┐
│                  BACKEND (Go)                   │
│  ├─ Express/FastAPI Routes                              │
│  ├─ Business Logic Layer                                │
│  ├─ Service Layer (Domain Services)                     │
│  ├─ Repository Layer (Data Access)                      │
│  └─ Middleware (Auth, Validation, Logging)              │
└────────┬────────────────────────────────────────────────┘
         │
    ┌────┴──────┐
    │  Database  │ (Read + Write)
    │  Queries   │
    │
┌───▼──────────────────────────────────────────────────────┐
│                  DATABASE (PostgreSQL)                   │
│  ├─ Transactional Storage (ACID)                         │
│  ├─ Event Log (Audit Trail)                              │
│  ├─ Cache Layer (Redis) [opcional]                       │
│  └─ Backup & Replication                                 │
└──────────────────────────────────────────────────────────┘
```

### Componentes Principais

| Componente | Responsabilidade | Tecnologia |
|-----------|------------------|-----------|
| Frontend Client | Renderizar UI, capturar input | Next.js, React |
| Server Actions | Orquestrar requisições API | Next.js Server Actions |
| Backend API | Processar lógica, persistir dados | Express/Fastify/FastAPI |
| Database | Armazenar e recuperar dados | PostgreSQL |
| Cache | Acelerar reads frequentes | Redis (opcional) |
| Message Queue | Processamento assíncrono | RabbitMQ/Kafka (opcional) |

---

## Arquitetura em Camadas

### Camada de Apresentação (Frontend)

```typescript
// src/app/(dashboard)/transfers/page.tsx (Server Component)
export default async function TransfersPage() {
  // Renderiza a página Server-side
  return <TransfersContainer />;
}

// src/modules/transfers/components/PIXTransferForm.tsx (Client Component)
"use client";
export function PIXTransferForm() {
  const [formData, setFormData] = useState();
  
  const handleSubmit = async (data) => {
    // Chama Server Action
    const result = await executePIXTransfer(data);
  };
}
```

**Responsabilidades**:
- ✅ Renderização de UI
- ✅ Validação client-side
- ✅ Estados locais (form, loading, errors)
- ✅ Chamadas a Server Actions

---

### Camada de Integração (Server Actions)

```typescript
// src/modules/transfers/actions/index.ts
"use server";

export async function executePIXTransfer(input: unknown) {
  try {
    // 1. Verificar sessão
    const session = await requireOrySession();
    
    // 2. Validar dados (Zod)
    const validated = pixTransferSchema.safeParse(input);
    if (!validated.success) return { success: false, error: "..." };
    
    // 3. Chamar API Backend
    const response = await fetch(
      `${BACKEND_URL}/api/transfers/pix`,
      {
        method: "POST",
        headers: getAuthHeaders(session),
        body: JSON.stringify(validated.data),
      }
    );
    
    // 4. Processar resposta
    const transfer = await response.json();
    
    // 5. Revalidar cache
    revalidatePath("/dashboard/transfers");
    
    return { success: true, data: transfer };
  } catch (error) {
    return { success: false, error: error.message };
  }
}
```

**Responsabilidades**:
- ✅ Validação Zod (schema)
- ✅ Autenticação (Ory session)
- ✅ Chamadas HTTP ao backend
- ✅ Transformação de dados
- ✅ Revalidação de cache

---

### Camada de API (Backend)

```
Backend/
├── src/
│   ├── routes/
│   │   ├── transfers.ts         # Endpoints de transferência
│   │   ├── cards.ts              # Endpoints de cartões
│   │   ├── bills.ts              # Endpoints de boletos
│   │   ├── budgets.ts            # Endpoints de orçamentos
│   │   ├── support.ts            # Endpoints de suporte
│   │   └── auth.ts               # Endpoints de autenticação
│   │
│   ├── services/
│   │   ├── TransferService.ts    # Lógica de negócio
│   │   ├── CardService.ts
│   │   ├── BillService.ts
│   │   ├── BudgetService.ts
│   │   └── SupportService.ts
│   │
│   ├── repositories/
│   │   ├── TransferRepository.ts # Data Access
│   │   ├── CardRepository.ts
│   │   ├── BillRepository.ts
│   │   ├── BudgetRepository.ts
│   │   └── UserRepository.ts
│   │
│   ├── models/
│   │   ├── Transfer.ts           # Entity models
│   │   ├── Card.ts
│   │   ├── Bill.ts
│   │   ├── Budget.ts
│   │   └── User.ts
│   │
│   ├── middleware/
│   │   ├── auth.ts              # Validação de token Ory
│   │   ├── validation.ts         # Validação de dados
│   │   ├── errorHandler.ts       # Tratamento de erros
│   │   └── logging.ts            # Logging estruturado
│   │
│   └── utils/
│       ├── validators.ts         # Helper de validação
│       ├── formatters.ts         # Formatação de dados
│       └── constants.ts          # Constantes
```

---

### Camada de Dados (Database)

```sql
-- Schema PostgreSQL
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transfers (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  type VARCHAR(20), -- PIX, TED, P2P
  amount DECIMAL(15,2),
  status VARCHAR(20), -- PENDING, COMPLETED, FAILED
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Índices para performance
CREATE INDEX idx_transfers_user_id ON transfers(user_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_created_at ON transfers(created_at DESC);
```

---

## Especificação de APIs REST

### 1. Módulo de Transferências

#### 1.1 PIX Transfer

```http
POST /api/transfers/pix HTTP/1.1
Host: api.lauraTech.com
Content-Type: application/json
Authorization: Bearer <token>
X-User-Id: <userId>
X-Request-Id: <requestId>

{
  "pixKey": "recipient@example.com",
  "amount": 150.00,
  "description": "Pagamento projeto",
  "scheduledFor": "2026-01-20T10:00:00Z" (opcional)
}
```

**Response (201 Created)**:
```json
{
  "id": "trans_123abc",
  "type": "PIX",
  "amount": 150.00,
  "status": "COMPLETED",
  "pixKey": "recipient@example.com",
  "createdAt": "2026-01-15T14:30:00Z",
  "completedAt": "2026-01-15T14:30:15Z"
}
```

**Error Response (400 Bad Request)**:
```json
{
  "error": "INVALID_PIX_KEY",
  "message": "Chave PIX inválida",
  "code": "PIX_001"
}
```

---

#### 1.2 TED Transfer

```http
POST /api/transfers/ted HTTP/1.1
Content-Type: application/json

{
  "bank": "001", // Banco do Brasil
  "agency": "0001",
  "account": "123456",
  "accountType": "CHECKING",
  "amount": 1000.00,
  "description": "Transferência",
  "recipientName": "João Silva"
}
```

**Response (201 Created)**:
```json
{
  "id": "trans_456def",
  "type": "TED",
  "amount": 1000.00,
  "status": "PENDING",
  "bank": "001",
  "fee": 8.50,
  "estimatedDelivery": "2026-01-16T09:00:00Z"
}
```

---

#### 1.3 P2P Transfer

```http
POST /api/transfers/p2p HTTP/1.1
Content-Type: application/json

{
  "recipientUserId": "user_789",
  "amount": 50.00,
  "description": "Divisão de despesa"
}
```

**Response (201 Created)**:
```json
{
  "id": "trans_789ghi",
  "type": "P2P",
  "amount": 50.00,
  "status": "COMPLETED",
  "recipientUserId": "user_789",
  "transferredAt": "2026-01-15T14:35:00Z"
}
```

---

#### 1.4 List Transfers

```http
GET /api/transfers?page=1&limit=20&status=COMPLETED&type=PIX HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "trans_123",
      "type": "PIX",
      "amount": 150.00,
      "status": "COMPLETED",
      "createdAt": "2026-01-15T14:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "totalPages": 3
  }
}
```

---

#### 1.5 Cancel Transfer

```http
POST /api/transfers/{transferId}/cancel HTTP/1.1
Content-Type: application/json

{
  "reason": "Mudança de ideia"
}
```

**Response (200 OK)**:
```json
{
  "id": "trans_123",
  "status": "CANCELLED",
  "cancelledAt": "2026-01-15T14:35:00Z",
  "reason": "Mudança de ideia"
}
```

---

### 2. Módulo de Cartões

#### 2.1 List User Cards

```http
GET /api/cards HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "card_123",
      "lastFourDigits": "4321",
      "brand": "VISA",
      "type": "PHYSICAL",
      "status": "ACTIVE",
      "expiryDate": "12/27",
      "daily_limit": 5000.00,
      "monthly_limit": 50000.00,
      "spent_today": 1200.00,
      "spent_month": 15000.00
    }
  ]
}
```

---

#### 2.2 Get Card Details (CVV/Número)

```http
GET /api/cards/{cardId}/details HTTP/1.1
X-Verification-Code: <OTP>
```

**Response (200 OK)**:
```json
{
  "id": "card_123",
  "number": "4532 XXXX XXXX 4321",
  "cvv": "XXX",
  "holderName": "JOÃO SILVA"
}
```

---

#### 2.3 Toggle Card Status

```http
POST /api/cards/{cardId}/block HTTP/1.1
Content-Type: application/json

{
  "action": "BLOCK", // ou UNBLOCK
  "reason": "Segurança"
}
```

**Response (200 OK)**:
```json
{
  "id": "card_123",
  "status": "BLOCKED",
  "blockedAt": "2026-01-15T14:40:00Z"
}
```

---

#### 2.4 Create Virtual Card

```http
POST /api/cards/virtual HTTP/1.1
Content-Type: application/json

{
  "limit": 500.00,
  "expiresAt": "2026-02-15T23:59:59Z",
  "description": "Compras online"
}
```

**Response (201 Created)**:
```json
{
  "id": "card_456",
  "number": "4532 1111 1111 4321",
  "cvv": "123",
  "type": "VIRTUAL",
  "limit": 500.00,
  "expiresAt": "2026-02-15T23:59:59Z"
}
```

---

### 3. Módulo de Boletos

#### 3.1 Validate Barcode

```http
POST /api/bills/validate HTTP/1.1
Content-Type: application/json

{
  "barcode": "34191.79001 01017 91020 150008 154500000123456",
  "amount": 154.50 (opcional - para verificação)
}
```

**Response (200 OK)**:
```json
{
  "valid": true,
  "barcode": "34191.79001 01017 91020 150008 154500000123456",
  "type": "BANK", // BANK ou UTILITY
  "amount": 154.50,
  "dueDate": "2026-02-20",
  "recipient": "Empresa XYZ"
}
```

---

#### 3.2 Pay Bill

```http
POST /api/bills/pay HTTP/1.1
Content-Type: application/json

{
  "barcode": "34191.79001 01017 91020 150008 154500000123456",
  "amount": 154.50,
  "paymentDate": "2026-01-20" (opcional)
}
```

**Response (201 Created)**:
```json
{
  "id": "bill_payment_123",
  "status": "COMPLETED",
  "amount": 154.50,
  "paidAt": "2026-01-15T14:45:00Z",
  "receipt": "REC_123456789"
}
```

---

### 4. Módulo de Orçamentos

#### 4.1 Create Budget

```http
POST /api/budgets HTTP/1.1
Content-Type: application/json

{
  "category": "FOOD",
  "limit": 500.00,
  "period": "MONTHLY", // WEEKLY, MONTHLY, ANNUAL
  "alertThresholds": [50, 75, 90],
  "startDate": "2026-01-01",
  "endDate": "2026-01-31"
}
```

**Response (201 Created)**:
```json
{
  "id": "budget_123",
  "category": "FOOD",
  "limit": 500.00,
  "spent": 0.00,
  "percentage": 0,
  "period": "MONTHLY",
  "createdAt": "2026-01-15T14:50:00Z"
}
```

---

#### 4.2 Get Budget Summary

```http
GET /api/budgets/summary HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "totalBudget": 5000.00,
  "totalSpent": 1500.00,
  "percentageSpent": 30,
  "budgets": [
    {
      "id": "budget_123",
      "category": "FOOD",
      "limit": 500.00,
      "spent": 250.00,
      "status": "SAFE" // SAFE, WARNING, DANGER
    }
  ]
}
```

---

#### 4.3 Get Category Spending (Analytics)

```http
GET /api/analytics/category-spending?period=MONTHLY HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "data": [
    {
      "category": "FOOD",
      "spent": 500.00,
      "percentageOfTotal": 35,
      "transactionCount": 25
    },
    {
      "category": "TRANSPORT",
      "spent": 250.00,
      "percentageOfTotal": 18,
      "transactionCount": 12
    }
  ],
  "total": 1500.00,
  "period": "2026-01-01_2026-01-31"
}
```

---

#### 4.4 Get Spending Trends

```http
GET /api/analytics/spending-trends?days=30 HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "data": [
    {
      "date": "2026-01-01",
      "amount": 50.00
    },
    {
      "date": "2026-01-02",
      "amount": 120.50
    }
  ]
}
```

---

### 5. Módulo de Suporte

#### 5.1 Create Support Ticket

```http
POST /api/support/tickets HTTP/1.1
Content-Type: application/json

{
  "category": "CARD", // ACCOUNT, CARD, TRANSFER, BILL, TECHNICAL, OTHER
  "priority": "HIGH", // LOW, MEDIUM, HIGH, URGENT
  "subject": "Cartão bloqueado",
  "description": "Meu cartão foi bloqueado sem motivo",
  "attachments": ["file_123"] (opcional)
}
```

**Response (201 Created)**:
```json
{
  "id": "ticket_123",
  "number": "TKT-2026-00123",
  "status": "OPEN",
  "category": "CARD",
  "priority": "HIGH",
  "createdAt": "2026-01-15T15:00:00Z"
}
```

---

#### 5.2 Get User Tickets

```http
GET /api/support/tickets?page=1&status=OPEN&priority=HIGH HTTP/1.1
```

**Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "ticket_123",
      "number": "TKT-2026-00123",
      "status": "OPEN",
      "category": "CARD",
      "priority": "HIGH",
      "subject": "Cartão bloqueado",
      "createdAt": "2026-01-15T15:00:00Z",
      "updatedAt": "2026-01-15T15:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "total": 5
  }
}
```

---

#### 5.3 Add Ticket Message

```http
POST /api/support/tickets/{ticketId}/messages HTTP/1.1
Content-Type: application/json

{
  "message": "Incluindo informações adicionais...",
  "attachments": [] (opcional)
}
```

**Response (201 Created)**:
```json
{
  "id": "msg_123",
  "ticketId": "ticket_123",
  "message": "Incluindo informações adicionais...",
  "authorType": "USER", // USER ou STAFF
  "createdAt": "2026-01-15T15:05:00Z"
}
```

---

## Modelo de Dados

### Entity Relationship Diagram (ERD)

```
┌─────────────┐         ┌──────────────┐
│   users     │────────→│  transfers   │
├─────────────┤         ├──────────────┤
│ id (PK)     │         │ id (PK)      │
│ email       │         │ user_id (FK) │
│ ory_id      │         │ type         │
│ name        │         │ amount       │
│ created_at  │         │ status       │
└─────────────┘         │ created_at   │
      │                 └──────────────┘
      │
      ├──────────→ ┌──────────────┐
      │            │   cards      │
      │            ├──────────────┤
      │            │ id (PK)      │
      │            │ user_id (FK) │
      │            │ last_4       │
      │            │ brand        │
      │            │ status       │
      └───────────→│ created_at   │
                   └──────────────┘
                         │
                         └────→ ┌──────────────┐
                                │ transactions │
                                ├──────────────┤
                                │ id (PK)      │
                                │ card_id (FK) │
                                │ amount       │
                                │ merchant     │
                                │ created_at   │
                                └──────────────┘

┌─────────────┐         ┌──────────────┐
│   users     │────────→│   budgets    │
└─────────────┘         ├──────────────┤
                        │ id (PK)      │
                        │ user_id (FK) │
                        │ category     │
                        │ limit        │
                        │ spent        │
                        │ period       │
                        └──────────────┘

┌─────────────┐         ┌──────────────┐
│   users     │────────→│   tickets    │
└─────────────┘         ├──────────────┤
                        │ id (PK)      │
                        │ user_id (FK) │
                        │ status       │
                        │ category     │
                        │ priority     │
                        │ created_at   │
                        └──────────────┘
                              │
                              └────→ ┌──────────────┐
                                     │ ticket_msgs  │
                                     ├──────────────┤
                                     │ id (PK)      │
                                     │ ticket_id(FK)│
                                     │ message      │
                                     │ author_type  │
                                     │ created_at   │
                                     └──────────────┘
```

### Schema SQL Detalhado

```sql
-- Users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ory_id VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255),
  phone VARCHAR(20),
  cpf VARCHAR(11) UNIQUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- Transfers
CREATE TABLE transfers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  type VARCHAR(20) NOT NULL, -- PIX, TED, P2P
  amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
  status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, COMPLETED, FAILED, CANCELLED
  description TEXT,
  
  -- PIX fields
  pix_key VARCHAR(255),
  
  -- TED fields
  bank_code VARCHAR(3),
  agency VARCHAR(5),
  account VARCHAR(12),
  account_type VARCHAR(10),
  recipient_name VARCHAR(255),
  
  -- P2P fields
  recipient_user_id UUID REFERENCES users(id),
  
  -- Scheduling
  scheduled_for TIMESTAMP,
  
  -- Timestamps
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  completed_at TIMESTAMP,
  
  CONSTRAINT valid_amount CHECK (amount > 0)
);

-- Cards
CREATE TABLE cards (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  type VARCHAR(20) NOT NULL, -- PHYSICAL, VIRTUAL
  brand VARCHAR(20) NOT NULL, -- VISA, MASTERCARD, ELO
  last_four VARCHAR(4) NOT NULL,
  full_number_encrypted BYTEA NOT NULL,
  cvv_encrypted BYTEA NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, BLOCKED, CANCELLED, LOST, STOLEN
  expiry_date VARCHAR(5) NOT NULL, -- MM/YY
  daily_limit DECIMAL(15,2) NOT NULL,
  monthly_limit DECIMAL(15,2) NOT NULL,
  spent_today DECIMAL(15,2) DEFAULT 0,
  spent_month DECIMAL(15,2) DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP,
  blocked_at TIMESTAMP
);

-- Transactions (despesas do cartão)
CREATE TABLE card_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  card_id UUID NOT NULL REFERENCES cards(id),
  user_id UUID NOT NULL REFERENCES users(id),
  amount DECIMAL(15,2) NOT NULL,
  merchant VARCHAR(255) NOT NULL,
  category VARCHAR(50),
  status VARCHAR(20) DEFAULT 'COMPLETED',
  created_at TIMESTAMP DEFAULT NOW()
);

-- Bills
CREATE TABLE bills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  barcode VARCHAR(50) NOT NULL UNIQUE,
  amount DECIMAL(15,2) NOT NULL,
  type VARCHAR(20), -- BANK, UTILITY, OTHER
  due_date DATE,
  recipient VARCHAR(255),
  status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PAID, OVERDUE
  paid_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Budgets
CREATE TABLE budgets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  category VARCHAR(50) NOT NULL, -- FOOD, TRANSPORT, LEISURE, etc
  limit_amount DECIMAL(15,2) NOT NULL,
  spent_amount DECIMAL(15,2) DEFAULT 0,
  period VARCHAR(20) NOT NULL, -- WEEKLY, MONTHLY, ANNUAL
  alert_thresholds JSONB DEFAULT '[50,75,90]',
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Support Tickets
CREATE TABLE support_tickets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  ticket_number VARCHAR(20) UNIQUE NOT NULL,
  category VARCHAR(50) NOT NULL, -- ACCOUNT, CARD, TRANSFER, BILL, TECHNICAL
  priority VARCHAR(20) NOT NULL, -- LOW, MEDIUM, HIGH, URGENT
  subject VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  status VARCHAR(20) DEFAULT 'OPEN', -- OPEN, IN_PROGRESS, WAITING, RESOLVED, CLOSED
  assigned_to UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  resolved_at TIMESTAMP
);

-- Ticket Messages
CREATE TABLE ticket_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_id UUID NOT NULL REFERENCES support_tickets(id),
  author_id UUID NOT NULL REFERENCES users(id),
  author_type VARCHAR(20) NOT NULL, -- USER, STAFF
  message TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Audit Log (rastreamento de todas as operações)
CREATE TABLE audit_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  action VARCHAR(100) NOT NULL, -- TRANSFER_CREATED, CARD_BLOCKED, etc
  resource_type VARCHAR(50) NOT NULL, -- TRANSFER, CARD, BUDGET, etc
  resource_id UUID NOT NULL,
  old_values JSONB,
  new_values JSONB,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_transfers_user_id ON transfers(user_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_created_at ON transfers(created_at DESC);
CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_transactions_card_id ON card_transactions(card_id);
CREATE INDEX idx_transactions_user_id ON card_transactions(user_id);
CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_period ON budgets(period);
CREATE INDEX idx_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_tickets_status ON support_tickets(status);
CREATE INDEX idx_audit_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_created_at ON audit_log(created_at DESC);
```

---

## Fluxos de Integração

### Fluxo 1: Transferência PIX

```
┌─────────────┐
│  User Form  │
│(PIX form)   │
└──────┬──────┘
       │ 1. Valida locally (regex, amount)
       │
       ▼
┌─────────────────────────┐
│  Server Action          │
│ executePIXTransfer()    │
└──────┬──────────────────┘
       │ 2. Validação Zod + Session
       │
       ▼
┌──────────────────────────────┐
│  Backend /api/transfers/pix  │
│  POST                        │
└──────┬───────────────────────┘
       │ 3. Validação + Business Logic
       │ 4. Busca usuário no DB
       │ 5. Valida PIX key
       │
       ▼
┌──────────────────────┐
│  External Service    │
│  (Banco/Ory)         │
│  Executa transferência
└──────┬───────────────┘
       │ 6. Confirmação
       │
       ▼
┌──────────────────────────┐
│  Database Transaction    │
│  INSERT INTO transfers   │
│  UPDATE user balance     │
└──────┬───────────────────┘
       │ 7. Sucesso
       │
       ▼
┌──────────────────────────┐
│  Response to Frontend    │
│  { success, data }       │
└──────┬───────────────────┘
       │ 8. Revalidate cache
       │ 9. Show success toast
       │
       ▼
┌─────────────────┐
│  User sees UI   │
│  updated (✓)    │
└─────────────────┘
```

### Fluxo 2: Pagamento de Boleto

```
┌────────────────────────┐
│  User Scans Barcode    │
│  (Camera)              │
└────────┬───────────────┘
         │ 1. Envia barcode
         │
         ▼
┌──────────────────────────────┐
│  Server Action               │
│  validateBarcode()           │
└────────┬─────────────────────┘
         │ 2. Zod validation
         │
         ▼
┌───────────────────────────────┐
│  Backend /api/bills/validate  │
│  POST                         │
└────────┬──────────────────────┘
         │ 3. Regex + Database lookup
         │
         ▼
┌──────────────────────────────┐
│  Response: Bill Details      │
│  { amount, dueDate, ... }    │
└────────┬─────────────────────┘
         │ 4. Exibe para confirmação
         │
         ▼
┌────────────────────────────┐
│  User Confirms Payment     │
│  (Click Pagar)             │
└────────┬───────────────────┘
         │ 5. Envia confirmação
         │
         ▼
┌────────────────────────────────┐
│  Server Action                 │
│  payBill()                     │
└────────┬───────────────────────┘
         │ 6. Validação
         │
         ▼
┌──────────────────────────────┐
│  Backend /api/bills/pay      │
│  POST                        │
└────────┬─────────────────────┘
         │ 7. DB Transaction:
         │    - INSERT bill_payments
         │    - UPDATE user balance
         │    - Log audit
         │
         ▼
┌──────────────────────┐
│  Response: Receipt   │
│  { receipt ID, ... } │
└────────┬─────────────┘
         │ 8. Download/Share
         │
         ▼
┌──────────────────────┐
│  User Success Page   │
└──────────────────────┘
```

### Fluxo 3: Sincronização de Orçamentos

```
Frontend Request:
GET /api/budgets/summary

Backend Processing:
1. requireOrySession() → get user_id
2. SELECT budgets WHERE user_id = ?
3. SELECT SUM(amount) FROM transactions 
   WHERE user_id = ? AND date >= ?
4. Calculate: percentage = spent/limit * 100
5. Determine: status = SAFE|WARNING|DANGER

Response:
{
  totalBudget: 5000,
  totalSpent: 1500,
  budgets: [
    {
      category: FOOD,
      spent: 250,
      percentage: 50,
      status: SAFE
    }
  ]
}

Frontend Update:
- Atualiza estado local
- Re-renderiza gráficos
- Cach é revalidado via revalidatePath()
```

---

## Padrões de Comunicação

### Request/Response Pattern

```typescript
// Frontend (Server Action)
const handleTransfer = async (data) => {
  try {
    const result = await executePIXTransfer(data);
    if (result.success) {
      toast.success("Transferência realizada!");
    } else {
      toast.error(result.error);
    }
  } catch (error) {
    toast.error("Erro inesperado");
  }
};

// Backend API
app.post('/api/transfers/pix', async (req, res) => {
  try {
    // 1. Validate
    if (!req.body.pixKey) {
      return res.status(400).json({
        error: 'INVALID_REQUEST',
        message: 'Chave PIX obrigatória'
      });
    }
    
    // 2. Execute
    const transfer = await transferService.executePIX(req.body);
    
    // 3. Return
    return res.status(201).json(transfer);
  } catch (error) {
    return res.status(500).json({
      error: 'INTERNAL_SERVER_ERROR',
      message: error.message
    });
  }
});
```

### Error Handling Pattern

```typescript
// Frontend
type ActionResult<T> = 
  | { success: true; data: T }
  | { success: false; error: string };

// Backend
type ApiResponse<T> = {
  data?: T;
  error?: string;
  code?: string;
  message?: string;
};

// Mapping
export const errorCodes = {
  PIX_INVALID_KEY: { status: 400, message: 'Chave PIX inválida' },
  INSUFFICIENT_BALANCE: { status: 402, message: 'Saldo insuficiente' },
  TRANSFER_LIMIT_EXCEEDED: { status: 402, message: 'Limite de transferência excedido' },
  CARD_BLOCKED: { status: 403, message: 'Cartão bloqueado' },
  USER_NOT_FOUND: { status: 404, message: 'Usuário não encontrado' },
  DATABASE_ERROR: { status: 500, message: 'Erro ao acessar o banco de dados' },
};
```

### Pagination Pattern

```typescript
// Request
GET /api/transfers?page=2&limit=20&sort=created_at&order=DESC

// Response
{
  data: Transfer[],
  pagination: {
    page: number,
    limit: number,
    total: number,
    totalPages: number,
    hasMore: boolean,
    nextCursor?: string
  }
}
```

---

## Persistência de Dados

### Estratégia de Armazenamento

#### Primary Database (PostgreSQL)

```
Dados Críticos:
- Transferências, cartões, contas
- Saldos, limites, transações
- Dados do usuário (PII)

Estratégia:
- ACID transactions
- Replicação para backup
- WAL (Write-Ahead Logging)
```

#### Cache Layer (Redis) - Opcional

```typescript
// Cache keys
const cacheKeys = {
  userBalance: (userId) => `balance:${userId}`,
  userCards: (userId) => `cards:${userId}`,
  transferHistory: (userId, page) => `transfers:${userId}:${page}`,
  budgetSummary: (userId) => `budget:${userId}`,
};

// TTL (Time-To-Live)
const cacheTTL = {
  userBalance: 5 * 60, // 5 minutos
  userCards: 30 * 60, // 30 minutos
  transferHistory: 10 * 60, // 10 minutos
  budgetSummary: 15 * 60, // 15 minutos
};

// Invalidação
await redis.del(cacheKeys.userBalance(userId)); // ao fazer transferência
await redis.del(cacheKeys.userCards(userId)); // ao bloquear cartão
```

### Transações Garantidas

```typescript
// Usar transações para operações críticas
async executeTransfer(userId, transferData) {
  const client = await pool.connect();
  
  try {
    await client.query('BEGIN');
    
    // 1. Verificar saldo
    const user = await client.query(
      'SELECT balance FROM users WHERE id = $1 FOR UPDATE',
      [userId]
    );
    
    if (user.rows[0].balance < transferData.amount) {
      throw new Error('INSUFFICIENT_BALANCE');
    }
    
    // 2. Debitar conta
    await client.query(
      'UPDATE users SET balance = balance - $1 WHERE id = $2',
      [transferData.amount, userId]
    );
    
    // 3. Creditar conta
    await client.query(
      'UPDATE users SET balance = balance + $1 WHERE id = $2',
      [transferData.amount, recipientId]
    );
    
    // 4. Log da transação
    const result = await client.query(
      'INSERT INTO transfers (user_id, type, amount, status) VALUES ($1, $2, $3, $4) RETURNING *',
      [userId, 'PIX', transferData.amount, 'COMPLETED']
    );
    
    await client.query('COMMIT');
    return result.rows[0];
  } catch (error) {
    await client.query('ROLLBACK');
    throw error;
  } finally {
    client.release();
  }
}
```

### Backup e Disaster Recovery

```bash
#!/bin/bash
# Backup diário PostgreSQL
BACKUP_DIR="/backups/postgres"
DB_NAME="lauraTech"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

pg_dump -h localhost -U postgres $DB_NAME | \
  gzip > "$BACKUP_DIR/backup_$TIMESTAMP.sql.gz"

# Retenção de 30 dias
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

# Upload para S3
aws s3 cp "$BACKUP_DIR/backup_$TIMESTAMP.sql.gz" \
  s3://lauraTech-backups/postgres/
```

---

## Autenticação e Segurança

### Integração Ory Kratos

```typescript
// Frontend: Verificar sessão
const session = await getOrySession();
if (!session.active) {
  redirect('/auth/login');
}

// Backend: Validar token Ory
import { FrontendApi } from "@ory/client";

const verifyOrySession = async (sessionToken: string) => {
  const client = new FrontendApi({
    basePath: process.env.ORY_SDK_URL,
  });
  
  try {
    const session = await client.toSession({
      sessionToken: sessionToken
    });
    return session.data;
  } catch (error) {
    throw new Error('INVALID_SESSION');
  }
};
```

### Headers de Autenticação

```typescript
// Server Action
const headers = {
  'Content-Type': 'application/json',
  'X-User-Id': session.identity.id,
  'X-Request-Id': crypto.randomUUID(),
  // Cookie com session token é enviado automaticamente pelo navegador
};

// Backend verifica
middleware.auth = (req, res, next) => {
  const userId = req.headers['x-user-id'];
  const requestId = req.headers['x-request-id'];
  
  if (!userId || !requestId) {
    return res.status(401).json({ error: 'UNAUTHORIZED' });
  }
  
  // Valida contra Ory
  const session = verifyOrySession(req.cookies.session);
  if (!session) {
    return res.status(401).json({ error: 'INVALID_SESSION' });
  }
  
  req.user = session;
  next();
};
```

### Criptografia de Dados Sensíveis

```typescript
import crypto from 'crypto';

const encryptCardData = (cardNumber: string) => {
  const cipher = crypto.createCipher('aes-256-cbc', process.env.ENCRYPTION_KEY);
  let encrypted = cipher.update(cardNumber, 'utf8', 'hex');
  encrypted += cipher.final('hex');
  return encrypted;
};

const decryptCardData = (encrypted: string) => {
  const decipher = crypto.createDecipher('aes-256-cbc', process.env.ENCRYPTION_KEY);
  let decrypted = decipher.update(encrypted, 'hex', 'utf8');
  decrypted += decipher.final('utf8');
  return decrypted;
};

// No banco de dados
INSERT INTO cards (full_number_encrypted, cvv_encrypted)
VALUES (encryptCardData(cardNumber), encryptCardData(cvv));
```

### Rate Limiting

```typescript
import rateLimit from 'express-rate-limit';

const apiLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutos
  max: 100, // 100 requests por IP
  message: 'Muitas requisições, tente novamente depois',
});

const transferLimiter = rateLimit({
  windowMs: 60 * 60 * 1000, // 1 hora
  max: 10, // máximo 10 transferências por hora
  skip: (req) => req.user.role === 'admin', // admins não têm limite
  keyGenerator: (req) => req.user.id, // por usuário, não por IP
});

app.post('/api/transfers/pix', transferLimiter, (req, res) => {
  // ...
});
```

---

## Tratamento de Erros

### Estratégia de Erro Comum

```typescript
// Definir enum de erros
enum ErrorCode {
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  UNAUTHORIZED = 'UNAUTHORIZED',
  FORBIDDEN = 'FORBIDDEN',
  NOT_FOUND = 'NOT_FOUND',
  CONFLICT = 'CONFLICT',
  INTERNAL_ERROR = 'INTERNAL_ERROR',
  EXTERNAL_SERVICE_ERROR = 'EXTERNAL_SERVICE_ERROR',
}

class ApiError extends Error {
  constructor(
    public code: ErrorCode,
    public status: number,
    public message: string,
    public details?: any
  ) {
    super(message);
  }
}

// Middleware de erro
app.use((err: ApiError, req, res, next) => {
  const status = err.status || 500;
  
  res.status(status).json({
    error: err.code,
    message: err.message,
    ...(process.env.NODE_ENV === 'development' && { details: err.details }),
  });
  
  // Log
  logger.error({
    error: err.code,
    message: err.message,
    userId: req.user?.id,
    path: req.path,
  });
});
```

### Erros Específicos por Módulo

```typescript
// Transferências
- INSUFFICIENT_BALANCE: 402
- TRANSFER_LIMIT_EXCEEDED: 402
- INVALID_PIX_KEY: 400
- ACCOUNT_NOT_FOUND: 404
- TRANSFER_SCHEDULED: 201 (sucesso)

// Cartões
- CARD_BLOCKED: 403
- CARD_EXPIRED: 403
- INVALID_CVV: 400
- LIMIT_EXCEEDED: 402
- CARD_NOT_FOUND: 404

// Boletos
- INVALID_BARCODE: 400
- BILL_NOT_FOUND: 404
- BILL_OVERDUE: 402
- BILL_ALREADY_PAID: 409

// Suporte
- TICKET_NOT_FOUND: 404
- TICKET_CLOSED: 409
- INVALID_CATEGORY: 400
```

---

## Caching e Performance

### Estratégia de Cache

```typescript
// Cache inverso (HTTP)
app.use((req, res, next) => {
  // GET requests para listagens
  if (req.method === 'GET' && req.path.startsWith('/api/')) {
    res.set('Cache-Control', 'private, max-age=300'); // 5 min
  }
  next();
});

// Application cache (Redis/Memory)
const cachedGetUserCards = async (userId: string) => {
  const cached = await redis.get(`cards:${userId}`);
  if (cached) return JSON.parse(cached);
  
  const cards = await cardService.getUserCards(userId);
  await redis.set(`cards:${userId}`, JSON.stringify(cards), 'EX', 1800);
  return cards;
};

// Query optimization
SELECT id, last_four, brand, status 
FROM cards 
WHERE user_id = $1
LIMIT 10; -- Nunca retornar tudo
```

### Lazy Loading

```typescript
// Frontend
export async function CardsList({ userId }: { userId: string }) {
  const cards = await getCardDetails(userId); // fetch inicial
  
  return (
    <div>
      {cards.map(card => (
        <Suspense key={card.id} fallback={<CardSkeleton />}>
          <CardItem cardId={card.id} />
        </Suspense>
      ))}
    </div>
  );
}

// Carrega cada card em paralelo
```

### CDN para Conteúdo Estático

```typescript
// next.config.ts
export default {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'cdn.lauraTech.com',
      },
    ],
    unoptimized: false, // usar Next.js image optimization
  },
};
```

---

## Monitoramento e Observabilidade

### Estrutura de Logging

```typescript
import winston from 'winston';

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.File({ 
      filename: 'error.log', 
      level: 'error' 
    }),
    new winston.transports.File({ filename: 'combined.log' }),
  ],
});

// Uso
logger.info('User login', {
  userId: user.id,
  timestamp: new Date(),
  ip: req.ip,
});

logger.error('Transfer failed', {
  userId: user.id,
  transferId: transfer.id,
  reason: error.message,
  stack: error.stack,
});
```

### Métricas

```typescript
import prometheus from 'prom-client';

// Métricas personalizadas
const transferCounter = new prometheus.Counter({
  name: 'transfers_total',
  help: 'Total de transferências processadas',
  labelNames: ['type', 'status'],
});

const responseTime = new prometheus.Histogram({
  name: 'http_request_duration_seconds',
  help: 'Tempo de resposta HTTP em segundos',
  labelNames: ['method', 'route', 'status'],
});

// Middleware
app.use((req, res, next) => {
  const start = Date.now();
  
  res.on('finish', () => {
    const duration = (Date.now() - start) / 1000;
    responseTime
      .labels(req.method, req.route?.path || req.path, res.statusCode)
      .observe(duration);
  });
  
  next();
});

// Endpoint metrics
app.get('/metrics', (req, res) => {
  res.set('Content-Type', prometheus.register.contentType);
  res.end(prometheus.register.metrics());
});
```

### Alertas

```yaml
# AlertManager config
groups:
  - name: api_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "Alta taxa de erros na API"
      
      - alert: DatabaseDown
        expr: up{job="postgres"} == 0
        for: 1m
        annotations:
          summary: "Banco de dados offline"
      
      - alert: SlowAPIs
        expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
        for: 5m
        annotations:
          summary: "APIs respondendo lentamente"
```

---

## Roadmap de Implementação

### Fase 1: Fundação (Semanas 1-4)

- [ ] Setup backend (Express/Fastify)
- [ ] Conectar PostgreSQL
- [ ] Implementar autenticação Ory
- [ ] Criar middleware base (auth, validation, errors)
- [ ] Implementar 3 endpoints de transferência (PIX, TED, P2P)

**Entregáveis**:
- Backend rodando localmente
- Endpoints testáveis via Postman
- Database schema criado
- Testes unitários básicos

---

### Fase 2: Módulos Principais (Semanas 5-8)

- [ ] Endpoints de cartões (CRUD, block/unblock)
- [ ] Endpoints de boletos (validate, pay)
- [ ] Endpoints de orçamentos (create, list, summary)
- [ ] Endpoints de suporte (tickets, messages)
- [ ] Integração com frontend (testar fluxos)

**Entregáveis**:
- Todos os endpoints implementados
- Integração completa com frontend
- Testes de integração
- Documentação OpenAPI/Swagger

---

### Fase 3: Otimização (Semanas 9-10)

- [ ] Implementar Redis cache
- [ ] Rate limiting
- [ ] Compressão gzip
- [ ] Índices de database
- [ ] Paginação eficiente

**Entregáveis**:
- Performance benchmarks
- Load testing report
- Cache hit ratio > 70%

---

### Fase 4: Segurança (Semana 11)

- [ ] Criptografia de dados sensíveis
- [ ] Auditoria (audit log)
- [ ] CORS configurado
- [ ] HTTPS/TLS
- [ ] Rate limiting por endpoint

**Entregáveis**:
- Security audit checklist ✓
- Certificado SSL instalado
- Audit logs em produção

---

### Fase 5: Observabilidade (Semana 12)

- [ ] Logging estruturado
- [ ] Métricas (Prometheus)
- [ ] Alertas (AlertManager)
- [ ] Rastreamento (Jaeger/Zipkin)
- [ ] Dashboards (Grafana)

**Entregáveis**:
- Dashboard Grafana funcional
- Alertas configurados
- Logs centralizados (ELK/Splunk)

---

### Fase 6: Produção (Semana 13+)

- [ ] Docker/Kubernetes
- [ ] CI/CD (GitHub Actions)
- [ ] Staging environment
- [ ] Health checks
- [ ] Disaster recovery plan

**Entregáveis**:
- Pipeline CI/CD automático
- Containers pronto para prod
- Runbook de operações

---

## Checklist de Implementação

### Pré-requisitos
- [ ] PostgreSQL 14+
- [ ] Node.js 18+ / Python 3.10+
- [ ] Redis (opcional)
- [ ] Ory Kratos rodando
- [ ] Variáveis de ambiente configuradas

### Backend
- [ ] Framework (Express/FastAPI) setup
- [ ] Autenticação Ory integrada
- [ ] Validação de input (Zod/Pydantic)
- [ ] Error handling middleware
- [ ] Logging estruturado

### Database
- [ ] Schema SQL criado
- [ ] Índices implementados
- [ ] Migrations automáticas
- [ ] Backup strategy definida

### APIs
- [ ] 30+ endpoints implementados
- [ ] Documentação OpenAPI
- [ ] Paginação suportada
- [ ] Tratamento de erros robusto

### Testes
- [ ] Testes unitários (>80% coverage)
- [ ] Testes de integração
- [ ] Testes E2E (fluxos críticos)
- [ ] Performance tests

### Segurança
- [ ] HTTPS/TLS
- [ ] Rate limiting
- [ ] CORS configurado
- [ ] Dados sensíveis criptografados
- [ ] Auditoria implementada

### Performance
- [ ] Cache implementado
- [ ] Queries otimizadas
- [ ] Compressão gzip
- [ ] Load testing passed

### Deployment
- [ ] Docker container
- [ ] Health checks
- [ ] Metrics/monitoring
- [ ] Alertas configurados
- [ ] Documentation completa

---

## Resumo

Este plano de integração fornece:

✅ **Arquitetura em 3 camadas** clara e escalável  
✅ **30+ endpoints REST** documentados e especificados  
✅ **Modelo de dados completo** com schema SQL  
✅ **Padrões de comunicação** bem definidos  
✅ **Segurança robusta** com Ory Kratos  
✅ **Estratégia de cache** para performance  
✅ **Observabilidade completa** com logging e métricas  
✅ **Roadmap prático** de 13 semanas  

**Próximos passos**: Selecionar tecnologia backend (Express/FastAPI), criar repositório, iniciar Fase 1.

