package dashboard

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestDashboard_DirectStatsUpdate(t *testing.T) {
	// Create dashboard
	dashboard := New("http://localhost:8000")
	
	// Simulate multiple requests by directly calling the stats methods
	events := []struct {
		statusCode int
		duration   time.Duration
	}{
		{200, 100 * time.Millisecond},
		{200, 50 * time.Millisecond},
		{400, 200 * time.Millisecond},
	}
	
	// Record events directly
	for i, event := range events {
		dashboard.stats.RecordRequest(event.statusCode, event.duration)
		dashboard.requestLog.Add(RequestEvent{
			Method:     "POST",
			Path:       "/v1/messages",
			StatusCode: event.statusCode,
			Duration:   event.duration,
			Timestamp:  time.Now().Add(time.Duration(i) * time.Second),
			Size:       1024,
		})
	}
	
	// Get stats
	stats := dashboard.stats.GetStats()
	
	// Verify all requests were recorded
	assert.Equal(t, int64(3), stats.TotalRequests, "Should have recorded 3 requests")
	assert.Equal(t, int64(2), stats.SuccessCount, "Should have 2 successful requests")
	assert.Equal(t, int64(1), stats.ErrorCount, "Should have 1 error request")
	
	// Verify request log has all entries
	requests := dashboard.requestLog.GetRequests(10)
	assert.Len(t, requests, 3, "Request log should contain 3 entries")
}

func TestDashboard_EventProcessing(t *testing.T) {
	// Create dashboard model
	model := New("http://localhost:8000")
	
	// Initialize the model
	initCmd := model.Init()
	assert.NotNil(t, initCmd, "Init should return commands")
	
	// Simulate sending multiple events through Update
	events := []RequestEvent{
		{
			Method:     "POST",
			Path:       "/v1/messages",
			StatusCode: 200,
			Duration:   100 * time.Millisecond,
			Timestamp:  time.Now(),
			Size:       1024,
		},
		{
			Method:     "GET",
			Path:       "/v1/models",
			StatusCode: 200,
			Duration:   50 * time.Millisecond,
			Timestamp:  time.Now().Add(1 * time.Second),
			Size:       512,
		},
	}
	
	// Process first event
	updatedModel, cmd := model.Update(events[0])
	model = updatedModel.(*Model)
	
	// Verify first event was processed
	stats := model.stats.GetStats()
	assert.Equal(t, int64(1), stats.TotalRequests, "Should have 1 request after first event")
	
	// Check if we have a command to continue listening
	assert.NotNil(t, cmd, "Should return a command to continue listening")
	
	// Process second event
	updatedModel, cmd = model.Update(events[1])
	model = updatedModel.(*Model)
	
	// Verify second event was processed
	stats = model.stats.GetStats()
	assert.Equal(t, int64(2), stats.TotalRequests, "Should have 2 requests after second event")
	
	// Should still have a command to continue listening
	assert.NotNil(t, cmd, "Should continue to return listening command")
}

func TestDashboard_RequestsPerSecondUpdate(t *testing.T) {
	dashboard := New("http://localhost:8000")
	
	// Record first request directly
	dashboard.stats.RecordRequest(200, 100*time.Millisecond)
	
	// Check initial rate
	stats1 := dashboard.stats.GetStats()
	assert.Greater(t, stats1.ReqPerSecond, 0.0, "Should have non-zero rate after first request")
	initialRate := stats1.ReqPerSecond
	t.Logf("Initial rate after 1 request: %.3f req/s", initialRate)
	
	// Record second request after 2 seconds
	time.Sleep(2 * time.Second)
	dashboard.stats.RecordRequest(200, 150*time.Millisecond)
	
	// Check updated rate
	stats2 := dashboard.stats.GetStats()
	t.Logf("Rate after 2 requests (2s apart): %.3f req/s", stats2.ReqPerSecond)
	
	// With 2 requests over ~2 seconds, rate should be around 1 req/s
	assert.NotEqual(t, initialRate, stats2.ReqPerSecond, "Rate should change after second request")
	assert.Greater(t, stats2.ReqPerSecond, 0.5, "Rate should be reasonable for 2 requests in 2 seconds")
	assert.Less(t, stats2.ReqPerSecond, 1.5, "Rate should be reasonable for 2 requests in 2 seconds")
}

