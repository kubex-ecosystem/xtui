package wrappers

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gl "github.com/kubex-ecosystem/logz"
)

type KbxProgressBarModel struct {
	ctx      context.Context
	msg      ProgressBarMsg
	err      error
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

type ProgressBar interface {
	tea.Model
}

type ProgressBarTickMsg struct{}
type ProgressBarMsg struct {
	current int
	total   int
	file    string
	title   string
	percent float32
}

var (
	kbxProgressBarCurrentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	kbxProgressBarDoneStyle           = lipgloss.NewStyle()
	kbxProgressBarCheckMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func NewKbXProgressBarModel(ctx context.Context, file string, current int, total int, title string) ProgressBar {
	p := progress.New(
		progress.WithScaledGradient("#F72100", "#F66600"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	return &KbxProgressBarModel{
		ctx:      ctx,
		msg:      ProgressBarMsg{current: current, file: file, total: total, title: title},
		spinner:  s,
		progress: p,
		done:     false,
		err:      nil,
	}
}

func (m *KbxProgressBarModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return ProgressBarTickMsg{}
		}),
	)
}

func (m *KbxProgressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			gl.Debug("Progress bar key pressed")
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		var cmd tea.Cmd
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case ProgressBarTickMsg:
		cmd := m.calculatePercent()

		var tickCmd tea.Cmd
		if !m.done {
			tickCmd = tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ProgressBarTickMsg{}
			})
		}

		return m, tea.Batch(cmd, tickCmd)
	case ProgressBarMsg:
		m.msg.current = msg.current
		m.msg.total = msg.total
		m.msg.percent = msg.percent
		m.msg.file = msg.file
		m.msg.title = msg.title
		return m, m.calculatePercent()
	}
	return m, nil
}

func (m *KbxProgressBarModel) View() string {
	if m.err != nil {
		return gl.Sprintf("Error: %s\n", m.err.Error())
	}

	pStr := gl.Sprintf(" %d/%d (%.1f%%)", m.msg.current, m.msg.total, m.msg.percent*100)
	spin := m.spinner.View() + " "
	prog := m.progress.View()

	infoMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Progress: ")
	info := "[XTui] " + infoMsg + spin + prog + pStr

	if m.width > 0 {
		info = lipgloss.NewStyle().MaxWidth(m.width).Render(info)
	}

	if m.done {
		return info + "\n" + m.renderProgressBar() + kbxProgressBarDoneStyle.Render("[XTui] Done! \n")
	}

	return info + "\n" + m.renderProgressBar()
}

func (m *KbxProgressBarModel) renderProgressBar() string {
	fileRendered := kbxProgressBarCurrentPkgNameStyle.Render(m.msg.title)
	infoBadge := lipgloss.NewStyle().Inline(true).Render("[XTui] 📦 -> ")
	infoMsg := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("6")).Render("Status: ")
	return gl.Sprintf("%s%s%s %s\n", infoBadge, infoMsg, kbxProgressBarCheckMark, fileRendered)
}

func (m *KbxProgressBarModel) calculatePercent() tea.Cmd {
	r, err := os.ReadFile(m.msg.file)
	if err != nil || len(r) == 0 {
		return nil
	}

	content := strings.TrimSpace(string(r))

	switch content {
	case "done":
		m.msg.percent = 1.0
		m.msg.current = m.msg.total
		m.done = true
		return tea.Quit
	case "error", "exit":
		m.err = gl.Errorf("process terminated with: %s", content)
		m.msg.percent = 1.0
		m.done = true
		return tea.Quit
	}

	val, err := strconv.Atoi(content)
	if err != nil {
		return nil
	}

	if val != m.msg.current {
		m.msg.current = val

		percent := float64(val) / float64(m.msg.total)
		if percent < 0 {
			percent = 0
		}
		if percent > 1.0 {
			percent = 1.0
		}
		m.msg.percent = float32(percent)

		progressCmd := m.progress.SetPercent(percent)

		if val >= m.msg.total {
			m.done = true
			return tea.Batch(progressCmd, tea.Quit)
		}

		return progressCmd
	}

	return nil
}

func KbxProgressBarRun(file string, current int, total int, percent float32, title string) tea.Cmd {
	return func() tea.Msg {
		return ProgressBarMsg{
			current: current,
			file:    file,
			total:   total,
			title:   title,
			percent: percent,
		}
	}
}

func kbxProgressBarMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ProgressBarWithUI installs dependencies in a terminal UI with a progress bar
func ProgressBarWithUI(ctx context.Context, file string, current int, total int, percent float32, title string) error {
	gl.Debug(gl.Sprintf("Progress bar: %s/%d/%s", file, total, title))
	if file == "" {
		return gl.Error("progress bar file not specified")
	}

	if total <= 0 || total > 100 {
		total = 100
	}

	model := NewKbXProgressBarModel(ctx, file, current, total, title)

	p := tea.NewProgram(
		model.(*KbxProgressBarModel),
		tea.WithInput(os.Stdin),
		tea.WithContext(ctx),
	)
	defer func() {
		gl.Debug(gl.Sprintf("Progress bar cleanup"))
		p.Quit()
	}()

	m, err := p.Run()
	if err != nil {
		return gl.Errorf("Failed to run progress bar: %s", err.Error())
	}
	if pb, ok := m.(*KbxProgressBarModel); ok {
		if pb.err != nil {
			return pb.err
		}
	}

	return nil
}
