package mock

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNewModelReturnsPointerModel(t *testing.T) {
	t.Parallel()

	model := NewModel(120, 40, nil, nil, nil)
	if _, ok := any(model).(*AppModel); !ok {
		t.Fatalf("NewModel should return *AppModel, got %T", model)
	}

	var _ tea.Model = model
}

func TestNewStylesUseAdaptiveColors(t *testing.T) {
	t.Parallel()

	styles := NewStyles()
	for name, color := range map[string]lipgloss.TerminalColor{
		"Statusline": styles.Statusline.GetForeground(),
		"Footer":     styles.Footer.GetForeground(),
	} {
		if _, ok := color.(lipgloss.AdaptiveColor); !ok {
			t.Fatalf("%s should use lipgloss.AdaptiveColor, got %T", name, color)
		}
	}
}
