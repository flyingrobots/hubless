package mock

import "github.com/charmbracelet/lipgloss"

// Styles mirrors the styling helpers used by the real TUI so mock rendering stays consistent and lintable.
type Styles struct {
	Statusline lipgloss.Style
	Footer     lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		Statusline: lipgloss.NewStyle().Foreground(lipgloss.Color("#00B3A4")).Bold(true),
		Footer:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6B9F7F")),
	}
}
