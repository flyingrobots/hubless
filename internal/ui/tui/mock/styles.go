package mock

import "github.com/charmbracelet/lipgloss"

// Styles groups lipgloss styles used across the mock TUI surfaces.
type Styles struct {
	Statusline     lipgloss.Style
	Footer         lipgloss.Style
	Body           lipgloss.Style
	Section        lipgloss.Style
	SectionSel     lipgloss.Style
	Hint           lipgloss.Style
	Key            lipgloss.Style
	Title          lipgloss.Style
	Subtitle       lipgloss.Style
	Border         lipgloss.Style
	Accent         lipgloss.Style
	StatusOpen     lipgloss.Style
	StatusProgress lipgloss.Style
	StatusClosed   lipgloss.Style
	PriorityHigh   lipgloss.Style
	PriorityMedium lipgloss.Style
	PriorityLow    lipgloss.Style
	Overlay        lipgloss.Style
}

func newStyles() Styles {
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	return Styles{
		Statusline:     lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Background(lipgloss.Color("236")).Padding(0, 1).Bold(true),
		Footer:         lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Background(lipgloss.Color("236")).Padding(0, 1),
		Body:           lipgloss.NewStyle().Padding(0, 1),
		Section:        lipgloss.NewStyle().Padding(0, 1).MarginBottom(0),
		SectionSel:     lipgloss.NewStyle().Padding(0, 1).MarginBottom(0).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("60")).Bold(true),
		Hint:           lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Key:            lipgloss.NewStyle().Foreground(lipgloss.Color("115")).Bold(true),
		Title:          lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		Subtitle:       lipgloss.NewStyle().Foreground(lipgloss.Color("110")).Bold(true),
		Border:         border,
		Accent:         lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true),
		StatusOpen:     lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true),
		StatusProgress: lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
		StatusClosed:   lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true),
		PriorityHigh:   lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
		PriorityMedium: lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
		PriorityLow:    lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Overlay:        lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Foreground(lipgloss.Color("250")).Background(lipgloss.Color("57")),
	}
}
