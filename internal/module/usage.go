package module

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kubex-ecosystem/xtui/internal/module/info"
	"github.com/spf13/cobra"
)

var banners = []string{
	`
    ██╗  ██╗████████╗██╗   ██╗██╗
    ╚██╗██╔╝╚══██╔══╝██║   ██║██║
     ╚███╔╝    ██║   ██║   ██║██║
     ██╔██╗    ██║   ██║   ██║██║
    ██╔╝ ██╗   ██║   ╚██████╔╝██║
    ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝

 %sPowered by - Kubex Ecosystem %s%s
`,
	`
   __    __ ________          ______
  |  \  |  \        \        |      \
  | ▓▓  | ▓▓\▓▓▓▓▓▓▓▓__    __ \▓▓▓▓▓▓
   \▓▓\/  ▓▓  | ▓▓  |  \  |  \ | ▓▓
    >▓▓  ▓▓   | ▓▓  | ▓▓  | ▓▓ | ▓▓
   /  ▓▓▓▓\   | ▓▓  | ▓▓  | ▓▓ | ▓▓
  |  ▓▓ \▓▓\  | ▓▓  | ▓▓__/ ▓▓_| ▓▓_
  | ▓▓  | ▓▓  | ▓▓   \▓▓    ▓▓   ▓▓ \
   \▓▓   \▓▓   \▓▓    \▓▓▓▓▓▓ \▓▓▓▓▓▓
 %sPowered by - Kubex Ecosystem %s%s
`,
	`

  \  / |    |   |_ _|
    /  __|  |   |  |
    \  |    |   |  |
 _/\_\\__| \___/ ___|


 %sKubex %s%s
`,
}

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}

	manifest, err := info.GetManifest()
	if err != nil {
		description += ""
	} else {
		if manifest.GetDescription() != "" && description == "" {
			description = manifest.GetDescription()
		}
	}

	if hideBanner {
		return map[string]string{"banner": "", "description": description}
	}

	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = fmt.Sprintf(banners[bannerRandIndex], "\033[1;34m", manifest.GetVersion(), "\033[0m")

	return map[string]string{"banner": banner, "description": description}
}

// colorYellow, colorGreen, colorBlue, colorRed, and colorHelp are utility functions
// that return a string formatted with the specified color using the fatih/color package.
// These functions are used to colorize output in the CLI usage template.
// They are registered as template functions in the CLI usage template to allow
// coloring specific parts of the command usage output.
func colorYellow(s string) string {
	return color.New(color.FgYellow).SprintFunc()(s)
}

func colorGreen(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}

func colorBlue(s string) string {
	return color.New(color.FgBlue).SprintFunc()(s)
}

func colorRed(s string) string {
	return color.New(color.FgRed).SprintFunc()(s)
}

func colorHelp(s string) string {
	return color.New(color.FgCyan).SprintFunc()(s)
}

func hasServiceCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] == "true" {
			return true
		}
	}
	return false
}

func hasModuleCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] != "true" {
			return true
		}
	}
	return false
}

// SetUsageDefinition set the usage definition for the command.
func SetUsageDefinition(cmd *cobra.Command) {
	cobra.AddTemplateFunc("colorYellow", colorYellow)
	cobra.AddTemplateFunc("colorGreen", colorGreen)
	cobra.AddTemplateFunc("colorRed", colorRed)
	cobra.AddTemplateFunc("colorBlue", colorBlue)
	cobra.AddTemplateFunc("colorHelp", colorHelp)
	cobra.AddTemplateFunc("hasServiceCommands", hasServiceCommands)
	cobra.AddTemplateFunc("hasModuleCommands", hasModuleCommands)

	// Altera o template de uso do cobra
	cmd.SetUsageTemplate(cliUsageTemplate)
}

var cliUsageTemplate = `{{- if index .Annotations "banner" }}{{colorBlue (index .Annotations "banner")}}{{end}}{{- if (index .Annotations "description") }}
{{index .Annotations "description"}}
{{- end }}

{{colorYellow "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [args]{{end}}{{if gt (len .Aliases) 0}}

{{colorYellow "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{colorYellow "Example:"}}
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{colorYellow "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{colorGreen (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{colorYellow "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasAvailableInheritedFlags}}

{{colorYellow "Global Options:"}}
  {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasHelpSubCommands}}

{{colorYellow "Additional help topics:"}}
{{range .Commands}}{{if .IsHelpCommand}}
  {{colorGreen (rpad .CommandPath .CommandPathPadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasSubCommands}}

{{colorYellow (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}
`
