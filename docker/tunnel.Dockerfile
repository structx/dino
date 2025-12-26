
FROM golang:1.25.3-alpine3.22

WORKDIR /usr/src/app

COPY go.sum go.mod ./
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/tunnel ./cmd/tunnel

CMD [ "tunnel" ]