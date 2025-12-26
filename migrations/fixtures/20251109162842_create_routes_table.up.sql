
CREATE TABLE IF NOT EXISTS dino.routes (
    id UUID PRIMARY KEY DEFAULT extensions.uuid_generate_v4(),
    tunnel_name VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) UNIQUE NOT NULL,
    destination_protocol VARCHAR(25) NOT NULL,
    destination_ip VARCHAR(25) NOT NULL,
    destination_port INTEGER NOT NULL,
    -- auth_token_hash
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (tunnel_name) REFERENCES dino.tunnels (identifier)
);
