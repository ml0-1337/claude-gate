package dashboard

import (
	"testing"
	"time"

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