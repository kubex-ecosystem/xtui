package wrappers

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gl "github.com/kubex-ecosystem/logz/logger"
)

type AppDepsModel struct {
	apps     []string
	path     string
	yes      bool
	quiet    bool
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

var (
	appDepsCurrentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	appDepsDoneStyle           = lipgloss.NewStyle()
	appDepsCheckMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func NewAppDepsModel(apps []string, path string, yes bool, quiet bool) AppDepsModel {
	p := progress.New(
		progress.WithScaledGradient("#F72100", "#F66600"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	return AppDepsModel{
		apps:     apps,
		path:     path,
		yes:      yes,
		quiet:    quiet,
		spinner:  s,
		progress: p,
	}
}

func (m *AppDepsModel) Init() tea.Cmd {
	return tea.Batch(downloadAndInstall(m.apps[m.index], m.path, m.yes, m.quiet))
}

func (m *AppDepsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
		if m.index >= len(m.apps)-1 {
			m.done = true
			return m, tea.Quit
		}

		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index+1) / float64(len(m.apps)))

		return m, tea.Batch(
			progressCmd,
			downloadAndInstall(m.apps[m.index], m.path, m.yes, m.quiet),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m *AppDepsModel) View() string {
	n := len(m.apps)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	depCount := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := appDepsMax(0, m.width-lipgloss.Width(spin+prog+depCount))

	infoMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Installation progress: ")
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("[pkgz] " + infoMsg + prog + depCount)

	cellsRemaining := appDepsMax(0, m.width-lipgloss.Width(info+spin))
	gap := strings.Repeat(" ", cellsRemaining)

	if m.done {
		depsQty := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(fmt.Sprintf("%d", n))
		return info + gap + "\n" + m.renderInstalledApps() + appDepsDoneStyle.Render(fmt.Sprintf("[pkgz] Done! Installed %s applications.\n", depsQty))
	}

	return info + gap + "\n" + m.renderInstalledApps()
}

func (m *AppDepsModel) renderInstalledApps() string {
	var installedApps strings.Builder
	for i := 0; i < m.index+1; i++ {
		dep := m.apps[i]
		depRendered := appDepsCurrentPkgNameStyle.Render(dep)
		infoBadge := lipgloss.NewStyle().Inline(true).Render("[pkgz] ")
		infoMsg := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("6")).Render("Installed application: ")
		installedApps.WriteString(fmt.Sprintf("%s%s%s %s\n", infoBadge, infoMsg, appDepsCheckMark, depRendered))
	}
	return installedApps.String()
}

type installedPkgMsg string

func downloadAndInstall(dep string, path string, yes bool, quiet bool) tea.Cmd {
	return func() tea.Msg {
		gl.Log("info", "Installing application: ")
		pkg := dep
		if strings.Contains(dep, "/") {
			pkg = strings.Split(dep, "/")[1]
		}
		pathFlag := ""
		if path != "" {
			pathFlag = "-p " + path
		}
		yesFlag := ""
		if yes {
			yesFlag = "-y"
		}
		quietFlag := ""
		if quiet {
			quietFlag = "-qq"
		}
		cmd := exec.Command("sudo", "apt-get", "install", pkg, pathFlag, yesFlag, quietFlag)
		stdin := strings.NewReader("s\n")
		cmd.Stdout = log.Writer()
		cmd.Stdin = stdin
		if err := cmd.Run(); err != nil {
			gl.Log("error", "error installing application: "+err.Error())
			return tea.Quit()
		}
		gl.Log("info", "Application installed.")
		return installedPkgMsg(dep)
	}
}

func appDepsMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func InstallDependenciesWithUI(args ...string) error {
	if len(args) < 4 {
		gl.Log("error", "missing arguments")
		return nil
	}
	apps := strings.Split(args[0], " ")
	if len(apps) == 0 {
		gl.Log("error", "no applications requested")
		return nil
	}
	path := args[1]
	yes := args[2] == "true"
	quiet := args[3] == "true"

	availableProperties := getAvailableProperties()
	if len(availableProperties) > 0 {
		adaptedArgs := adaptArgsToProperties(args, availableProperties)
		return InstallDependenciesWithUI(adaptedArgs...)
	}

	model := NewAppDepsModel(apps, path, yes, quiet)
	p := tea.NewProgram(&model)
	_, err := p.Run()
	defer p.Quit()
	if err != nil {
		gl.Log("error", "error running dependencies installation: "+err.Error())
		return nil
	}
	return nil
}

func getAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

func adaptArgsToProperties(args []string, properties map[string]string) []string {
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}

func NavigateAndExecuteApplication(apps []string, path string, yes bool, quiet bool) error {
	model := NewAppDepsModel(apps, path, yes, quiet)
	p := tea.NewProgram(&model)
	_, err := p.Run()
	defer p.Quit()
	if err != nil {
		gl.Log("error", "error running application navigation and execution: "+err.Error())
		return nil
	}
	return nil
}
