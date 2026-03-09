package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"mini-asm/internal/model"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {

    response := model.HealthCheckResponse{
        Status:    "ok",                    
        Timestamp: time.Now().UTC(),
    }

    dbStatus := "disconnected"
    statusCode := http.StatusServiceUnavailable

    // Check database connection
    if h.db != nil && h.db.Ping() == nil {
        dbStatus = "connected"
        statusCode = http.StatusOK
        
        // Only get stats if database is connected
        stats := h.db.Stats()
        response.Database = model.DatabaseStats{
            Status:          dbStatus,
            OpenConnections: stats.OpenConnections,
            InUse:           stats.InUse,
            Idle:            stats.Idle,
            MaxOpen:         stats.MaxOpenConnections,
        }
    } else {
        // Database disconnected - return minimal info
        response.Status = "degraded"
        statusCode = http.StatusServiceUnavailable
        response.Database = model.DatabaseStats{
            Status:          dbStatus,
            OpenConnections: 0,
            InUse:           0,
            Idle:            0,
            MaxOpen:         0,
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}

/*
🎓 NOTES:

Refactored từ Buổi 1:
- Buổi 1: Health check logic trong main.go
- Buổi 2: Extracted to separate handler

Benefits:
- Consistent with other handlers
- Can add more health checks (database, etc.) in Buổi 3
- Reusable and testable
*/
