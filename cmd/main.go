package main

import (
	"os"

	"github.com/kubex-ecosystem/xtui/internal/module"
)

func main() {
	if err := module.RegX().Execute(); err != nil {
		os.Exit(1)
	}
}
