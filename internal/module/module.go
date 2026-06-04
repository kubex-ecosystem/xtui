// Package module provides internal types and functions for the Grompt application.
package module

import (
	"context"
	"os"
	"strings"

	flow "github.com/kubex-ecosystem/ethyr/tools/flow"
	control "github.com/kubex-ecosystem/ethyr/tools/flow/control"
	fsm "github.com/kubex-ecosystem/ethyr/tools/flow/fsm"
	gl "github.com/kubex-ecosystem/logz"

	kbxmod "github.com/kubex-ecosystem/xtui/internal/module/kbx"
	"github.com/kubex-ecosystem/xtui/internal/module/version"

	"github.com/spf13/cobra"
)

// XTui representa a estrutura do módulo ui.
type XTui struct {
	*flow.Flow

	engine  *flow.Engine[fsm.State]
	pass    flow.Pass[fsm.State]
	control *control.ManagerControl[control.StepFlag]

	Command *cobra.Command

	HideBanner bool
	Debug      bool
}

func NewXTui(ctx context.Context, cfg kbxmod.InitArgs) *XTui {
	x := &XTui{
		Flow:       flow.NewFlow(),
		HideBanner: cfg.HideBanner,
	}

	x.Command = &cobra.Command{
		Use:         x.Module(),
		Short:       x.ShortDescription(),
		Aliases:     []string{x.Alias()},
		Example:     x.concatenateExamples(),
		Annotations: GetDescriptions([]string{x.ShortDescription(), x.LongDescription()}, x.HideBanner),
	}

	SetUsageDefinition(x.Command)

	x.Command.PersistentFlags().BoolVarP(&x.HideBanner, "hide-banner", "H", false, "Hide/disable the ASCII banner.")
	x.Command.PersistentFlags().BoolVarP(&x.Debug, "debug", "D", false, "Enable debug mode.")

	return x
}

// Init initializes the module
func (m *XTui) Init(ctx context.Context, subCmds []*cobra.Command) error {
	gl.Debugf("Creating command for XTuI with flags: %v", os.Args)

	SetUsageDefinition(m.Command)

	for _, cmd := range subCmds {
		if cmd != nil {
			SetUsageDefinition(cmd)
			for _, subCmd := range cmd.Commands() {
				if subCmd != nil {
					SetUsageDefinition(subCmd)
					if len(subCmd.Annotations) == 0 {
						subCmd.Annotations = GetDescriptions([]string{subCmd.Short, subCmd.Long}, m.HideBanner)
						if !strings.Contains(strings.Join(os.Args, " "), subCmd.Use) {
							if subCmd.Short == "" {
								subCmd.Short = subCmd.Annotations["description"]
							}
						}
					}
				}
			}
			m.Command.AddCommand(cmd)
		}
	}

	// Set initial state
	if m.FSM == nil {
		m.FSM = fsm.NewFSM(kbxmod.XTUIIdle, []fsm.Transition{
			{From: kbxmod.XTUIIdle, Event: kbxmod.EvBoot, To: kbxmod.XTUIBootstrapping},
			{From: kbxmod.XTUIBootstrapping, Event: kbxmod.EvConfigure, To: kbxmod.XTUIConfiguring},
			{From: kbxmod.XTUIConfiguring, Event: kbxmod.EvRead, To: kbxmod.XTUIReading},
			{From: kbxmod.XTUIReading, Event: kbxmod.EvProcess, To: kbxmod.XTUIProcessing},
			{From: kbxmod.XTUIProcessing, Event: kbxmod.EvWait, To: kbxmod.XTUIWaiting},
			{From: kbxmod.XTUIWaiting, Event: kbxmod.EvStop, To: kbxmod.XTUIStopping},
		})
	}

	// Set initial state flag
	if !m.FSM.Can(kbxmod.EvBoot) {
		return gl.Errorf("failed to trigger event %d", kbxmod.EvBoot)
	}

	// Transição para estado inicial
	gl.Debugf("Triggering event %d", kbxmod.EvBoot)
	if !m.FSM.Trigger(kbxmod.EvBoot) {
		return gl.Errorf("failed to trigger event %d", kbxmod.EvBoot)
	}

	gl.Debugf("State machine: %d", m.FSM.Current())
	return nil
}

// Alias retorna o alias do módulo ui.
func (m *XTui) Alias() string {
	return ""
}

// ShortDescription retorna uma descrição curta do módulo ui.
func (m *XTui) ShortDescription() string {
	return "Terminal UI"
}

// LongDescription retorna uma descrição longa do módulo ui.
func (m *XTui) LongDescription() string {
	return "Terminal XTUI module. It allows you to interact with the terminal using a graphical interface."
}

// Usage retorna a forma de uso do módulo ui.
func (m *XTui) Usage() string {
	return "xui [command] [args]"
}

// Examples retorna exemplos de uso do módulo ui.
func (m *XTui) Examples() []string {
	return []string{"xtui [command] [args]", "xtui logz -o 'file.log'", "xtui deps -o 'install'", "xtui tcp-status '127.0.0.1:8080'"}
}

// Active verifica se o módulo ui está ativo.
func (m *XTui) Active() bool {
	return true
}

// Module retorna o nome do módulo ui.
func (m *XTui) Module() string {
	return "xtui"
}

// Execute executa o comando especificado para o módulo ui.
func (m *XTui) Execute() error {
	if m.Command == nil {
		return nil
	}

	return m.Command.Execute()
}

// concatenateExamples concatena os exemplos de uso do módulo.
func (m *XTui) concatenateExamples() string {
	examples := ""
	for _, example := range m.Examples() {
		examples += string(example) + "\n  "
	}
	return examples
}

// Command retorna o comando cobra para o módulo.
func (m *XTui) loadCommands(cmd *cobra.Command, cms []*cobra.Command) {
	gl.GetLogger("XTuI")

	cmd.Annotations = GetDescriptions([]string{cmd.Short, cmd.Long}, m.HideBanner)

	// Set usage definitions for the command and its subcommands
	SetUsageDefinition(cmd)
	for _, subCmd := range cms {
		SetUsageDefinition(subCmd)
		if !strings.Contains(strings.Join(os.Args, " "), subCmd.Use) {
			if subCmd.Short == "" {
				subCmd.Short = subCmd.Annotations["description"]
			}
		}
	}
}

func (m *XTui) Version() kbxmod.VersionService {
	return version.NewVersionService()
}
