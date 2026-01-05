
ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.25.5

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

WORKDIR /usr/src/app

COPY go.sum go.mod ./
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/tunnel ./cmd/tunnel

CMD [ "tunnel" ]