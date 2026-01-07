
ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.25.5

FROM registry.opensuse.org/opensuse/bci/golang:1.25.5 AS builder

WORKDIR /usr/src/app

COPY go.sum go.mod ./
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/tunnel ./cmd/tunnel

FROM registry.opensuse.org/opensuse/bci/bci-micro-fips:20260105.0-10.15

COPY --chown=65532:65532 --chmod=0755 --from=builder /bin/dino /usr/bin/dino

USER 65532

ENTRYPOINT [ "tunnel" ]
CMD [ ]