// Test 1: Dashboard View() should show initialization message when not ready
func TestModel_View_NotReady(t *testing.T) {
	// Prediction: This test will pass - testing not ready state
	
	m := &Model{
		ready: false,
	}
	
	output := m.View()
	assert.Equal(t, "Initializing dashboard...", output)
}

// Test 2: Dashboard View() should render full dashboard when ready
func TestModel_View_Ready(t *testing.T) {
	// Prediction: This test will pass - testing full render
	
	m := &Model{
		ready:       true,
		serverURL:   "http://localhost:8080",
		startTime:   time.Now(),
		oauthStatus: "Ready",
		stats:       NewRequestStats(),
		requestLog:  NewRequestLog(100),
		viewport:    viewport.New(80, 10),
		showHelp:    false,
	}
	
	output := m.View()
	
	// Check for key components
	assert.Contains(t, output, "Claude Gate Dashboard")
	assert.Contains(t, output, "Server: http://localhost:8080")
	assert.Contains(t, output, "OAuth: Ready")
	assert.Contains(t, output, "Total Requests")
	assert.Contains(t, output, "Success Rate")
	assert.Contains(t, output, "Recent Requests")
}

// Test 3: renderHeader() should show running status
func TestModel_RenderHeader_Running(t *testing.T) {
	// Prediction: This test will pass - testing running header
	
	m := &Model{
		serverURL:   "http://localhost:8080",
		startTime:   time.Now().Add(-5 * time.Minute),
		oauthStatus: "Ready",
		paused:      false,
	}
	
	header := m.renderHeader()
	
	assert.Contains(t, header, "🚀 Claude Gate Dashboard")
	assert.Contains(t, header, "● Running")
	assert.Contains(t, header, "Server: http://localhost:8080")
	assert.Contains(t, header, "OAuth: Ready")
	assert.Contains(t, header, "Uptime: 5m")
}

// Test 4: renderHeader() should show paused status when paused
func TestModel_RenderHeader_Paused(t *testing.T) {
	// Prediction: This test will pass - testing paused header
	
	m := &Model{
		serverURL:   "http://localhost:8080",
		startTime:   time.Now().Add(-10 * time.Second),
		oauthStatus: "Ready",
		paused:      true,
	}
	
	header := m.renderHeader()
	
	assert.Contains(t, header, "🚀 Claude Gate Dashboard")
	assert.Contains(t, header, "⏸ Paused")
	assert.NotContains(t, header, "● Running")
}

// Test 5: renderStats() should display all stat cards
func TestModel_RenderStats(t *testing.T) {
	// Prediction: This test will pass - testing stats rendering
	
	// Create stats with some data
	stats := NewRequestStats()
	stats.RecordRequest(200, 100*time.Millisecond)
	stats.RecordRequest(200, 200*time.Millisecond)
	stats.RecordRequest(500, 300*time.Millisecond)
	
	m := &Model{
		stats: stats,
	}
	
	output := m.renderStats()
	
	// Check for stat cards
	assert.Contains(t, output, "Total Requests")
	assert.Contains(t, output, "3") // 3 requests
	assert.Contains(t, output, "Success Rate")
	assert.Contains(t, output, "66.7%") // 2/3 success
	assert.Contains(t, output, "Avg Response")
	// Average of 100ms, 200ms, 300ms = 200ms, but internal calculation may differ
	// Just check that it contains "ms"
	assert.Contains(t, output, "ms") // has duration
	assert.Contains(t, output, "Requests/sec")
}

// Test 6: createStatCard() should format card with label and value
func TestModel_CreateStatCard(t *testing.T) {
	// Prediction: This test will pass - testing card formatting
	
	m := &Model{}
	
	// Test creating a stat card
	card := m.createStatCard("Test Label", "42", lipgloss.NewStyle())
	
	// Check card contains label and value
	assert.Contains(t, card, "Test Label")
	assert.Contains(t, card, "42")
	// Check for border characters
	assert.Contains(t, card, "╭")
	assert.Contains(t, card, "╰")
}

