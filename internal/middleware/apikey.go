package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/repository"
)

// APIKeyAuth authenticates requests that carry an X-API-Key header. It must be
// registered BEFORE AuthMiddleware on routes that accept either credential.
// On success it injects organizationID into the context. It does not abort the
// chain if the header is absent, so a downstream JWT middleware can still run.
func APIKeyAuth(repo *repository.APIKeyRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-API-Key")
		if raw == "" {
			c.Next()
			return
		}
		key, err := repo.FindByRawKey(raw)
		if err != nil || key.Status != "active" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		_ = repo.TouchLastUsed(key.ID)
		c.Set("organizationID", key.OrganizationID)
		c.Set("authMethod", "apikey")
		c.Next()
	}
}
