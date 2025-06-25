-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    identity_number TEXT NOT NULL UNIQUE,
    is_delinquent BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Loans Table
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(20, 2) NOT NULL,
    interest_rate NUMERIC(5, 2) NOT NULL,
    term_months INTEGER NOT NULL,
    total_payment NUMERIC(20, 2) NOT NULL,
    outstanding_balance NUMERIC(20, 2) NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'paid', 'defaulted')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Billing Schedules Table
CREATE TABLE billing_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    payment_due_date TIMESTAMPTZ NOT NULL,
    amount_due NUMERIC(20, 2) NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'paid', 'overdue')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payments Table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    billing_schedule_id UUID NOT NULL REFERENCES billing_schedules(id) ON DELETE CASCADE,
    amount NUMERIC(20, 2) NOT NULL,
    payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('bank transfer', 'QRIS')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);