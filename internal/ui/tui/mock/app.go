package mock

import (
	"time"

	flexbox "github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen identifies which buffer is currently active in the mock TUI.
type Screen int

const (
	ScreenStatus Screen = iota
	ScreenIssues
	ScreenDetail
	ScreenKanban
)

// AppModel hosts the mocked Bubbletea state used to exercise layouts.
type AppModel struct {
	screen Screen
	ready  bool

	width  int
	height int

	sections []mockStatusSection
	issues   []mockIssue
	board    []mockBoardColumn

	sectionIndex int
	issueIndex   int

	profile layoutProfile

	keyPrefix rune

	styles Styles
}

func NewModel(width, height int, sections []mockStatusSection, issues []mockIssue, board []mockBoardColumn) AppModel {
	prof := profileForWidth(width)
	return AppModel{
		screen:       ScreenStatus,
		width:        width,
		height:       height,
		sections:     sections,
		issues:       issues,
		board:        board,
		profile:      prof,
		styles:       newStyles(),
		sectionIndex: 0,
		issueIndex:   0,
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg {
		return readyMsg{}
	})
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.profile = profileForWidth(msg.Width)
		return m, nil
	case readyMsg:
		m.ready = true
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m AppModel) View() string {
	box := flexbox.New(m.width, m.height)
	headerRow := box.NewRow()
	headerRow.AddCells(flexbox.NewCell(1, 1).SetContent(m.styles.Statusline.Render("Hubless Mock · " + m.profile.Name())))
	bodyRow := box.NewRow()
	bodyRow.AddCells(flexbox.NewCell(1, 6).SetContent("mocked body"))
	footerRow := box.NewRow()
	footerRow.AddCells(flexbox.NewCell(1, 1).SetContent(m.styles.Footer.Render("Press q to quit")))
	box.SetRows([]*flexbox.Row{headerRow, bodyRow, footerRow})
	return box.Render()
}

func (m AppModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

type readyMsg struct{}

// Supporting mock-friendly copies of data structures to avoid coupling to domain package.
type mockStatusSection struct {
	Title   string
	Items   []string
	Hint    string
	Counter int
}

type mockIssue struct {
	ID          string
	Title       string
	Status      string
	Priority    string
	Assignee    string
	LastUpdated time.Time
	Body        string
	Comments    []mockComment
	Events      []mockEvent
}

type mockComment struct {
	Author    string
	Body      string
	CreatedAt time.Time
}

type mockEvent struct {
	Label     string
	Actor     string
	Timestamp time.Time
	Note      string
}

type mockBoardColumn struct {
	Name      string
	Limit     int
	Issues    []mockBoardCard
	Highlight string
}

type mockBoardCard struct {
	ID       string
	Title    string
	Assignee string
	Priority string
}
