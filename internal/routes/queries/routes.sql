
-- name: InsertRoute :one
-- InsertRoute insert new route record
INSERT INTO dino.routes ( 
    tunnel_name,
    hostname,
    destination_protocol,
    destination_ip,
    destination_port
) VALUES ( 
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: SelectRoute :one
SELECT
    *
FROM 
    dino.routes
WHERE
    id = $1;

-- name: SelectManyRoutes :many
SELECT
    id,
    hostname,
    is_active,
    created_at
FROM
    dino.routes
WHERE
    tunnel_name = $1
ORDER BY (id, created_at) ASC
LIMIT $2 OFFSET $3;

-- name: UpdateRoute :one
UPDATE dino.routes
SET
    hostname = $2,
    destination_ip = $3,
    destination_port = $4,
    destination_protocol = $5,
    is_active = $6
WHERE
    id = $1 
RETURNING *;

-- name: DeleteRoute :exec
DELETE FROM dino.routes WHERE id = $1;

-- name: SelectActiveRoute :one
-- SelectActiveRoute
SELECT
    t.id
FROM
    dino.routes as r
INNER JOIN 
    dino.tunnels as t
ON
    r.tunnel_name = t.identifier
WHERE
    hostname = $1 AND is_active = TRUE;

-- name: SelectRoutesMany :many
-- SelectRoutesMany
SELECT
    r.*,
    t.id
FROM 
    dino.routes as r
INNER JOIN
    dino.tunnels as t
ON
    r.tunnel_name = t.identitifer
WHERE
    r.tunnel_name = $1;
