CREATE TABLE roles (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE
);

CREATE TABLE members (
    eoa_address  TEXT PRIMARY KEY,
    role_id      INTEGER NOT NULL REFERENCES roles(id),
    created_at   INTEGER NOT NULL
);

CREATE TABLE user_attributes (
    eoa_address  TEXT PRIMARY KEY REFERENCES members(eoa_address) ON DELETE CASCADE,
    kyc          INTEGER NOT NULL DEFAULT 0 CHECK (kyc IN (0, 1)),
    blacklisted  INTEGER NOT NULL DEFAULT 0 CHECK (blacklisted IN (0, 1))
);

CREATE TABLE auth_tokens (
    token_hash   TEXT PRIMARY KEY,
    eoa_address  TEXT NOT NULL REFERENCES members(eoa_address) ON DELETE CASCADE,
    issued_at    INTEGER NOT NULL,
    expires_at   INTEGER NOT NULL
);

CREATE INDEX idx_auth_tokens_eoa ON auth_tokens(eoa_address);

CREATE TABLE challenges (
    eoa_address  TEXT PRIMARY KEY,
    nonce        TEXT NOT NULL,
    expires_at   INTEGER NOT NULL
);

CREATE TABLE lambdas (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT NOT NULL,
    role_id      INTEGER NOT NULL REFERENCES roles(id),
    description  TEXT,
    created_at   INTEGER NOT NULL,
    UNIQUE (name, role_id)
);

CREATE INDEX idx_lambdas_role ON lambdas(role_id);

CREATE TABLE lambda_rules (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    lambda_id      INTEGER NOT NULL REFERENCES lambdas(id) ON DELETE CASCADE,
    selector       TEXT NOT NULL,
    lhs_kind       TEXT NOT NULL CHECK (lhs_kind IN ('calldata', 'attribute')),
    lhs_offset     INTEGER,
    lhs_attribute  TEXT,
    condition      TEXT NOT NULL CHECK (condition IN ('eq', 'neq', 'gt', 'lt', 'gte', 'lte')),
    rhs_kind       TEXT NOT NULL CHECK (rhs_kind IN ('tx_origin', 'msg_sender', 'literal')),
    rhs_value      TEXT,
    CHECK (
        (lhs_kind = 'calldata' AND lhs_offset IS NOT NULL AND lhs_attribute IS NULL)
        OR (lhs_kind = 'attribute' AND lhs_attribute IS NOT NULL AND lhs_offset IS NULL)
    ),
    CHECK (
        (rhs_kind = 'literal' AND rhs_value IS NOT NULL)
        OR (rhs_kind IN ('tx_from', 'tx_to') AND rhs_value IS NULL)
    )
);

CREATE INDEX idx_lambda_rules_lambda ON lambda_rules(lambda_id);
CREATE INDEX idx_lambda_rules_selector ON lambda_rules(lambda_id, selector);

CREATE TABLE access_rules (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    contract_address  TEXT NOT NULL,
    function_selector TEXT NOT NULL,
    mode              TEXT NOT NULL CHECK (mode IN ('allow','deny')),
    UNIQUE (contract_address, function_selector)
);

CREATE INDEX idx_access_rules_contract ON access_rules(contract_address);

CREATE TABLE access_rule_entries (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    rule_id    INTEGER NOT NULL REFERENCES access_rules(id) ON DELETE CASCADE,
    role_id    INTEGER NOT NULL REFERENCES roles(id),
    lambda_id  INTEGER REFERENCES lambdas(id) ON DELETE RESTRICT,
    UNIQUE (rule_id, role_id)
);

CREATE INDEX idx_access_rule_entries_lambda ON access_rule_entries(lambda_id);
