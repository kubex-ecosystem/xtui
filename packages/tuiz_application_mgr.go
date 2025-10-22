package packages

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gl "github.com/kubex-ecosystem/logz/logger"
)

type KbxDepsModel struct {
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
	kbxDepsCurrentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	kbxDepsDoneStyle           = lipgloss.NewStyle()
	kbxDepsCheckMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func KbxDepsNewModel(apps []string, path string, yes bool, quiet bool) KbxDepsModel {
	p := progress.New(
		progress.WithScaledGradient("#F72100", "#F66600"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	return KbxDepsModel{
		apps:     apps,
		path:     path,
		yes:      yes,
		quiet:    quiet,
		spinner:  s,
		progress: p,
	}
}

func (m *KbxDepsModel) Init() tea.Cmd {
	return tea.Batch(kbxDepsDownloadAndInstall(m.apps[m.index], m.path, m.yes, m.quiet))
}

func (m *KbxDepsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case kbxDepsInstalledPkgMsg:
		if m.index >= len(m.apps)-1 {
			m.done = true
			return m, tea.Quit
		}

		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index+1) / float64(len(m.apps)))

		return m, tea.Batch(
			progressCmd,
			kbxDepsDownloadAndInstall(m.apps[m.index], m.path, m.yes, m.quiet),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		kbxDepsNewModel, cmd := m.progress.Update(msg)
		if kbxDepsNewModel, ok := kbxDepsNewModel.(progress.Model); ok {
			m.progress = kbxDepsNewModel
		}
		return m, cmd
	}
	return m, nil
}

func (m *KbxDepsModel) View() string {
	n := len(m.apps)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	depCount := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := kbxDepsMax(0, m.width-lipgloss.Width(spin+prog+depCount))

	infoMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Installation progress: ")
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("[pkgz] " + infoMsg + prog + depCount)

	cellsRemaining := kbxDepsMax(0, m.width-lipgloss.Width(info+spin))
	gap := strings.Repeat(" ", cellsRemaining)

	if m.done {
		depsQty := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(fmt.Sprintf("%d", n))
		return info + gap + "\n" + m.renderInstalledApps() + kbxDepsDoneStyle.Render(fmt.Sprintf("[pkgz] Done! Installed %s applications.\n", depsQty))
	}

	return info + gap + "\n" + m.renderInstalledApps()
}

func (m *KbxDepsModel) renderInstalledApps() string {
	var installedApps strings.Builder
	for i := 0; i < m.index+1; i++ {
		dep := m.apps[i]
		depRendered := kbxDepsCurrentPkgNameStyle.Render(dep)
		// Insere um ícone que simbolize a instalação do pacote
		infoBadge := lipgloss.NewStyle().Inline(true).Render("[XTui] 📦 -> ")
		infoMsg := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("6")).Render("Installed application: ")
		installedApps.WriteString(fmt.Sprintf("%s%s%s %s\n", infoBadge, infoMsg, kbxDepsCheckMark, depRendered))
	}
	return installedApps.String()
}

type kbxDepsInstalledPkgMsg string

func kbxDepsDownloadAndInstall(dep string, path string, yes bool, quiet bool) tea.Cmd {
	return func() tea.Msg {
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
		cmd.Stdout = Writer("XTui")
		cmd.Stdin = stdin
		if err := cmd.Run(); err != nil {
			log.Println("error installing application: "+dep+" "+err.Error(), "XTui")
			return tea.Quit()
		}
		return kbxDepsInstalledPkgMsg(dep)
	}
}

func kbxDepsMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// InstallDepsWithUI installs dependencies in a terminal UI with a progress bar
func InstallDepsWithUI(args ...string) error {
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
	model := KbxDepsNewModel(apps, path, yes, quiet)
	p := tea.NewProgram(&model)
	_, err := p.Run()
	defer p.Quit()
	if err != nil {
		gl.Log("error", "error running dependencies installation: "+err.Error())
		return nil
	}
	return nil
}

func Writer(module string) io.Writer {
	currentPid := os.Getpid()
	logFileName := strings.Join([]string{module, "logz_", strconv.Itoa(currentPid), ".log"}, "")
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil || cacheDir == "" {
		cacheDir = os.TempDir()
	}
	logFilePath := filepath.Join(cacheDir, logFileName)
	if logFileStatErr := os.Remove(logFilePath); logFileStatErr == nil {
		cmdRm := fmt.Sprintf("rm -f %s", logFilePath)
		if _, cmdRmErr := exec.Command("bash", "-c", cmdRm).Output(); cmdRmErr != nil {
			fmt.Println(cmdRmErr)
			return os.Stdout
		}
	}
	if logFilePathErr := os.MkdirAll(filepath.Dir(logFilePath), 0777); logFilePathErr != nil {
		fmt.Println(logFilePathErr)
		return os.Stdout
	}
	logFile, logFileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if logFileErr != nil {
		fmt.Println(logFileErr)
		return os.Stdout
	}
	return logFile
}
