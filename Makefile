
# binary build command args
BINARY ?= dino
BIN_OUT ?= bin/$(BINARY)
ABS_BIN_PATH := $(shell pwd)/$(BIN_OUT)

CERT_OUT ?= .certs

# version control docker build args
GO_VERSION ?= "1.25.5"
ALPINE_VERSION ?= "3.23"

DINO_HOSTNAME ?= tunnel.dino.local

# fips at build time
# on 	: enable
# off	: disable
FIPS_MODE ?= "on"

.PHONY: protos lint test cli server tunnel certs all

.PHONY: clean
clean:
	@rm -rf .certs .local-volumes
	@mkdir -p .certs .local-volumes .local-volumes/pgdata

certs:
	@echo "Generating ECDSA P-256 certs for $(DINO_HOSTNAME)..." \
  mkdir -p $(CERT_OUT) \
	openssl req -x509 -nodes -days 365 \
		-newkey ec:<(openssl ecparam -name prime256v1) \
		-keyout $(CERT_OUT)/backend.key \
		-out $(CERT_OUT)/backend.cert \
		-sha384 \
		-subj "/C=US/ST=State/L=City/O=Dev/CN=$(DINO_HOSTNAME)" \
	chmod 0644 $(CERT_OUT)/backend.cert $(CERT_OUT)/backend.key

protos:
	@protoc --go_out=. --go_opt=paths=source_relative 			\
    --go-grpc_out=. --go-grpc_opt=paths=source_relative	\
    pb/tunnels/v1/tunnel_service.proto
	@protoc --go_out=. --go_opt=paths=source_relative 			\
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/rtunnel/v1/rtunnel_service.proto
	@protoc --go_out=. --go_opt=paths=source_relative 			\
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/routes/v1/route_service.proto

lint:
	@golangci-lint run ./...

test:
	@go test ./...

.PHONY: ci
ci: protos lint test 

alias: cli
	@echo "Configuring alias for $(BINARY)..."
	@ALIAS_CMD="alias dino='$(ABS_BIN_PATH)'"; \
	if [ -f "$$HOME/.zshrc" ]; then \
		CONF="$$HOME/.zshrc"; \
	elif [ -f "$$HOME/.bashrc" ]; then \
		CONF="$$HOME/.bashrc"; \
	else \
		echo "No config file found"; exit 1; \
	fi; \
	if grep -q "alias dino=" "$$CONF"; then \
		echo "Alias already exists in $$CONF. Skipping."; \
	else \
		echo "$$ALIAS_CMD" >> "$$CONF"; \
		echo "Added alias to $$CONF. Run 'source $$CONF' to activate."; \
	fi

unalias:
	@if [ -f "$$HOME/.zshrc" ]; then sed -i '/alias dino=/d' ~/.zshrc; fi
	@if [ -f "$$HOME/.bashrc" ]; then sed -i '/alias dino=/d' ~/.bashrc; fi
	@echo "Removed dino alias from shell config files."

cli:
	@CGO_ENABLED=0 go build -o $(BIN_OUT) ./cmd/cli	

server:
	docker build -t dino/server:latest 		 			    \
		--build-arg FIPS_ON=${FIPS_MODE} 	 			      \
		--build-arg GO_VERSION=${GO_VERSION}			    \
		--build-arg ALPINE_VERSION=${ALPINE_VERSION}	\
		-f docker/server.Dockerfile .

tunnel:
	docker build -t dino/tunnel:latest 		 			    \
		--build-arg FIPS_ON=${FIPS_MODE} 	 			      \
		--build-arg GO_VERSION=${GO_VERSION}			    \
		--build-arg ALPINE_VERSION=${ALPINE_VERSION}	\
		-f docker/tunnel.Dockerfile .

all: server tunnel