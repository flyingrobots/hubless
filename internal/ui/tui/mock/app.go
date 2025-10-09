package mock

import (
	"fmt"
	"strings"
	"time"

	flexbox "github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	mockdata "github.com/flyingrobots/hubless/internal/mock"
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

	sections []mockdata.StatusSection
	issues   []mockdata.Issue
	board    []mockdata.BoardColumn

	sectionIndex int
	issueIndex   int
	boardColumn  int
	boardCard    int

	keyPrefix rune
	flash     string

	filterActive bool
	filterQuery  string

	lastHydrated time.Time
	profile      layoutProfile

	styles Styles
}

// NewModel constructs a mock TUI model seeded with mocked catalog data.
func NewModel(width, height int, sections []mockdata.StatusSection, issues []mockdata.Issue, board []mockdata.BoardColumn) AppModel {
	prof := profileForWidth(width)
	return AppModel{
		screen:       ScreenStatus,
		width:        width,
		height:       height,
		sections:     sections,
		issues:       issues,
		board:        board,
		sectionIndex: 0,
		issueIndex:   0,
		boardColumn:  0,
		boardCard:    0,
		lastHydrated: time.Now(),
		profile:      prof,
		styles:       newStyles(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg {
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
	if m.width == 0 || m.height == 0 {
		return "loading layout…"
	}

	switch m.screen {
	case ScreenStatus:
		return m.renderStatus()
	case ScreenIssues:
		return m.renderIssueList()
	case ScreenDetail:
		return m.renderIssueDetail()
	case ScreenKanban:
		return m.renderKanban()
	default:
		return ""
	}
}

func (m AppModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.filterActive && (key == "esc" || key == "ctrl+c" || key == "enter") {
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		m.filterActive = false
		m.flash = "Filter overlay closed"
		return m, nil
	}

	if m.keyPrefix != 0 {
		switch m.keyPrefix {
		case 'g':
			switch key {
			case "s":
				m.screen = ScreenStatus
				m.flash = "Jumped to status"
			case "i":
				m.screen = ScreenIssues
				m.flash = "Jumped to issue list"
			case "b":
				m.screen = ScreenKanban
				m.flash = "Jumped to kanban"
			case "/":
				m.filterActive = true
				if m.filterQuery == "" {
					m.filterQuery = "status:open assignee:me priority:high"
				}
				m.flash = "Filter overlay"
			case "g":
				m.lastHydrated = time.Now()
				m.flash = "Refreshed mocked data"
			default:
				m.flash = fmt.Sprintf("Unknown go binding g %s", key)
			}
		}
		m.keyPrefix = 0
		return m, nil
	}

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "g":
		m.keyPrefix = 'g'
		m.flash = "go prefix"
		return m, nil
	}

	switch m.screen {
	case ScreenStatus:
		m = m.handleStatusKeys(key)
	case ScreenIssues:
		m = m.handleIssueKeys(key)
	case ScreenDetail:
		m = m.handleDetailKeys(key)
	case ScreenKanban:
		m = m.handleKanbanKeys(key)
	}

	return m, nil
}

func (m AppModel) handleStatusKeys(key string) AppModel {
	switch key {
	case "tab":
		if len(m.sections) > 0 {
			m.sectionIndex = (m.sectionIndex + 1) % len(m.sections)
		}
	case "shift+tab":
		if len(m.sections) > 0 {
			m.sectionIndex = (m.sectionIndex - 1 + len(m.sections)) % len(m.sections)
		}
	case "enter", "l":
		m.screen = ScreenIssues
		m.flash = "Opened issue focus"
	case "j", "down":
		if len(m.sections) > 0 {
			m.sectionIndex = clampIndex(m.sectionIndex+1, len(m.sections))
		}
	case "k", "up":
		if len(m.sections) > 0 {
			m.sectionIndex = clampIndex(m.sectionIndex-1, len(m.sections))
		}
	}
	return m
}

func (m AppModel) handleIssueKeys(key string) AppModel {
	switch key {
	case "j", "down":
		if len(m.issues) > 0 {
			m.issueIndex = clampIndex(m.issueIndex+1, len(m.issues))
		}
	case "k", "up":
		if len(m.issues) > 0 {
			m.issueIndex = clampIndex(m.issueIndex-1, len(m.issues))
		}
	case "enter", "l":
		if len(m.issues) > 0 {
			m.screen = ScreenDetail
			m.flash = fmt.Sprintf("Viewing %s", m.issues[m.issueIndex].ID)
		}
	case "b", "h", "ctrl+o":
		m.screen = ScreenStatus
		m.flash = "Returned to status"
	}
	return m
}

func (m AppModel) handleDetailKeys(key string) AppModel {
	switch key {
	case "b", "h", "ctrl+o":
		m.screen = ScreenIssues
		m.flash = "Back to list"
	case "j", "down":
		if len(m.issues) > 0 {
			m.issueIndex = clampIndex(m.issueIndex+1, len(m.issues))
		}
	case "k", "up":
		if len(m.issues) > 0 {
			m.issueIndex = clampIndex(m.issueIndex-1, len(m.issues))
		}
	case "l":
		m.screen = ScreenKanban
	case "g":
		// allow g prefix from detail state
		m.keyPrefix = 'g'
		m.flash = "go prefix"
	}
	return m
}

func (m AppModel) handleKanbanKeys(key string) AppModel {
	switch key {
	case "b", "ctrl+o":
		m.screen = ScreenStatus
		m.flash = "Back to status"
	case "enter":
		m.flash = "Popup actions are mocked"
	case "j", "down":
		if len(m.board) > 0 {
			col := m.board[m.boardColumn]
			if len(col.Issues) > 0 {
				m.boardCard = clampIndex(m.boardCard+1, len(col.Issues))
			}
		}
	case "k", "up":
		if len(m.board) > 0 {
			col := m.board[m.boardColumn]
			if len(col.Issues) > 0 {
				m.boardCard = clampIndex(m.boardCard-1, len(col.Issues))
			}
		}
	case "h", "left":
		if len(m.board) > 0 {
			m.boardColumn = clampIndex(m.boardColumn-1, len(m.board))
			m.boardCard = 0
		}
	case "l", "right":
		if len(m.board) > 0 {
			m.boardColumn = clampIndex(m.boardColumn+1, len(m.board))
			m.boardCard = 0
		}
	}
	return m
}

func (m AppModel) renderStatus() string {
	var body strings.Builder
	for idx, section := range m.sections {
		block := m.renderStatusSection(section)
		if idx == m.sectionIndex {
			body.WriteString(m.styles.SectionSel.Render(block))
		} else {
			body.WriteString(m.styles.Section.Render(block))
		}
		body.WriteString("\n\n")
	}
	content := strings.TrimRight(body.String(), "\n")
	if m.filterActive {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", m.renderFilterOverlay())
	}
	footer := "tab cycle • enter issues • g i issues • g b board"
	return m.renderFrame("Status Buffer", content, footer)
}

func (m AppModel) renderIssueList() string {
	left := m.renderIssueListPanel()
	right := m.renderIssuePreview()

	var body string
	if m.profile.twoColumn {
		innerHeight := max(0, m.height-4)
		inner := flexbox.New(m.width-2, innerHeight)
		row := inner.NewRow()
		row.AddCells(
			flexbox.NewCell(m.profile.listRatio, 1).SetContent(left),
			flexbox.NewCell(m.profile.previewRatio, 1).SetContent(right),
		)
		inner.SetRows([]*flexbox.Row{row})
		body = inner.Render()
	} else {
		body = lipgloss.JoinVertical(lipgloss.Left, left, "", right)
	}

	if m.filterActive {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", m.renderFilterOverlay())
	}

	footer := "j/k move • enter detail • b back • g s status"
	return m.renderFrame("Issues", body, footer)
}

func (m AppModel) renderIssueDetail() string {
	if len(m.issues) == 0 {
		return m.renderFrame("Issue Detail", "No issues loaded", "g s status")
	}
	issue := m.issues[m.issueIndex]

	var body strings.Builder
	headline := fmt.Sprintf("%s  %s  %s  updated %s",
		m.statusBadge(issue.Status),
		m.priorityBadge(issue.Priority),
		m.styles.Accent.Render(issue.Assignee),
		m.styles.Hint.Render(formatAgo(issue.LastUpdated)),
	)
	body.WriteString(m.styles.Title.Render(issue.Title))
	body.WriteString("\n")
	body.WriteString(headline)
	body.WriteString("\n\n")
	body.WriteString(issue.Body)

	if len(issue.Comments) > 0 {
		body.WriteString("\n\n")
		body.WriteString(m.styles.Subtitle.Render("Recent comments"))
		for _, comment := range issue.Comments {
			body.WriteString("\n")
			body.WriteString(m.styles.Key.Render(comment.Author))
			body.WriteString(" · ")
			body.WriteString(m.styles.Hint.Render(formatAgo(comment.CreatedAt)))
			body.WriteString("\n  ")
			body.WriteString(comment.Body)
		}
	}

	if len(issue.Events) > 0 {
		body.WriteString("\n\n")
		body.WriteString(m.styles.Subtitle.Render("Timeline"))
		for _, evt := range issue.Events {
			body.WriteString("\n")
			body.WriteString(m.styles.Key.Render(evt.Label))
			body.WriteString(" → ")
			body.WriteString(evt.Note)
			body.WriteString(" (" + m.styles.Hint.Render(formatAgo(evt.Timestamp)) + ")")
		}
	}

	wrapped := m.styles.Border.Render(body.String())
	if m.filterActive {
		wrapped = lipgloss.JoinVertical(lipgloss.Left, wrapped, "", m.renderFilterOverlay())
	}

	footer := "b back • g s status • g b board"
	title := fmt.Sprintf("Detail · %s", issue.ID)
	return m.renderFrame(title, wrapped, footer)
}

func (m AppModel) renderKanban() string {
	if len(m.board) == 0 {
		return m.renderFrame("Kanban", "No board data", "g s status")
	}

	canShowRow := m.width >= len(m.board)*m.profile.boardMinWidth
	var body string

	if canShowRow {
		inner := flexbox.New(m.width-2, max(0, m.height-4))
		row := inner.NewRow()
		for idx, column := range m.board {
			block := m.renderBoardColumn(column, idx == m.boardColumn)
			row.AddCells(flexbox.NewCell(1, 1).SetContent(block))
		}
		inner.SetRows([]*flexbox.Row{row})
		body = inner.Render()
	} else {
		var builder strings.Builder
		for idx, column := range m.board {
			block := m.renderBoardColumn(column, idx == m.boardColumn)
			builder.WriteString(block)
			builder.WriteString("\n\n")
		}
		body = builder.String()
	}

	if m.filterActive {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", m.renderFilterOverlay())
	}

	footer := "h/l move columns • j/k move cards • b back"
	return m.renderFrame("Kanban", body, footer)
}

func (m AppModel) renderFrame(title, body, footer string) string {
	headerContent := m.renderStatusline(title)
	footerContent := m.renderFooter(footer)

	box := flexbox.New(m.width, m.height)
	headerRow := box.NewRow()
	headerRow.AddCells(flexbox.NewCell(1, 1).SetContent(headerContent))
	bodyRow := box.NewRow()
	bodyRow.AddCells(flexbox.NewCell(1, 8).SetContent(body))
	footerRow := box.NewRow()
	footerRow.AddCells(flexbox.NewCell(1, 1).SetContent(footerContent))
	box.SetRows([]*flexbox.Row{headerRow, bodyRow, footerRow})
	return box.Render()
}

func (m AppModel) renderStatusline(title string) string {
	env := fmt.Sprintf("hubless/main · %s", m.profile.Name())
	refreshed := m.styles.Hint.Render("refreshed " + formatAgo(m.lastHydrated))
	content := fmt.Sprintf("%s · %s", env, refreshed)
	label := m.styles.Statusline.Render(content)
	return label + "\n" + m.styles.Subtitle.Render(title)
}

func (m AppModel) renderFooter(hint string) string {
	if m.flash != "" {
		hint = fmt.Sprintf("%s   %s", hint, m.styles.Accent.Render(m.flash))
	}
	return m.styles.Footer.Render(hint)
}

func (m AppModel) renderStatusSection(section mockdata.StatusSection) string {
	var builder strings.Builder
	heading := fmt.Sprintf("%s %s", m.styles.Title.Render(section.Title), m.styles.Hint.Render(fmt.Sprintf("(%d)", section.Counter)))
	builder.WriteString(heading)
	for _, item := range section.Items {
		builder.WriteString("\n  - ")
		builder.WriteString(item)
	}
	if section.Hint != "" {
		builder.WriteString("\n")
		builder.WriteString(m.styles.Hint.Render(section.Hint))
	}
	return builder.String()
}

func (m AppModel) renderIssueListPanel() string {
	if len(m.issues) == 0 {
		return "No issues in catalog"
	}
	var builder strings.Builder
	for idx, issue := range m.issues {
		prefix := "  "
		lineStyle := m.styles.Body
		if idx == m.issueIndex {
			prefix = "→ "
			lineStyle = m.styles.SectionSel
		}
		line := fmt.Sprintf("%s%-36s %s %s %s",
			prefix,
			issue.Title,
			m.statusBadge(issue.Status),
			m.priorityBadge(issue.Priority),
			m.styles.Hint.Render(formatAgo(issue.LastUpdated)),
		)
		builder.WriteString(lineStyle.Render(line))
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m AppModel) renderIssuePreview() string {
	if len(m.issues) == 0 {
		return m.styles.Border.Render("Select an issue to preview")
	}
	issue := m.issues[m.issueIndex]
	var builder strings.Builder
	builder.WriteString(m.styles.Subtitle.Render(issue.ID))
	builder.WriteString("\n")
	builder.WriteString(m.styles.Title.Render(issue.Title))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("%s  %s  %s",
		m.statusBadge(issue.Status),
		m.priorityBadge(issue.Priority),
		m.styles.Hint.Render("Updated "+formatAgo(issue.LastUpdated)),
	))
	builder.WriteString("\n\n")
	builder.WriteString(issue.Body)
	if len(issue.Comments) > 0 {
		builder.WriteString("\n\n")
		builder.WriteString(m.styles.Subtitle.Render("Comments"))
		for _, comment := range issue.Comments {
			builder.WriteString("\n")
			builder.WriteString(m.styles.Key.Render(comment.Author))
			builder.WriteString(" · ")
			builder.WriteString(m.styles.Hint.Render(formatAgo(comment.CreatedAt)))
			builder.WriteString("\n  ")
			builder.WriteString(comment.Body)
		}
	}
	return m.styles.Border.Render(builder.String())
}

func (m AppModel) renderBoardColumn(column mockdata.BoardColumn, selected bool) string {
	var builder strings.Builder
	header := fmt.Sprintf("%s (%d)", column.Name, column.Limit)
	if column.Highlight == "over" {
		header += " !"
	}
	if selected {
		builder.WriteString(m.styles.Title.Copy().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("60")).Render(header))
	} else {
		builder.WriteString(m.styles.Title.Render(header))
	}
	for idx, card := range column.Issues {
		marker := "  "
		style := m.styles.Body
		if selected && idx == m.boardCard {
			marker = "→ "
			style = m.styles.SectionSel
		}
		line := fmt.Sprintf("%s%-24s %s %s", marker, card.Title, m.styles.Hint.Render(card.Assignee), m.priorityBadge(card.Priority))
		builder.WriteString("\n")
		builder.WriteString(style.Render(line))
	}
	return m.styles.Border.Render(builder.String())
}

func (m AppModel) renderFilterOverlay() string {
	query := m.filterQuery
	if query == "" {
		query = "status:open"
	}
	return m.styles.Overlay.Render("/ " + query)
}

func (m AppModel) statusBadge(status string) string {
	switch strings.ToLower(status) {
	case "open", "planned":
		return m.styles.StatusOpen.Render(status)
	case "in-progress", "started":
		return m.styles.StatusProgress.Render(status)
	case "done", "closed":
		return m.styles.StatusClosed.Render(status)
	default:
		return m.styles.Hint.Render(status)
	}
}

func (m AppModel) priorityBadge(priority string) string {
	switch strings.ToLower(priority) {
	case "high":
		return m.styles.PriorityHigh.Render("🔥 high")
	case "medium":
		return m.styles.PriorityMedium.Render("● medium")
	case "low":
		return m.styles.PriorityLow.Render("· low")
	default:
		return m.styles.Hint.Render(priority)
	}
}

type readyMsg struct{}

func clampIndex(idx, size int) int {
	if size == 0 {
		return 0
	}
	if idx < 0 {
		return (idx%size + size) % size
	}
	if idx >= size {
		return idx % size
	}
	return idx
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func formatAgo(t time.Time) string {
	if t.IsZero() {
		return "n/a"
	}
	diff := time.Since(t)
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm", int(diff/time.Minute))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh", int(diff/time.Hour))
	}
	days := int(diff / (24 * time.Hour))
	return fmt.Sprintf("%dd", days)
}