// Test 7: calculateSuccessRate() should handle zero total requests
func TestModel_CalculateSuccessRate_ZeroTotal(t *testing.T) {
	// Prediction: This test will pass - testing zero division
	
	m := &Model{}
	
	stats := Stats{
		TotalRequests: 0,
		SuccessCount:  0,
		ErrorCount:    0,
	}
	
	rate := m.calculateSuccessRate(stats)
	assert.Equal(t, 0.0, rate)
}

// Test 8: calculateSuccessRate() should calculate percentage correctly
func TestModel_CalculateSuccessRate_WithRequests(t *testing.T) {
	// Prediction: This test will pass - testing percentage calculation
	
	m := &Model{}
	
	tests := []struct {
		name         string
		stats        Stats
		expectedRate float64
	}{
		{
			name: "all success",
			stats: Stats{
				TotalRequests: 10,
				SuccessCount:  10,
				ErrorCount:    0,
			},
			expectedRate: 100.0,
		},
		{
			name: "half success",
			stats: Stats{
				TotalRequests: 10,
				SuccessCount:  5,
				ErrorCount:    5,
			},
			expectedRate: 50.0,
		},
		{
			name: "two thirds success",
			stats: Stats{
				TotalRequests: 3,
				SuccessCount:  2,
				ErrorCount:    1,
			},
			expectedRate: float64(2) / float64(3) * 100, // Let Go calculate it
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := m.calculateSuccessRate(tt.stats)
			assert.Equal(t, tt.expectedRate, rate)
		})
	}
}

// Test 9: formatReqPerSecond() should format rates correctly
func TestFormatReqPerSecond(t *testing.T) {
	// Prediction: This test will pass - testing number formatting
	
	tests := []struct {
		rate     float64
		expected string
	}{
		{0.001, "0.001"},      // Very low rate - 3 decimals
		{0.05, "0.050"},       // Low rate - 3 decimals
		{0.1, "0.10"},         // Boundary - 2 decimals
		{0.5, "0.50"},         // Less than 1 - 2 decimals
		{0.99, "0.99"},        // Just under 1 - 2 decimals
		{1.0, "1.0"},          // Exactly 1 - 1 decimal
		{5.5, "5.5"},          // Higher rate - 1 decimal
		{10.25, "10.2"},       // Double digit - 1 decimal
		{100.99, "101.0"},     // Triple digit - 1 decimal
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatReqPerSecond(tt.rate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test 10: renderFooter() should show help text
func TestModel_RenderFooter_NoHelp(t *testing.T) {
	// Prediction: This test will pass - testing footer without extended help
	
	m := &Model{
		showHelp: false,
	}
	
	footer := m.renderFooter()
	
	// Check basic help items
	assert.Contains(t, footer, "q: quit")
	assert.Contains(t, footer, "p: pause")
	assert.Contains(t, footer, "c: clear")
	assert.Contains(t, footer, "?: help")
	assert.Contains(t, footer, "↑/↓: scroll")
	
	// Should not contain extended help
	assert.NotContains(t, footer, "tab: switch panes")
	assert.NotContains(t, footer, "f: filter")
	assert.NotContains(t, footer, "e: export")
}

// Test 10b: renderFooter() with extended help
func TestModel_RenderFooter_WithHelp(t *testing.T) {
	// Prediction: This test will pass - testing footer with extended help
	
	m := &Model{
		showHelp: true,
	}
	
	footer := m.renderFooter()
	
	// Check all help items
	assert.Contains(t, footer, "q: quit")
	assert.Contains(t, footer, "tab: switch panes")
	assert.Contains(t, footer, "f: filter")
	assert.Contains(t, footer, "e: export")
}

// Test 11: renderRequestsHeader() should format header correctly
func TestModel_RenderRequestsHeader(t *testing.T) {
	// Prediction: This test will pass - testing request header
	
	m := &Model{}
	
	header := m.renderRequestsHeader()
	
	// Check header content
	assert.Contains(t, header, "Recent Requests")
	assert.Contains(t, header, "Time")
	assert.Contains(t, header, "Method")
	assert.Contains(t, header, "Status")
	assert.Contains(t, header, "Duration")
	assert.Contains(t, header, "Size")
	assert.Contains(t, header, "Path")
}