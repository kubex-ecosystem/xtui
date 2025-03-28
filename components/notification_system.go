package components

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type NotificationType string

const (
	Info    NotificationType = "info"
	Warning NotificationType = "warning"
	Error   NotificationType = "error"
)

type Notification struct {
	Message string
	Type    NotificationType
}

func ShowNotification(notification Notification) {
	var style lipgloss.Style

	switch notification.Type {
	case Info:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#75FBAB"))
	case Warning:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FDFF90"))
	case Error:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7698"))
	default:
		style = lipgloss.NewStyle()
	}

	fmt.Println(style.Render(notification.Message))
}

func DisplayInfoNotification(message string) {
	ShowNotification(Notification{Message: message, Type: Info})
}

func DisplayWarningNotification(message string) {
	ShowNotification(Notification{Message: message, Type: Warning})
}

func DisplayErrorNotification(message string) {
	ShowNotification(Notification{Message: message, Type: Error})
}
