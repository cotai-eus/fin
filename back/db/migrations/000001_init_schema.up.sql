-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ========================================
-- USERS TABLE
-- ========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kratos_identity_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    cpf VARCHAR(11) UNIQUE,

    -- Balance in cents (avoid floating point issues)
    balance_cents BIGINT DEFAULT 0 CHECK (balance_cents >= 0),

    -- Transfer limits
    daily_transfer_limit_cents BIGINT DEFAULT 100000,   -- R$ 1,000
    monthly_transfer_limit_cents BIGINT DEFAULT 500000, -- R$ 5,000

    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'closed')),
    kyc_status VARCHAR(20) DEFAULT 'pending' CHECK (kyc_status IN ('pending', 'approved', 'rejected')),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- TRANSFERS TABLE
-- ========================================
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    type VARCHAR(20) NOT NULL CHECK (type IN ('pix', 'ted', 'p2p', 'deposit', 'withdrawal')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),

    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    fee_cents BIGINT DEFAULT 0 CHECK (fee_cents >= 0),
    currency VARCHAR(3) DEFAULT 'BRL',

    -- PIX specific fields
    pix_key VARCHAR(255),
    pix_key_type VARCHAR(20) CHECK (pix_key_type IN ('cpf', 'cnpj', 'email', 'phone', 'random')),

    -- TED specific fields
    recipient_name VARCHAR(255),
    recipient_document VARCHAR(14),
    recipient_bank VARCHAR(3),
    recipient_branch VARCHAR(5),
    recipient_account VARCHAR(12),
    recipient_account_type VARCHAR(10) CHECK (recipient_account_type IN ('checking', 'savings')),

    -- P2P specific fields
    recipient_user_id UUID REFERENCES users(id) ON DELETE RESTRICT,

    -- Scheduling and completion
    scheduled_for TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    authentication_code VARCHAR(50),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- CARDS TABLE
-- ========================================
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    type VARCHAR(20) NOT NULL CHECK (type IN ('physical', 'virtual')),
    brand VARCHAR(20) NOT NULL CHECK (brand IN ('visa', 'mastercard', 'elo')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'cancelled', 'lost', 'stolen', 'expired')),

    -- Encrypted fields (BYTEA for binary data)
    card_number_encrypted BYTEA NOT NULL,
    cvv_encrypted BYTEA NOT NULL,
    pin_hash VARCHAR(255),

    -- Unencrypted metadata
    last_four_digits VARCHAR(4) NOT NULL,
    holder_name VARCHAR(255) NOT NULL,
    expiry_month SMALLINT NOT NULL CHECK (expiry_month >= 1 AND expiry_month <= 12),
    expiry_year SMALLINT NOT NULL CHECK (expiry_year >= 2024),

    -- Limits in cents
    daily_limit_cents BIGINT DEFAULT 500000 CHECK (daily_limit_cents >= 0),
    monthly_limit_cents BIGINT DEFAULT 5000000 CHECK (monthly_limit_cents >= 0),
    current_daily_spent_cents BIGINT DEFAULT 0 CHECK (current_daily_spent_cents >= 0),
    current_monthly_spent_cents BIGINT DEFAULT 0 CHECK (current_monthly_spent_cents >= 0),

    -- Security settings
    is_contactless BOOLEAN DEFAULT TRUE,
    is_international BOOLEAN DEFAULT FALSE,
    block_international BOOLEAN DEFAULT FALSE,
    block_online BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    blocked_at TIMESTAMP WITH TIME ZONE
);

-- ========================================
-- CARD TRANSACTIONS TABLE
-- ========================================
CREATE TABLE card_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    merchant_name VARCHAR(255) NOT NULL,
    merchant_category VARCHAR(50),

    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'declined', 'refunded')),
    is_international BOOLEAN DEFAULT FALSE,

    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- BILLS TABLE
-- ========================================
CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    type VARCHAR(20) NOT NULL CHECK (type IN ('bank', 'utility', 'tax', 'other')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'overdue', 'cancelled')),

    barcode VARCHAR(50) NOT NULL UNIQUE,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    fee_cents BIGINT DEFAULT 0 CHECK (fee_cents >= 0),
    final_amount_cents BIGINT NOT NULL CHECK (final_amount_cents > 0),

    recipient_name VARCHAR(255) NOT NULL,
    due_date DATE NOT NULL,
    payment_date TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- BUDGETS TABLE
-- ========================================
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    category VARCHAR(50) NOT NULL,
    period VARCHAR(20) NOT NULL CHECK (period IN ('weekly', 'monthly', 'annual')),

    limit_cents BIGINT NOT NULL CHECK (limit_cents > 0),
    current_spent_cents BIGINT DEFAULT 0 CHECK (current_spent_cents >= 0),

    alert_threshold SMALLINT DEFAULT 75 CHECK (alert_threshold >= 0 AND alert_threshold <= 100),
    alerts_enabled BOOLEAN DEFAULT TRUE,

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT valid_date_range CHECK (end_date > start_date)
);

-- ========================================
-- SUPPORT TICKETS TABLE
-- ========================================
CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('account', 'card', 'transfer', 'bill', 'technical', 'other')),
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'waiting', 'resolved', 'closed')),

    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- TICKET MESSAGES TABLE
-- ========================================
CREATE TABLE ticket_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    message TEXT NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- AUDIT LOGS TABLE (IMMUTABLE)
-- ========================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,

    old_values JSONB,
    new_values JSONB,

    ip_address INET,
    user_agent TEXT,
    request_id VARCHAR(50),

    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failure')),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Make audit_logs immutable (compliance requirement)
CREATE RULE audit_logs_no_update AS ON UPDATE TO audit_logs DO INSTEAD NOTHING;
CREATE RULE audit_logs_no_delete AS ON DELETE TO audit_logs DO INSTEAD NOTHING;

-- ========================================
-- INDEXES FOR PERFORMANCE
-- ========================================

-- Users
CREATE INDEX idx_users_kratos_id ON users(kratos_identity_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_cpf ON users(cpf) WHERE cpf IS NOT NULL;

-- Transfers
CREATE INDEX idx_transfers_user_id ON transfers(user_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_created_at ON transfers(created_at DESC);
CREATE INDEX idx_transfers_scheduled ON transfers(scheduled_for) WHERE scheduled_for IS NOT NULL;

-- Cards
CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_cards_last_four ON cards(last_four_digits);

-- Card Transactions
CREATE INDEX idx_card_txn_card_id ON card_transactions(card_id);
CREATE INDEX idx_card_txn_user_id ON card_transactions(user_id);
CREATE INDEX idx_card_txn_date ON card_transactions(transaction_date DESC);
CREATE INDEX idx_card_txn_category ON card_transactions(merchant_category);

-- Bills
CREATE INDEX idx_bills_user_id ON bills(user_id);
CREATE INDEX idx_bills_barcode ON bills(barcode);
CREATE INDEX idx_bills_status ON bills(status);

-- Budgets
CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_period ON budgets(period);

-- Support Tickets
CREATE INDEX idx_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_tickets_status ON support_tickets(status);
CREATE INDEX idx_tickets_number ON support_tickets(ticket_number);

-- Ticket Messages
CREATE INDEX idx_ticket_msgs_ticket_id ON ticket_messages(ticket_id);

-- Audit Logs
CREATE INDEX idx_audit_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_request_id ON audit_logs(request_id) WHERE request_id IS NOT NULL;
