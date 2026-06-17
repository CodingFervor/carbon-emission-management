package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/handler"
	"github.com/CodingFervor/carbon-emission-management/internal/middleware"
)

// New builds the Gin engine with the full route table wired to the handlers.
func New(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.APIKeyAuth(h.APIKey))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")
	{
		// Auth (public)
		api.POST("/auth/login", h.Login)
		api.POST("/auth/register", h.Register)

		// Auth (authenticated) + all protected resources
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/auth/profile", h.Profile)

			// Organizations
			auth.GET("/organizations", h.ListOrganizations)
			auth.POST("/organizations", h.CreateOrganization)
			auth.GET("/organizations/:id", h.GetOrganization)
			auth.PUT("/organizations/:id", h.UpdateOrganization)
			auth.DELETE("/organizations/:id", h.DeleteOrganization)

			// Facilities
			auth.GET("/facilities", h.ListFacilities)
			auth.POST("/facilities", h.CreateFacility)
			auth.GET("/facilities/:id", h.GetFacility)
			auth.PUT("/facilities/:id", h.UpdateFacility)
			auth.DELETE("/facilities/:id", h.DeleteFacility)
			auth.GET("/organizations/:id/facilities", h.ListFacilitiesByOrg)

			// Emission Sources
			auth.GET("/emission-sources", h.ListEmissionSources)
			auth.POST("/emission-sources", h.CreateEmissionSource)
			auth.GET("/emission-sources/:id", h.GetEmissionSource)
			auth.PUT("/emission-sources/:id", h.UpdateEmissionSource)
			auth.DELETE("/emission-sources/:id", h.DeleteEmissionSource)

			// Emission Factors
			auth.GET("/emission-factors", h.ListEmissionFactors)
			auth.POST("/emission-factors", h.CreateEmissionFactor)
			auth.GET("/emission-factors/:id", h.GetEmissionFactor)
			auth.PUT("/emission-factors/:id", h.UpdateEmissionFactor)
			auth.DELETE("/emission-factors/:id", h.DeleteEmissionFactor)

			// Emission Records
			auth.GET("/emission-records", h.ListEmissionRecords)
			auth.POST("/emission-records", h.CreateEmissionRecord)
			auth.GET("/emission-records/:id", h.GetEmissionRecord)
			auth.PUT("/emission-records/:id", h.UpdateEmissionRecord)
			auth.DELETE("/emission-records/:id", h.DeleteEmissionRecord)
			auth.GET("/facilities/:id/emission-records", h.ListRecordsByFacility)
			auth.POST("/emission-records/calculate", h.CalculateEmission)

			// Carbon Credits
			auth.GET("/carbon-credits", h.ListCarbonCredits)
			auth.POST("/carbon-credits", h.CreateCarbonCredit)
			auth.GET("/carbon-credits/:id", h.GetCarbonCredit)
			auth.PUT("/carbon-credits/:id", h.UpdateCarbonCredit)
			auth.DELETE("/carbon-credits/:id", h.DeleteCarbonCredit)
			auth.POST("/carbon-credits/:id/retire", h.RetireCarbonCredit)

			// Reduction Targets
			auth.GET("/reduction-targets", h.ListReductionTargets)
			auth.POST("/reduction-targets", h.CreateReductionTarget)
			auth.GET("/reduction-targets/:id", h.GetReductionTarget)
			auth.PUT("/reduction-targets/:id", h.UpdateReductionTarget)
			auth.DELETE("/reduction-targets/:id", h.DeleteReductionTarget)

			// Carbon Reports
			auth.GET("/reports", h.ListReports)
			auth.POST("/reports", h.CreateReport)
			auth.GET("/reports/:id", h.GetReport)
			auth.PUT("/reports/:id", h.UpdateReport)
			auth.DELETE("/reports/:id", h.DeleteReport)
			auth.POST("/reports/:id/generate", h.GenerateReport)

			// Analytics
			auth.GET("/analytics/dashboard", h.Dashboard)
			auth.GET("/analytics/by-scope", h.EmissionsByScope)
			auth.GET("/analytics/trend", h.EmissionsTrend)
			auth.GET("/analytics/comparison", h.EmissionsComparison)
			auth.GET("/analytics/facility-breakdown", h.FacilityBreakdown)

			// Audit Logs (admin only)
			auth.GET("/audit-logs", middleware.RequireRole("admin"), h.ListAuditLogs)
			auth.POST("/audit-logs/:id/rollback", middleware.RequireRole("admin"), h.RollbackAuditLog)
			auth.GET("/rollbacks", middleware.RequireRole("admin"), h.ListRollbacks)

			// ---- 10 ops & management features ----

			// 1. Data import / export
			auth.GET("/data/imports", h.ListDataImports)
			auth.POST("/data/imports", h.CreateDataImport)
			auth.POST("/data/imports/:id/process", h.ProcessDataImport)
			auth.GET("/data/exports/:entity", h.ExportEntity)

			// 2. Scheduled tasks
			auth.GET("/scheduled-tasks", h.ListScheduledTasks)
			auth.POST("/scheduled-tasks", h.CreateScheduledTask)
			auth.GET("/scheduled-tasks/:id", h.GetScheduledTask)
			auth.PUT("/scheduled-tasks/:id", h.UpdateScheduledTask)
			auth.DELETE("/scheduled-tasks/:id", h.DeleteScheduledTask)
			auth.POST("/scheduled-tasks/:id/run", h.RunScheduledTask)
			auth.POST("/scheduled-tasks/:id/pause", h.PauseScheduledTask)

			// 3. Alerts
			auth.GET("/alerts", h.ListAlerts)
			auth.POST("/alerts", h.CreateAlert)
			auth.GET("/alerts/:id", h.GetAlert)
			auth.POST("/alerts/:id/acknowledge", h.AcknowledgeAlert)
			auth.POST("/alerts/:id/resolve", h.ResolveAlert)

			// 4. Notifications
			auth.GET("/notifications", h.ListNotifications)
			auth.POST("/notifications", h.CreateNotification)
			auth.POST("/notifications/:id/send", h.SendNotification)
			auth.GET("/notifications/templates", h.ListNotificationTemplates)

			// 5. API keys
			auth.GET("/api-keys", h.ListAPIKeys)
			auth.POST("/api-keys", h.CreateAPIKey)
			auth.POST("/api-keys/:id/revoke", h.RevokeAPIKey)

			// 6. Webhooks
			auth.GET("/webhooks", h.ListWebhooks)
			auth.POST("/webhooks", h.CreateWebhook)
			auth.GET("/webhooks/:id", h.GetWebhook)
			auth.PUT("/webhooks/:id", h.UpdateWebhook)
			auth.DELETE("/webhooks/:id", h.DeleteWebhook)
			auth.POST("/webhooks/:id/test", h.TestWebhook)

			// 7. Attachments
			auth.GET("/attachments", h.ListAttachments)
			auth.POST("/attachments", h.CreateAttachment)
			auth.GET("/attachments/:id", h.GetAttachment)
			auth.DELETE("/attachments/:id", h.DeleteAttachment)

			// 8. Report exports
			auth.GET("/report-exports", h.ListReportExports)
			auth.POST("/reports/:id/export", h.ExportReport)

			// 9. Rollback (routes above under audit-logs)

			// 10. System settings
			auth.GET("/settings", h.ListSettings)
			auth.POST("/settings", h.CreateSetting)
			auth.GET("/settings/:id", h.GetSetting)
			auth.PUT("/settings/:id", h.UpdateSetting)
			auth.DELETE("/settings/:id", h.DeleteSetting)
		}
	}

	return r
}
