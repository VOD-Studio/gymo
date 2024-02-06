CREATE TABLE IF NOT EXISTS users
(
    id        serial PRIMARY KEY,
    username  varchar     NOT NULL,
    password  varchar     NOT NULL,
    create_at timestamptz NOT NULL DEFAULT NOW(),
    update_at timestamptz NOT NULL DEFAULT NOW()
)
