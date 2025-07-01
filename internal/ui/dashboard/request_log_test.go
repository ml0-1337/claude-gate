package dashboard

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 11: NewRequestLog creates log with correct capacity
func TestNewRequestLog_Initialization(t *testing.T) {
	// Prediction: This test will pass - testing initialization
	
	t.Run("with positive size", func(t *testing.T) {
		log := NewRequestLog(100)
		assert.NotNil(t, log)
		assert.Equal(t, 100, log.maxSize)
		assert.Equal(t, 100, len(log.requests))
		assert.Equal(t, 0, log.head)
		assert.Equal(t, 0, log.count)
	})
	
	t.Run("with zero size defaults to 1000", func(t *testing.T) {
		log := NewRequestLog(0)
		assert.Equal(t, 1000, log.maxSize)
		assert.Equal(t, 1000, len(log.requests))
	})
	
	t.Run("with negative size defaults to 1000", func(t *testing.T) {
		log := NewRequestLog(-10)
		assert.Equal(t, 1000, log.maxSize)
	})
}

// Test 12: Add request appends to log correctly
func TestRequestLog_Add(t *testing.T) {
	// Prediction: This test will pass - testing add functionality
	
	log := NewRequestLog(5)
	
	event := RequestEvent{
		ID:         "req-1",
		Method:     "GET",
		Path:       "/api/users",
		StatusCode: 200,
		Duration:   100 * time.Millisecond,
		Timestamp:  time.Now(),
		Size:       1024,
	}
	
	log.Add(event)
	
	assert.Equal(t, 1, log.count)
	assert.Equal(t, 1, log.head)
	assert.Equal(t, event.ID, log.requests[0].ID)
}

// Test 13: Add respects max capacity (circular buffer)
func TestRequestLog_CircularBuffer(t *testing.T) {
	// Prediction: This test will pass - testing circular buffer behavior
	
	log := NewRequestLog(3)
	
	// Add 5 events to a buffer of size 3
	for i := 0; i < 5; i++ {
		event := RequestEvent{
			ID:         fmt.Sprintf("req-%d", i),
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v%d", i),
			StatusCode: 200,
			Timestamp:  time.Now(),
		}
		log.Add(event)
	}
	
	// Should only have 3 items (max capacity)
	assert.Equal(t, 3, log.count)
	assert.Equal(t, 2, log.head) // (5 % 3) = 2
	
	// Verify we have the last 3 items (req-2, req-3, req-4)
	requests := log.GetRequests(10)
	assert.Len(t, requests, 3)
	assert.Equal(t, "req-4", requests[0].ID) // Most recent
	assert.Equal(t, "req-3", requests[1].ID)
	assert.Equal(t, "req-2", requests[2].ID)
}

// Test 14: GetRequests returns all entries
func TestRequestLog_GetRequests(t *testing.T) {
	// Prediction: This test will pass - testing retrieval
	
	log := NewRequestLog(10)
	
	// Add some events
	for i := 0; i < 5; i++ {
		event := RequestEvent{
			ID:         fmt.Sprintf("req-%d", i),
			Method:     "GET",
			Path:       "/api/test",
			StatusCode: 200,
			Timestamp:  time.Now().Add(time.Duration(i) * time.Second),
		}
		log.Add(event)
	}
	
	t.Run("get all requests", func(t *testing.T) {
		requests := log.GetRequests(0) // 0 means get all
		assert.Len(t, requests, 5)
		// Most recent first
		assert.Equal(t, "req-4", requests[0].ID)
		assert.Equal(t, "req-0", requests[4].ID)
	})
	
	t.Run("get limited requests", func(t *testing.T) {
		requests := log.GetRequests(3)
		assert.Len(t, requests, 3)
		assert.Equal(t, "req-4", requests[0].ID)
		assert.Equal(t, "req-3", requests[1].ID)
		assert.Equal(t, "req-2", requests[2].ID)
	})
	
	t.Run("get more than available", func(t *testing.T) {
		requests := log.GetRequests(100)
		assert.Len(t, requests, 5) // Only 5 available
	})
}

// Test 15: GetRequests filters by status code
func TestRequestLog_FilterByStatusCode(t *testing.T) {
	// Prediction: This test will pass - testing status code filtering
	
	log := NewRequestLog(10)
	
	// Add events with different status codes
	statuses := []int{200, 404, 200, 500, 201}
	for i, status := range statuses {
		event := RequestEvent{
			ID:         fmt.Sprintf("req-%d", i),
			Method:     "GET",
			Path:       "/api/test",
			StatusCode: status,
			Timestamp:  time.Now(),
		}
		log.Add(event)
	}
	
	// Filter for 200 status
	log.SetFilter(RequestFilter{StatusCode: 200})
	requests := log.GetRequests(10)
	
	assert.Len(t, requests, 2)
	for _, req := range requests {
		assert.Equal(t, 200, req.StatusCode)
	}
}

