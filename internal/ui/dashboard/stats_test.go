package dashboard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestStats_SingleRequest(t *testing.T) {
	stats := NewRequestStats()
	
	// Record a single request
	stats.RecordRequest(200, 100*time.Millisecond)
	
	result := stats.GetStats()
	
	assert.Equal(t, int64(1), result.TotalRequests)
	assert.Equal(t, int64(1), result.SuccessCount)
	assert.Equal(t, int64(0), result.ErrorCount)
	assert.Equal(t, 100*time.Millisecond, result.AvgDuration)
	
	// Check that req/sec is NOT zero for a single request
	assert.Greater(t, result.ReqPerSecond, 0.0, "ReqPerSecond should be greater than 0 for a single request")
	
	// For a 60-second window, 1 request should give ~0.0167 req/s
	expectedRate := 1.0 / 60.0
	assert.InDelta(t, expectedRate, result.ReqPerSecond, 0.001, "Rate should be approximately 1/60 for a single request")
}

func TestRequestStats_MultipleRequests(t *testing.T) {
	stats := NewRequestStats()
	
	// Record multiple requests quickly
	stats.RecordRequest(200, 100*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	stats.RecordRequest(200, 150*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	stats.RecordRequest(200, 200*time.Millisecond)
	
	result := stats.GetStats()
	
	assert.Equal(t, int64(3), result.TotalRequests)
	assert.Equal(t, int64(3), result.SuccessCount)
	
	// For 3 requests over ~0.2 seconds, rate should be around 15 req/s
	// But with minimum timespan of 1 second, it should be 3 req/s
	assert.Greater(t, result.ReqPerSecond, 0.0)
	assert.LessOrEqual(t, result.ReqPerSecond, 3.0, "Rate should not exceed 3 req/s with 1-second minimum")
}

func TestRequestStats_ErrorRequests(t *testing.T) {
	stats := NewRequestStats()
	
	stats.RecordRequest(200, 100*time.Millisecond)
	stats.RecordRequest(404, 50*time.Millisecond)
	stats.RecordRequest(500, 200*time.Millisecond)
	
	result := stats.GetStats()
	
	assert.Equal(t, int64(3), result.TotalRequests)
	assert.Equal(t, int64(1), result.SuccessCount)
	assert.Equal(t, int64(2), result.ErrorCount)
}

func TestRequestStats_WindowCleaning(t *testing.T) {
	stats := &RequestStats{
		recentRequests: make([]time.Time, 0, 1000),
		windowSize:     5 * time.Second, // Short window for testing
		lastUpdate:     time.Now(),
	}
	
	// Add a request
	stats.RecordRequest(200, 100*time.Millisecond)
	
	// Check initial state
	result1 := stats.GetStats()
	assert.Greater(t, result1.ReqPerSecond, 0.0)
	
	// Wait for window to expire
	time.Sleep(6 * time.Second)
	
	// Add another request
	stats.RecordRequest(200, 100*time.Millisecond)
	
	// Should have cleaned old request but still show rate for new one
	result2 := stats.GetStats()
	assert.Greater(t, result2.ReqPerSecond, 0.0)
}

func TestRequestStats_ReqPerSecondCalculation(t *testing.T) {
	tests := []struct {
		name           string
		windowSize     time.Duration
		requests       int
		timeBetween    time.Duration
		expectedMinRate float64
		expectedMaxRate float64
	}{
		{
			name:           "Single request in 60s window",
			windowSize:     60 * time.Second,
			requests:       1,
			timeBetween:    0,
			expectedMinRate: 0.016,  // 1/60
			expectedMaxRate: 0.018,
		},
		{
			name:           "Two requests 1 second apart",
			windowSize:     60 * time.Second,
			requests:       2,
			timeBetween:    1 * time.Second,
			expectedMinRate: 1.9,    // ~2 requests / 1 second
			expectedMaxRate: 2.1,
		},
		{
			name:           "Five requests 100ms apart",
			windowSize:     60 * time.Second,
			requests:       5,
			timeBetween:    100 * time.Millisecond,
			expectedMinRate: 4.5,    // 5 requests / 1 second (minimum timespan)
			expectedMaxRate: 5.5,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &RequestStats{
				recentRequests: make([]time.Time, 0, 1000),
				windowSize:     tt.windowSize,
				lastUpdate:     time.Now(),
			}
			
			// Record requests
			for i := 0; i < tt.requests; i++ {
				stats.RecordRequest(200, 100*time.Millisecond)
				if i < tt.requests-1 {
					time.Sleep(tt.timeBetween)
				}
			}
			
			result := stats.GetStats()
			
			assert.GreaterOrEqual(t, result.ReqPerSecond, tt.expectedMinRate,
				"Rate should be at least %f, got %f", tt.expectedMinRate, result.ReqPerSecond)
			assert.LessOrEqual(t, result.ReqPerSecond, tt.expectedMaxRate,
				"Rate should be at most %f, got %f", tt.expectedMaxRate, result.ReqPerSecond)
		})
	}
}

// Test to debug the exact issue
func TestRequestStats_DebugZeroRate(t *testing.T) {
	stats := NewRequestStats()
	
	// Manually inspect what happens
	t.Logf("Initial state: recentRequests=%d, windowSize=%v", 
		len(stats.recentRequests), stats.windowSize)
	
	// Record one request
	stats.RecordRequest(200, 3601*time.Millisecond)
	
	// Get the stats
	result := stats.GetStats()
	
	t.Logf("After request: totalRequests=%d, recentRequests=%d, reqPerSecond=%f",
		result.TotalRequests, len(stats.recentRequests), result.ReqPerSecond)
	
	// This MUST not be zero
	assert.NotEqual(t, 0.0, result.ReqPerSecond, 
		"ReqPerSecond is 0.0! recentRequests length: %d", len(stats.recentRequests))
}