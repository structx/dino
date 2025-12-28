
-- name: InsertTunnel :one
INSERT INTO dino.tunnels (
    identifier,
    token_hash
) VALUES (
    $1, $2
) RETURNING *;

-- name: SelectTunnel :one
SELECT
    *
FROM
    dino.tunnels
WHERE
    identifier = $1;

-- name: SelectTunnelToken :one
SELECT
    token_hash
FROM
    dino.tunnels
WHERE
    id = $1;

-- name: ListTunnels :many
SELECT
    id,
    identifier,
    created_at
FROM
    dino.tunnels
ORDER BY (id, created_at)
LIMIT $1 OFFSET $2;

-- name: UpdateTunnel :one
UPDATE dino.tunnels
SET
    identifier = $2
WHERE 
    identifier = $1
RETURNING *;

-- name: DeleteTunnel :execresult
DELETE FROM dino.tunnels WHERE identifier = $1;