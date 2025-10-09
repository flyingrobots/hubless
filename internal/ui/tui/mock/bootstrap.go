package mock

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	mockdata "github.com/flyingrobots/hubless/internal/mock"
)

// NewProgram constructs a Bubbletea program for the mocked Hubless TUI.
func NewProgram() *tea.Program {
	now := time.Now()
	sections := mockdata.MockStatusSections(now)
	issues := mockdata.MockCatalog(now)
	board := mockdata.MockBoard()
	model := NewModel(0, 0, sections, issues, board)
	return tea.NewProgram(model, tea.WithAltScreen())
}
