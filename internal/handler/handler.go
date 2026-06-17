package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/repository"
	"github.com/CodingFervor/carbon-emission-management/pkg/response"
)

// Handler bundles all repositories so each resource handler can embed it and
// share access to the data layer.
type Handler struct {
	Org       *repository.OrganizationRepo
	Facility  *repository.FacilityRepo
	Source    *repository.EmissionSourceRepo
	Factor    *repository.EmissionFactorRepo
	Record    *repository.EmissionRecordRepo
	Credit    *repository.CarbonCreditRepo
	Target    *repository.ReductionTargetRepo
	Report    *repository.CarbonReportRepo
	Audit     *repository.AuditLogRepo
	Analytics *repository.AnalyticsRepo
	User      *repository.UserRepo

	// Ops & management features
	DataImport *repository.DataImportRepo
	Task       *repository.ScheduledTaskRepo
	Alert      *repository.AlertRepo
	Notify     *repository.NotificationRepo
	APIKey     *repository.APIKeyRepo
	Webhook    *repository.WebhookRepo
	Attachment *repository.AttachmentRepo
	Export     *repository.ReportExportRepo
	Rollback   *repository.RollbackRepo
	Setting    *repository.SystemSettingRepo
}

// New builds a Handler backed by the given repos.
func New(h *Handler) *Handler { return h }

// page parses paging query params (?page=&page_size=).
func page(c *gin.Context) repository.Page {
	p, _ := strconv.Atoi(c.Query("page"))
	ps, _ := strconv.Atoi(c.Query("page_size"))
	return repository.PageFromValues(p, ps)
}

// idParam parses the :id path parameter.
func idParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "invalid id")
		return 0, false
	}
	return id, true
}

// orgIDOf returns the current user's organization id from the auth context.
func orgIDOf(c *gin.Context) int64 {
	if v, ok := c.Get("organizationID"); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// userIDOf returns the current user's id from the auth context.
func userIDOf(c *gin.Context) int64 {
	if v, ok := c.Get("userID"); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
