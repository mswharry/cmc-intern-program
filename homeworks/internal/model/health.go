package model

import "time"

type DatabaseStats struct {
    Status           string `json:"status"`
    OpenConnections  int    `json:"open_connections"`
    InUse            int    `json:"in_use"`
    Idle             int    `json:"idle"`
    MaxOpen          int    `json:"max_open"`
}

type HealthCheckResponse struct {
    Status    string         `json:"status"`
    Database  DatabaseStats  `json:"database"`
    Timestamp time.Time      `json:"timestamp"`
}