package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ml0-1337/claude-gate/internal/ui/styles"
)

// Model represents the dashboard state
type Model struct {
	width     int
	height    int
	ready     bool
	paused    bool
	
	// Components
	stats      *RequestStats
	requestLog *RequestLog
	viewport   viewport.Model
	
	// Server info
	serverURL   string
	startTime   time.Time
	oauthStatus string
	
	// UI state
	showHelp     bool
	selectedPane int // 0: stats, 1: requests
	
	// Update channel
	eventChan chan RequestEvent
}

// New creates a new dashboard model
func New(serverURL string) *Model {
	return &Model{
		stats:       NewRequestStats(),
		requestLog:  NewRequestLog(1000),
		serverURL:   serverURL,
		startTime:   time.Now(),
		oauthStatus: "Ready",
		eventChan:   make(chan RequestEvent, 100),
	}
}

// Init initializes the dashboard
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.listenForEvents(),
		tickCmd(),
	)
}

// Update handles dashboard updates
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
		case key.Matches(msg, keys.Pause):
			m.paused = !m.paused
		case key.Matches(msg, keys.Clear):
			m.requestLog.Clear()
		case key.Matches(msg, keys.Tab):
			m.selectedPane = (m.selectedPane + 1) % 2
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport
		headerHeight := 10 // Approximate header height
		footerHeight := 3
		m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
		m.viewport.SetContent(m.renderRequests())
		m.ready = true

	case tickMsg:
		// Update UI periodically
		if !m.paused {
			m.viewport.SetContent(m.renderRequests())
		}
		cmds = append(cmds, tickCmd())

	case RequestEvent:
		// Record the request
		if !m.paused {
			m.stats.RecordRequest(msg.StatusCode, msg.Duration)
			m.requestLog.Add(msg)
			m.viewport.SetContent(m.renderRequests())
		}
		// Continue listening for more events
		cmds = append(cmds, m.listenForEvents())
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the dashboard
func (m *Model) View() string {
	if !m.ready {
		return "Initializing dashboard..."
	}

	var s strings.Builder

	// Header
	s.WriteString(m.renderHeader())
	s.WriteString("\n\n")

	// Stats panel
	s.WriteString(m.renderStats())
	s.WriteString("\n\n")

	// Request log viewport
	s.WriteString(m.renderRequestsHeader())
	s.WriteString("\n")
	s.WriteString(m.viewport.View())
	s.WriteString("\n")

	// Footer
	s.WriteString(m.renderFooter())

	return s.String()
}

// renderHeader renders the dashboard header
func (m *Model) renderHeader() string {
	uptime := time.Since(m.startTime).Round(time.Second)
	
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Render("🚀 Claude Gate Dashboard")
	
	status := lipgloss.NewStyle().
		Foreground(styles.Success).
		Render("● Running")
	
	if m.paused {
		status = lipgloss.NewStyle().
			Foreground(styles.Warning).
			Render("⏸ Paused")
	}
	
	info := fmt.Sprintf("Server: %s | OAuth: %s | Uptime: %s",
		m.serverURL, m.oauthStatus, uptime)
	
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		title,
		strings.Repeat(" ", 4),
		status,
	)
	
	return header + "\n" + styles.DescriptionStyle.Render(info)
}

// renderStats renders the statistics panel
func (m *Model) renderStats() string {
	stats := m.stats.GetStats()
	
	// Create stat cards
	cards := []string{
		m.createStatCard("Total Requests", fmt.Sprintf("%d", stats.TotalRequests), styles.InfoStyle),
		m.createStatCard("Success Rate", fmt.Sprintf("%.1f%%", m.calculateSuccessRate(stats)), styles.SuccessStyle),
		m.createStatCard("Avg Response", stats.AvgDuration.Round(time.Millisecond).String(), styles.InfoStyle),
		m.createStatCard("Requests/sec", formatReqPerSecond(stats.ReqPerSecond), styles.InfoStyle),
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, cards...)
}

// createStatCard creates a single stat card
func (m *Model) createStatCard(label, value string, valueStyle lipgloss.Style) string {
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(0, 1).
		Width(20).
		Height(3)
	
	content := fmt.Sprintf("%s\n%s",
		styles.DescriptionStyle.Render(label),
		valueStyle.Render(value),
	)
	
	return card.Render(content)
}

// renderRequestsHeader renders the request log header
func (m *Model) renderRequestsHeader() string {
	title := styles.SubtitleStyle.Render("Recent Requests")
	header := styles.DescriptionStyle.Render("Time     Method Status Duration Size Path")
	
	return title + "\n" + header
}

// renderRequests renders the request log content
func (m *Model) renderRequests() string {
	requests := m.requestLog.GetRequests(100)
	
	if len(requests) == 0 {
		return styles.DescriptionStyle.Render("No requests yet...")
	}
	
	var lines []string
	for _, req := range requests {
		line := FormatRequest(req)
		
		// Color based on status
		var style lipgloss.Style
		if req.StatusCode >= 200 && req.StatusCode < 300 {
			style = styles.SuccessStyle
		} else if req.StatusCode >= 400 && req.StatusCode < 500 {
			style = styles.WarningStyle
		} else if req.StatusCode >= 500 {
			style = styles.ErrorStyle
		} else {
			style = styles.InfoStyle
		}
		
		lines = append(lines, style.Render(line))
	}
	
	return strings.Join(lines, "\n")
}

// renderFooter renders the dashboard footer
func (m *Model) renderFooter() string {
	help := []string{
		"q: quit",
		"p: pause",
		"c: clear",
		"?: help",
		"↑/↓: scroll",
	}
	
	if m.showHelp {
		help = append(help, []string{
			"tab: switch panes",
			"f: filter",
			"e: export",
		}...)
	}
	
	return styles.HelpStyle.Render(strings.Join(help, " • "))
}

// calculateSuccessRate calculates the success rate from stats
func (m *Model) calculateSuccessRate(stats Stats) float64 {
	if stats.TotalRequests == 0 {
		return 0
	}
	return float64(stats.SuccessCount) / float64(stats.TotalRequests) * 100
}

// formatReqPerSecond formats the requests per second value
func formatReqPerSecond(rate float64) string {
	if rate < 0.1 {
		// For very low rates, show 3 decimal places
		return fmt.Sprintf("%.3f", rate)
	} else if rate < 1.0 {
		// For rates less than 1, show 2 decimal places
		return fmt.Sprintf("%.2f", rate)
	} else {
		// For higher rates, show 1 decimal place
		return fmt.Sprintf("%.1f", rate)
	}
}

// listenForEvents listens for request events
func (m *Model) listenForEvents() tea.Cmd {
	return func() tea.Msg {
		event := <-m.eventChan
		return event
	}
}

// tickMsg is sent periodically to update the UI
type tickMsg time.Time

// tickCmd returns a command that sends a tick every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Key bindings
type keyMap struct {
	Quit  key.Binding
	Help  key.Binding
	Pause key.Binding
	Clear key.Binding
	Tab   key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Pause: key.NewBinding(
		key.WithKeys("p", " "),
		key.WithHelp("p", "pause"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch pane"),
	),
}

// SendEvent sends a request event to the dashboard
func (m *Model) SendEvent(event RequestEvent) {
	select {
	case m.eventChan <- event:
	default:
		// Channel full, drop event
	}
}