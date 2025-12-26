package migrations

import (
	"embed"

	"go.uber.org/fx"
)

//go:embed fixtures/*.sql
var fs embed.FS

// Result
type Result struct {
	fx.Out

	EmbedFS embed.FS `group:"migrations"`
}

// Module
var Module = fx.Module("migrations", fx.Provide(invokeModule))

func invokeModule() Result {
	return Result{
		EmbedFS: fs,
	}
}