// Test 16: GetRequests filters by path
func TestRequestLog_FilterByPath(t *testing.T) {
	// Prediction: This test will pass - testing path filtering
	
	log := NewRequestLog(10)
	
	// Add events with different paths
	paths := []string{"/api/users", "/api/posts", "/health", "/api/users/123", "/metrics"}
	for i, path := range paths {
		event := RequestEvent{
			ID:         fmt.Sprintf("req-%d", i),
			Method:     "GET",
			Path:       path,
			StatusCode: 200,
			Timestamp:  time.Now(),
		}
		log.Add(event)
	}
	
	// Filter for paths containing "users"
	log.SetFilter(RequestFilter{PathContains: "users"})
	requests := log.GetRequests(10)
	
	assert.Len(t, requests, 2)
	assert.Equal(t, "/api/users/123", requests[0].Path)
	assert.Equal(t, "/api/users", requests[1].Path)
}

// Test 17: SetFilter updates filter correctly
func TestRequestLog_SetFilter(t *testing.T) {
	// Prediction: This test will pass - testing filter updates
	
	log := NewRequestLog(10)
	
	// Add various events
	events := []RequestEvent{
		{ID: "1", Method: "GET", Path: "/api/users", StatusCode: 200},
		{ID: "2", Method: "POST", Path: "/api/users", StatusCode: 201},
		{ID: "3", Method: "GET", Path: "/api/posts", StatusCode: 200},
		{ID: "4", Method: "DELETE", Path: "/api/users/1", StatusCode: 204},
		{ID: "5", Method: "GET", Path: "/api/users", StatusCode: 404, Error: "Not found"},
	}
	
	for _, event := range events {
		event.Timestamp = time.Now()
		log.Add(event)
	}
	
	t.Run("filter by method", func(t *testing.T) {
		log.SetFilter(RequestFilter{Method: "GET"})
		requests := log.GetRequests(10)
		assert.Len(t, requests, 3)
		for _, req := range requests {
			assert.Equal(t, "GET", req.Method)
		}
	})
	
	t.Run("filter by multiple criteria", func(t *testing.T) {
		log.SetFilter(RequestFilter{
			Method:       "GET",
			PathContains: "users",
			StatusCode:   200,
		})
		requests := log.GetRequests(10)
		assert.Len(t, requests, 1)
		assert.Equal(t, "1", requests[0].ID)
	})
	
	t.Run("filter for errors only", func(t *testing.T) {
		log.SetFilter(RequestFilter{ShowErrors: true})
		requests := log.GetRequests(10)
		assert.Len(t, requests, 1)
		assert.Equal(t, "5", requests[0].ID)
		assert.NotEmpty(t, requests[0].Error)
	})
	
	t.Run("clear filter", func(t *testing.T) {
		log.SetFilter(RequestFilter{}) // Empty filter
		requests := log.GetRequests(10)
		assert.Len(t, requests, 5) // All events
	})
}

// Test 18: Clear removes all entries
func TestRequestLog_Clear(t *testing.T) {
	// Prediction: This test will pass - testing clear functionality
	
	log := NewRequestLog(10)
	
	// Add some events
	for i := 0; i < 5; i++ {
		event := RequestEvent{
			ID:        fmt.Sprintf("req-%d", i),
			Method:    "GET",
			Path:      "/api/test",
			Timestamp: time.Now(),
		}
		log.Add(event)
	}
	
	assert.Equal(t, 5, log.count)
	
	// Clear the log
	log.Clear()
	
	assert.Equal(t, 0, log.count)
	assert.Equal(t, 0, log.head)
	
	// Get requests should return empty
	requests := log.GetRequests(10)
	assert.Empty(t, requests)
}

// Test 19: Concurrent access is thread-safe
func TestRequestLog_ConcurrentAccess(t *testing.T) {
	// Prediction: This test will pass - testing thread safety
	
	log := NewRequestLog(100)
	wg := sync.WaitGroup{}
	
	// Multiple writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				event := RequestEvent{
					ID:         fmt.Sprintf("writer-%d-req-%d", id, j),
					Method:     "GET",
					Path:       fmt.Sprintf("/api/test/%d", id),
					StatusCode: 200,
					Timestamp:  time.Now(),
				}
				log.Add(event)
				time.Sleep(time.Millisecond)
			}
		}(i)
	}
	
	// Multiple readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				requests := log.GetRequests(10)
				_ = len(requests) // Just access it
				time.Sleep(time.Millisecond)
			}
		}()
	}
	
	// Filter updater
	wg.Add(1)
	go func() {
		defer wg.Done()
		filters := []RequestFilter{
			{StatusCode: 200},
			{Method: "GET"},
			{PathContains: "test"},
			{},
		}
		for i := 0; i < 10; i++ {
			log.SetFilter(filters[i%len(filters)])
			time.Sleep(5 * time.Millisecond)
		}
	}()
	
	wg.Wait()
	
	// Verify we have the expected number of events
	assert.Equal(t, 100, log.count) // Buffer is full
}

