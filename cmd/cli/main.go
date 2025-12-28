package main

import (
	"soft.structx.io/dino/cmd/cli/sub"

	// do not alter order
	_ "soft.structx.io/dino/cmd/cli/sub/route"
	_ "soft.structx.io/dino/cmd/cli/sub/tunnel"
	_ "soft.structx.io/dino/cmd/cli/sub/tunnel/credentials"
)

func main() {
	sub.Execute()
}
