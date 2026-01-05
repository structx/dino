
# Dino 🦕

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/structx/dino)


Routing & Tunnel Management

## 📖 Table of Contents

[🚀 Quick Start](#-quick-start)\
[🏗️ Architecture](docs/architecture.md)\
[⚙️ Config](docs/config.md)

## 🚀 Quick Start

```bash
# clone and enter repository
git clone https://github.com/structx/dino.git
cd dino

# create local volumes
mkdir -p .certs .local-volumes .local-volumes/pgdata

# set desired hostname for local tunnel
HOSTNAME="tunnel.dino.local"

# generate certs for local tunnel
openssl req -x509 -nodes -days 365 \
  -newkey ec:<(openssl ecparam -name prime256v1) \
  -keyout ./backend.key \
  -out ./backend.cert \
  -sha384 \
  -subj "/C=US/ST=State/L=City/O=DinoDev/CN=$HOSTNAME"

# update local hosts file (requires sudo)
echo "127.0.0.1 api.dino.local traefik.dino.local tunnel.dino.local whoami.dino.local" | sudo tee -a /etc/hosts

# dockerize local server and tunnel
docker build -t dino/server:latest \
  -f docker/server.Dockerfile .
docker build -t dino/tunnel:latest \
  -f docker/tunnel.Dockerfile .

# start tunnel infra snd servers
docker compose -f tunnel.compose.yaml up -d

# verify dino has started
docker logs dino
```
