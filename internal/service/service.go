package service

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/database"
)

// Context carries shared dependencies (the DB handle) across the application layer.
type Context struct {
	db *sql.DB
}

func NewContext() *Context {
	return &Context{db: database.DB}
}

// DB exposes the underlying database handle.
func (c *Context) DB() *sql.DB { return c.db }

// HealthCheck reports the connection status of infrastructure dependencies.
func (c *Context) HealthCheck() map[string]string {
	status := map[string]string{}
	if c.db != nil {
		if err := c.db.Ping(); err != nil {
			status["database"] = "unhealthy: " + err.Error()
		} else {
			status["database"] = "healthy"
		}
	} else {
		status["database"] = "not connected"
	}
	return status
}
