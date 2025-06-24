package dashboard

import (
	"fmt"
	"sync"
	"time"
)

// RequestStats holds statistics about requests
type RequestStats struct {
	mu            sync.RWMutex
	totalRequests int64
	successCount  int64
	errorCount    int64
	avgDuration   time.Duration
	reqPerSecond  float64
	lastUpdate    time.Time
	
	// Time buckets for rate calculation
	recentRequests []time.Time
	windowSize     time.Duration
}

// NewRequestStats creates a new statistics tracker
func NewRequestStats() *RequestStats {
	return &RequestStats{
		recentRequests: make([]time.Time, 0, 1000),
		windowSize:     time.Minute,
		lastUpdate:     time.Now(),
	}
}

// RecordRequest records a new request
func (s *RequestStats) RecordRequest(statusCode int, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	s.totalRequests++
	
	if statusCode >= 200 && statusCode < 400 {
		s.successCount++
	} else if statusCode >= 400 {
		s.errorCount++
	}
	
	// Update average duration
	if s.avgDuration == 0 {
		s.avgDuration = duration
	} else {
		// Simple moving average
		s.avgDuration = (s.avgDuration + duration) / 2
	}
	
	// Add to recent requests
	s.recentRequests = append(s.recentRequests, now)
	
	// Clean old requests outside window
	cutoff := now.Add(-s.windowSize)
	i := 0
	for i < len(s.recentRequests) && s.recentRequests[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		s.recentRequests = s.recentRequests[i:]
	}
	
	// Calculate requests per second
	if len(s.recentRequests) > 1 {
		timeSpan := now.Sub(s.recentRequests[0]).Seconds()
		if timeSpan > 0 {
			s.reqPerSecond = float64(len(s.recentRequests)) / timeSpan
		}
	}
	
	s.lastUpdate = now
}

// GetStats returns current statistics
func (s *RequestStats) GetStats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return Stats{
		TotalRequests: s.totalRequests,
		SuccessCount:  s.successCount,
		ErrorCount:    s.errorCount,
		AvgDuration:   s.avgDuration,
		ReqPerSecond:  s.reqPerSecond,
		LastUpdate:    s.lastUpdate,
	}
}

// Stats represents a snapshot of statistics
type Stats struct {
	TotalRequests int64
	SuccessCount  int64
	ErrorCount    int64
	AvgDuration   time.Duration
	ReqPerSecond  float64
	LastUpdate    time.Time
}

// String returns a formatted string of the stats
func (s Stats) String() string {
	successRate := float64(0)
	if s.TotalRequests > 0 {
		successRate = float64(s.SuccessCount) / float64(s.TotalRequests) * 100
	}
	
	return fmt.Sprintf(
		"Total: %d | Success: %.1f%% | Avg: %s | Rate: %.1f req/s",
		s.TotalRequests,
		successRate,
		s.AvgDuration.Round(time.Millisecond),
		s.ReqPerSecond,
	)
}