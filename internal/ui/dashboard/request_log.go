package dashboard

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RequestEvent represents a single request
type RequestEvent struct {
	ID         string
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	Timestamp  time.Time
	Error      string
	Size       int64
}

// RequestLog maintains a ring buffer of recent requests
type RequestLog struct {
	mu       sync.RWMutex
	requests []RequestEvent
	maxSize  int
	head     int
	count    int
	filter   RequestFilter
}

// RequestFilter defines filtering criteria
type RequestFilter struct {
	StatusCode   int
	Method       string
	PathContains string
	ShowErrors   bool
}

// NewRequestLog creates a new request log
func NewRequestLog(maxSize int) *RequestLog {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &RequestLog{
		requests: make([]RequestEvent, maxSize),
		maxSize:  maxSize,
	}
}

// Add adds a new request to the log
func (l *RequestLog) Add(event RequestEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.requests[l.head] = event
	l.head = (l.head + 1) % l.maxSize
	if l.count < l.maxSize {
		l.count++
	}
}

// GetRequests returns filtered requests
func (l *RequestLog) GetRequests(limit int) []RequestEvent {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if limit <= 0 || limit > l.count {
		limit = l.count
	}
	
	result := make([]RequestEvent, 0, limit)
	
	// Start from the most recent
	start := l.head - 1
	if start < 0 {
		start = l.maxSize - 1
	}
	
	for i := 0; i < l.count && len(result) < limit; i++ {
		idx := (start - i + l.maxSize) % l.maxSize
		req := l.requests[idx]
		
		// Apply filter
		if l.matchesFilter(req) {
			result = append(result, req)
		}
	}
	
	return result
}

// SetFilter updates the filter
func (l *RequestLog) SetFilter(filter RequestFilter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.filter = filter
}

// matchesFilter checks if a request matches the current filter
func (l *RequestLog) matchesFilter(req RequestEvent) bool {
	// Status code filter
	if l.filter.StatusCode > 0 && req.StatusCode != l.filter.StatusCode {
		return false
	}
	
	// Method filter
	if l.filter.Method != "" && req.Method != l.filter.Method {
		return false
	}
	
	// Path filter
	if l.filter.PathContains != "" && !strings.Contains(req.Path, l.filter.PathContains) {
		return false
	}
	
	// Error filter
	if l.filter.ShowErrors && req.Error == "" {
		return false
	}
	
	return true
}

// Clear clears the log
func (l *RequestLog) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.head = 0
	l.count = 0
	l.requests = make([]RequestEvent, l.maxSize)
}

// FormatRequest formats a request for display
func FormatRequest(req RequestEvent) string {
	// Status with color indicator
	var statusStr string
	if req.StatusCode >= 200 && req.StatusCode < 300 {
		statusStr = fmt.Sprintf("2%02d", req.StatusCode%100)
	} else if req.StatusCode >= 300 && req.StatusCode < 400 {
		statusStr = fmt.Sprintf("3%02d", req.StatusCode%100)
	} else if req.StatusCode >= 400 && req.StatusCode < 500 {
		statusStr = fmt.Sprintf("4%02d", req.StatusCode%100)
	} else {
		statusStr = fmt.Sprintf("5%02d", req.StatusCode%100)
	}
	
	// Format duration
	durStr := fmt.Sprintf("%4dms", req.Duration.Milliseconds())
	
	// Format size
	sizeStr := formatBytes(req.Size)
	
	// Timestamp
	timeStr := req.Timestamp.Format("15:04:05")
	
	// Path (truncate if needed)
	path := req.Path
	if len(path) > 40 {
		path = path[:37] + "..."
	}
	
	return fmt.Sprintf("%s %s %s %s %s %-40s",
		timeStr, req.Method, statusStr, durStr, sizeStr, path)
}

// formatBytes formats bytes into human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%4dB ", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%3.0f%cB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}