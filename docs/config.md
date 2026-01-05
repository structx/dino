# ⚙️ Config

servers and tunnels can be configured using environment variables. this page will serve a dictionary for all variables and their purpose. 

## Server

`LOG_LEVEL`         `DEBUG`             global log level

`SERVER_HOST`       `127.0.0.1`         server api bind host\
`SERVER_PORT`       `50051`             server api bind port

`DB_USERNAME`       `dino`              database user\
`DB_PASSWORD`       `dino`              database user password\
`DB_HOST`           `postgres`          database host\
`DB_PORT`           `5432`              database port\
`DB_NAME`           `dino`              database name\
`DB_EXRA_PARAMS`    `sslmode=disable`   pgxpool dial string params

`AUTH_JWT_DURATION` `-1`                jwt duration\
`AUTH_JWT_ISSUER`   `dino.local`        jwt issuer\
`AUTH_JWT_AUD`      `dino`              jwt audience (`,` split list)

## Tunnel

`TUNNEL_ID`                                     tunnel id\
`TUNNEL_TOKEN`                                  tunnel token\
`TUNNEL_ENDPOINT`   `tunnel.dino.local:4222`    tunnel endpoint

## Proxy

`http` proxy configuration

`PROXY_HOST`                `127.0.0.1`     http server host\
`PROXY_PORT`                `8080`          http server port\
`PROXY_TIMEOUT`             `15`            max duration to read entire request (seconds)\
`PROXY_READ_HEADER_TIMEOUT` `15`            duration to read request headers (seconds)\
`PROXY_WRITE_TIMEOUT`       `15`            max duration to write response\
`PROXY_IDLE_TIMEOUT`        `30`            duration to wait for next request when keepalive is enabled