package mock

import "github.com/charmbracelet/lipgloss"

// Styles mirrors the styling helpers used by the real TUI so mock rendering stays consistent and lintable.
type Styles struct {
	Statusline lipgloss.Style
	Footer     lipgloss.Style
}

// NewStyles constructs and returns a Styles value with the mock TUI's default styling.
//
// The returned Styles sets:
// - Statusline: foreground color #00B3A4 and bold text.
// - Footer: foreground color #6B9F7F.
func NewStyles() Styles {
	return Styles{
		Statusline: lipgloss.NewStyle().Foreground(lipgloss.Color("#00B3A4")).Bold(true),
		Footer:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6B9F7F")),
	}
}
