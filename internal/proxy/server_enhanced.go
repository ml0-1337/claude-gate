package proxy

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/yourusername/claude-gate/internal/auth"
	"github.com/yourusername/claude-gate/internal/ui/dashboard"
)

// EnhancedProxyServer extends ProxyServer with dashboard functionality
type EnhancedProxyServer struct {
	*ProxyServer
	dashboard *dashboard.Model
}

// NewEnhancedProxyServer creates a proxy server with dashboard
func NewEnhancedProxyServer(config *ProxyConfig, address string, storage *auth.TokenStorage) *EnhancedProxyServer {
	// Create base proxy server components
	handler := NewProxyHandler(config)
	healthHandler := NewHealthHandler(storage)
	
	// Create dashboard
	dashboardModel := dashboard.New(fmt.Sprintf("http://%s", address))
	
	// Create middleware that logs to dashboard
	middleware := &dashboardMiddleware{
		handler:   CreateMux(handler, healthHandler),
		dashboard: dashboardModel,
	}
	
	// Create enhanced server
	server := &http.Server{
		Addr:         address,
		Handler:      middleware,
		ReadTimeout:  time.Minute,
		WriteTimeout: 10 * time.Minute, // Long timeout for streaming
	}
	
	return &EnhancedProxyServer{
		ProxyServer: &ProxyServer{
			handler: handler,
			server:  server,
		},
		dashboard: dashboardModel,
	}
}

// GetDashboard returns the dashboard model
func (s *EnhancedProxyServer) GetDashboard() *dashboard.Model {
	return s.dashboard
}

// dashboardMiddleware logs requests to the dashboard
type dashboardMiddleware struct {
	handler   http.Handler
	dashboard *dashboard.Model
}

func (m *dashboardMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skip health checks and root endpoint
	if r.URL.Path == "/health" || r.URL.Path == "/" {
		m.handler.ServeHTTP(w, r)
		return
	}
	
	// Create a response writer wrapper to capture status and size
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	start := time.Now()
	
	// Serve the request
	m.handler.ServeHTTP(rw, r)
	
	// Log to dashboard
	duration := time.Since(start)
	event := dashboard.RequestEvent{
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: rw.statusCode,
		Duration:   duration,
		Timestamp:  start,
		Size:       rw.written,
	}
	
	m.dashboard.SendEvent(event)
}

// responseWriter wraps http.ResponseWriter to capture response details
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}