package main

import (
	"context"
	"os"
	"strings"

	"github.com/kubex-ecosystem/logz"
	"github.com/kubex-ecosystem/xtui/cmd/cli"
	"github.com/kubex-ecosystem/xtui/internal/module"
	"github.com/kubex-ecosystem/xtui/internal/module/kbx"
	"github.com/spf13/cobra"
)

var (
	lgr           = logz.GetLoggerZ("XTUI") // Global logger
	mainCtx       = context.Background()
	xtuiContainer kbx.KbxModule
)

func wrap(ctx context.Context, mod *module.XTui) kbx.KbxModule {
	var hideBanner = os.Getenv("XTUI_HIDE_BANNER")
	if hideBanner == "" {
		hideBanner = "false"
	}

	if mod == nil {
		mod = module.NewXTui(ctx, kbx.InitArgs{
			HideBanner: strings.ToLower(hideBanner) == "true",
		})

		err := mod.Init(ctx, getSubCmdsList(mod))
		if err != nil {
			_ = logz.Errorf("Failed to initialize XTui module: %v", err)
		}
	}

	return mod
}

func UnWrap(mod kbx.KbxModule) *module.XTui {
	return mod.(*module.XTui)
}

func main() {
	// Set the global logger
	logz.SetGlobalLoggerZ(lgr)
	logz.SetGlobalLogger(lgr.Logger)
	logz.Logger.Logger = lgr.Logger.Logger

	xtuiContainer = wrap(mainCtx, nil)

	if err := xtuiContainer.Execute(); err != nil {
		os.Exit(1)
	}
}

func getSubCmdsList(x *module.XTui) []*cobra.Command {
	cmds := []*cobra.Command{}

	// Adiciona os comandos relacionados ao módulo
	cmds = append(cmds, cli.FormsRootCmd())
	cmds = append(cmds, cli.ViewsRootCmd())
	cmds = append(cmds, cli.AppsRootCmd())
	cmds = append(cmds, cli.PkgRootCmd())

	return cmds
}
