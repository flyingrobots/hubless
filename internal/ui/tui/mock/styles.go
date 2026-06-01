package mock

import "github.com/charmbracelet/lipgloss"

// Styles mirrors the styling helpers used by the real TUI so mock rendering stays consistent and lintable.
type Styles struct {
	Statusline lipgloss.Style
	Footer     lipgloss.Style
}

// NewStyles constructs and returns a Styles value with the mock TUI's default styling.
//
// The returned Styles sets adaptive foreground colors for light and dark terminals.
func NewStyles() Styles {
	return Styles{
		Statusline: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#005F5A",
			Dark:  "#00B3A4",
		}).Bold(true),
		Footer: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#3B6F4D",
			Dark:  "#6B9F7F",
		}),
	}
}
