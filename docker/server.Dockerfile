
ARG FIPS_ON=on

ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.25.5

FROM registry.opensuse.org/opensuse/bci/bci-base-fips:20260106.0-16.10 AS grpc_health_probe

RUN zypper --non-interactive install --no-recommends curl ca-certificates

ENV GRPC_HEALTH_PROBE_VERSION=v0.4.43
ENV TARGETOS=linux
ENV TARGETARCH=amd64

RUN curl -sfLo /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-${TARGETOS}-${TARGETARCH} && \
    chmod +x /bin/grpc_health_probe

FROM registry.opensuse.org/opensuse/bci/golang:1.25.5 AS builder

WORKDIR /usr/src/app

COPY go.sum go.mod ./
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/dino ./cmd/server

FROM registry.opensuse.org/opensuse/bci/bci-micro-fips:20260105.0-10.15

COPY --chown=65532:65532 --chmod=0755 --from=grpc_health_probe /bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --chown=65532:65532 --chmod=0755 --from=builder /bin/dino /usr/bin/dino

USER 65532

VOLUME [ "/qlog" ]

HEALTHCHECK --interval=30s --timeout=10s \
    CMD [ "/usr/bin/grpc_health_probe", "-addr", "127.0.0.1:50051" ]

EXPOSE 50051 4242

ENTRYPOINT [ "dino" ]
CMD [ ]
