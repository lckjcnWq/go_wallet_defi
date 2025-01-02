package main

import (
	_ "go-wallet-defi/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"go-wallet-defi/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