// Test FormatRequest function
func TestFormatRequest(t *testing.T) {
	// Prediction: This test will pass - testing request formatting
	
	req := RequestEvent{
		Method:     "GET",
		Path:       "/api/v1/messages",
		StatusCode: 200,
		Duration:   123 * time.Millisecond,
		Timestamp:  time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC),
		Size:       2048,
	}
	
	formatted := FormatRequest(req)
	
	// Check various parts of the formatted string
	assert.Contains(t, formatted, "12:30:45")
	assert.Contains(t, formatted, "GET")
	assert.Contains(t, formatted, "200")
	assert.Contains(t, formatted, "123ms")
	assert.Contains(t, formatted, "2KB")
	assert.Contains(t, formatted, "/api/v1/messages")
	
	t.Run("long path truncation", func(t *testing.T) {
		longPath := "/api/v1/very/long/path/that/exceeds/forty/characters/limit"
		req.Path = longPath
		formatted := FormatRequest(req)
		assert.Contains(t, formatted, "...")
		// The truncated path in output should be 40 chars (37 + "...")
		assert.Greater(t, len(longPath), 40) // Original path is longer than 40
		// Check that the formatted string contains truncated path
		assert.Contains(t, formatted, longPath[:37])
	})
	
	t.Run("different status codes", func(t *testing.T) {
		testCases := []struct {
			status   int
			expected string
		}{
			{200, "200"},
			{301, "301"},
			{404, "404"},
			{500, "500"},
		}
		
		for _, tc := range testCases {
			req.StatusCode = tc.status
			formatted := FormatRequest(req)
			assert.Contains(t, formatted, tc.expected)
		}
	})
}

// Test formatBytes function
func TestFormatBytes(t *testing.T) {
	// Prediction: This test will pass - testing byte formatting
	
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{0, "   0B "},
		{100, " 100B "},
		{1024, "  1KB"},
		{1536, "  2KB"}, // 1.5KB rounds to 2
		{1048576, "  1MB"},
		{5242880, "  5MB"},
		{1073741824, "  1GB"},
		{1099511627776, "  1TB"},
	}
	
	for _, tc := range testCases {
		result := formatBytes(tc.bytes)
		assert.Equal(t, tc.expected, result, "formatBytes(%d)", tc.bytes)
	}
}

// Test edge cases
func TestRequestLog_EdgeCases(t *testing.T) {
	// Prediction: This test will pass - testing edge cases
	
	t.Run("empty log operations", func(t *testing.T) {
		log := NewRequestLog(5)
		
		// Getting from empty log
		requests := log.GetRequests(10)
		assert.Empty(t, requests)
		
		// Clear empty log
		assert.NotPanics(t, func() {
			log.Clear()
		})
	})
	
	t.Run("single item circular buffer", func(t *testing.T) {
		log := NewRequestLog(1)
		
		// Add multiple items
		for i := 0; i < 3; i++ {
			event := RequestEvent{
				ID:        fmt.Sprintf("req-%d", i),
				Method:    "GET",
				Path:      "/test",
				Timestamp: time.Now(),
			}
			log.Add(event)
		}
		
		requests := log.GetRequests(10)
		assert.Len(t, requests, 1)
		assert.Equal(t, "req-2", requests[0].ID) // Only the last one
	})
	
	t.Run("filter with no matches", func(t *testing.T) {
		log := NewRequestLog(10)
		
		// Add some events
		for i := 0; i < 5; i++ {
			event := RequestEvent{
				ID:         fmt.Sprintf("req-%d", i),
				Method:     "GET",
				Path:       "/api/test",
				StatusCode: 200,
				Timestamp:  time.Now(),
			}
			log.Add(event)
		}
		
		// Set filter that matches nothing
		log.SetFilter(RequestFilter{StatusCode: 999})
		requests := log.GetRequests(10)
		assert.Empty(t, requests)
	})
}

// Test request with all fields populated
func TestRequestEvent_AllFields(t *testing.T) {
	// Prediction: This test will pass - testing complete request event
	
	log := NewRequestLog(5)
	
	event := RequestEvent{
		ID:         "test-123",
		Method:     "POST",
		Path:       "/api/v1/users",
		StatusCode: 201,
		Duration:   456 * time.Millisecond,
		Timestamp:  time.Now(),
		Error:      "partial error",
		Size:       10240,
	}
	
	log.Add(event)
	
	retrieved := log.GetRequests(1)
	require.Len(t, retrieved, 1)
	
	assert.Equal(t, event.ID, retrieved[0].ID)
	assert.Equal(t, event.Method, retrieved[0].Method)
	assert.Equal(t, event.Path, retrieved[0].Path)
	assert.Equal(t, event.StatusCode, retrieved[0].StatusCode)
	assert.Equal(t, event.Duration, retrieved[0].Duration)
	assert.Equal(t, event.Error, retrieved[0].Error)
	assert.Equal(t, event.Size, retrieved[0].Size)
}