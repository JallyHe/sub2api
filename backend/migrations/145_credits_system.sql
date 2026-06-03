-- 145: Credits system tables and user credit fields
-- Adds model_credit_rates, credit_ledger, and credit fields on users/subscription_plans.

-- Model credit consumption rates (admin-configurable)
CREATE TABLE model_credit_rates (
    id                            BIGSERIAL PRIMARY KEY,
    model_pattern                 VARCHAR(200) NOT NULL,
    credits_per_1k_tokens_input   BIGINT NOT NULL DEFAULT 1,
    credits_per_1k_tokens_output  BIGINT NOT NULL DEFAULT 3,
    priority                      INT NOT NULL DEFAULT 0,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_model_credit_rates_priority ON model_credit_rates (priority DESC);

-- Seed default rates
INSERT INTO model_credit_rates (model_pattern, credits_per_1k_tokens_input, credits_per_1k_tokens_output, priority)
VALUES
    ('claude-opus-4*',    15, 75, 100),
    ('claude-sonnet-4*',   3, 15,  90),
    ('claude-haiku-4*',    1,  5,  80),
    ('gpt-4o*',            5, 15,  70),
    ('gpt-4o-mini*',       1,  3,  60),
    ('gemini-2*-pro*',     5, 15,  50),
    ('gemini-2*-flash*',   1,  3,  40),
    ('*',                  2,  8,   0);

-- Append-only credit transaction ledger
CREATE TABLE credit_ledger (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    delta         BIGINT NOT NULL,
    reason        VARCHAR(50) NOT NULL,
    ref_id        VARCHAR(100),
    balance_after BIGINT NOT NULL,
    model         VARCHAR(100),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credit_ledger_user_created ON credit_ledger (user_id, created_at DESC);
CREATE INDEX idx_credit_ledger_user_id ON credit_ledger (user_id);

-- Credit fields on users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS credit_balance    BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS credit_expires_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS credit_plan_id    BIGINT;

-- Credits field on subscription_plans
ALTER TABLE subscription_plans
    ADD COLUMN IF NOT EXISTS credits BIGINT NOT NULL DEFAULT 0;
