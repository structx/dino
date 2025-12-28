package client

import "context"

type contextKey string

var clientKey contextKey = "dino_client"

// WithContext
func WithContext(ctx context.Context, client Client) context.Context {
	if client == nil {
		panic("client is nil")
	}
	return context.WithValue(ctx, clientKey, client)
}

// FromContext
func FromContext(ctx context.Context) Client {
	client, ok := ctx.Value(clientKey).(Client)
	if !ok {
		panic("invalid context value")
	}
	return client
}
