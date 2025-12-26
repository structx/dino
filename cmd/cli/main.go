package main

import (
	"github.com/structx/dino/cmd/cli/sub"

	// do not alter order
	_ "github.com/structx/dino/cmd/cli/sub/route"
	_ "github.com/structx/dino/cmd/cli/sub/tunnel"
	_ "github.com/structx/dino/cmd/cli/sub/tunnel/credentials"
)

func main() {
	sub.Execute()
}
