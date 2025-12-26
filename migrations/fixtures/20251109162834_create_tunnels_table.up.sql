
CREATE TABLE IF NOT EXISTS dino.tunnels (
    id UUID PRIMARY KEY DEFAULT extensions.uuid_generate_v4(),
    identifier VARCHAR(255) UNIQUE NOT NULL,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
