CREATE TABLE admin_users
(
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,

    display_name  TEXT,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,

    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admin_sessions
(
    id            BIGSERIAL PRIMARY KEY,
    admin_user_id BIGINT      NOT NULL REFERENCES admin_users (id) ON DELETE CASCADE,

    token_hash    TEXT        NOT NULL UNIQUE,

    expires_at    TIMESTAMPTZ NOT NULL,
    revoked_at    TIMESTAMPTZ,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at  TIMESTAMPTZ
);

CREATE INDEX idx_admin_sessions_user_id
    ON admin_sessions (admin_user_id);

CREATE INDEX idx_admin_sessions_active
    ON admin_sessions (admin_user_id, expires_at)
    WHERE revoked_at IS NULL;
