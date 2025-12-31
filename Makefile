
# binary build command args
BINARY ?= cli
BIN_OUT ?= bin/$BINARY

# version control docker build args
GO_VERSION ?= "1.25.5"
ALPINE_VERSION ?= "3.23"

# fips at build time
# on 	: enable
# off	: disable
FIPS_MODE ?= "on"

protos:
	protoc --go_out=. --go_opt=paths=source_relative 		\
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative	\
    	pb/tunnels/v1/tunnel_service.proto
	protoc --go_out=. --go_opt=paths=source_relative 		\
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	pb/rtunnel/v1/rtunnel_service.proto
	protoc --go_out=. --go_opt=paths=source_relative 		\
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	pb/routes/v1/route_service.proto

lint:
	@golangci-lint run ./...

test:
	@go test -v ./...

cli:
	@GODEBUG=fips140=${FIPS_MODE} CGO_ENABLED=0 go build -o ./bin/dino ./cmd/cli	

server:
	docker build -t dino/server:latest 		 			\
		--build-arg FIPS_ON=${FIPS_MODE} 	 			\
		--build-arg GO_VERSION=${GO_VERSION}			\
		--build-arg ALPINE_VERSION=${ALPINE_VERSION}	\
		-f docker/server.Dockerfile .
