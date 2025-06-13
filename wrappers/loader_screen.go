package wrappers

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rafa-mori/xtui/types"
)

var (
	loaderSpinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
	loaderHelpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	loaderDotStyle      = loaderHelpStyle.UnsetMargins()
	loaderDurationStyle = loaderDotStyle
	loaderAppStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	loaderErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	loaderTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	loaderMessageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))

	loaderSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	loaderWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Italic(true)
	loaderInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Underline(true)
)

type LoaderMsg struct {
	Message string
}

type LoaderCloseMsg struct{}

type loaderModel struct {
	spinner  spinner.Model
	messages []types.LoaderMessage
	quitting bool
	err      error
}

func newLoaderModel() loaderModel {
	s := spinner.New()
	s.Style = loaderSpinnerStyle
	return loaderModel{
		spinner:  s,
		messages: []types.LoaderMessage{},
	}
}

func (m loaderModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m loaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case LoaderMsg:
		m.messages = append(m.messages, types.LoaderMessage{Message: msg.Message})
		return m, nil
	case LoaderCloseMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m loaderModel) View() string {
	var s string

	if m.quitting {
		s += loaderTitleStyle.Render("Done!")
	} else {
		s += loaderTitleStyle.Render(m.spinner.View() + " Working...")
	}

	s += "\n\n"

	for _, msg := range m.messages {
		if strings.Contains(msg.Message, "Error: ") {
			s += loaderErrorStyle.Render(msg.Message) + "\n"
		} else if strings.Contains(msg.Message, "Success: ") {
			s += loaderSuccessStyle.Render(msg.Message) + "\n"
		} else if strings.Contains(msg.Message, "Warning: ") {
			s += loaderWarningStyle.Render(msg.Message) + "\n"
		} else {
			s += loaderMessageStyle.Render(msg.Message) + "\n"
		}
	}

	if !m.quitting {
		s += loaderHelpStyle.Render("Press any key to quit")
	}

	if m.quitting {
		s += "\n"
	}

	return loaderAppStyle.Render(s)
}

func StartLoader(messages chan tea.Msg) error {
	p := tea.NewProgram(newLoaderModel())

	go func() {
		for msg := range messages {
			p.Send(msg)
		}
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func NavigateAndExecuteLoader(messages chan tea.Msg) error {
	p := tea.NewProgram(newLoaderModel())

	go func() {
		for msg := range messages {
			p.Send(msg)
		}
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
