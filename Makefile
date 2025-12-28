
BINARY ?= cli
BIN_OUT ?= bin/$BINARY

protos:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	pb/tunnels/v1/tunnel_service.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	pb/rtunnel/v1/rtunnel_service.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	pb/routes/v1/route_service.proto

lint:
	@golangci-lint run ./...

test:
	@go test -v ./...

bin:
	rm ./bin/*
	GODEBUG=fips140=on CGO_ENABLED=0 go build -o ./bin/dino ./cmd/cli

server:
	docker build -t dino/server:latest --build-arg FIPS_ON=on -f docker/server.Dockerfile .

server/boring:
	docker build -t dino/server:latest --build-arg FIPS_ON=on -f docker/server.Dockerfile .

whoami:
	docker build -t dino/whoami:latest -f docker/whoami.Dockerfile services