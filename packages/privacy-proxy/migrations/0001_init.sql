-- Initial schema for privacy-proxy. See docs/system-design.md for rationale.

CREATE TABLE roles (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE
);
-- Role names are reconciled from src/roles.rs::ROLES at boot. No seed insert.

-- Every authenticated account: admin or user. The role-specific
-- attributes (if any) live in role-named attribute tables.
CREATE TABLE members (
    eoa_address  TEXT PRIMARY KEY,
    role_id      INTEGER NOT NULL REFERENCES roles(id),
    created_at   INTEGER NOT NULL
);

-- One row per `user`-role member. Admin role has no attribute table;
-- an admin's only attribute is their EOA, already in members.eoa_address.
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

CREATE TABLE access_rules (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    contract_address  TEXT NOT NULL,
    function_selector TEXT NOT NULL,
    mode              TEXT NOT NULL CHECK (mode IN ('allow','deny')),
    UNIQUE (contract_address, function_selector)
);

CREATE INDEX idx_access_rules_contract ON access_rules(contract_address);

CREATE TABLE access_rule_entries (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    rule_id      INTEGER NOT NULL REFERENCES access_rules(id) ON DELETE CASCADE,
    role_id      INTEGER NOT NULL REFERENCES roles(id),
    lambda_name  TEXT,
    UNIQUE (rule_id, role_id)
);
