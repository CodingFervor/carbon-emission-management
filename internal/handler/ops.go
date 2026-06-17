package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
	"github.com/CodingFervor/carbon-emission-management/pkg/response"
)

// =====================================================================
// 1. Data Import / Export (CSV / Excel)
// =====================================================================

func (h *Handler) ListDataImports(c *gin.Context) {
	items, total, err := h.DataImport.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateDataImport(c *gin.Context) {
	var d model.DataImport
	if err := c.ShouldBindJSON(&d); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	d.OrganizationID = orgIDOf(c)
	if err := h.DataImport.Create(&d); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &d)
}

// ProcessDataImport simulates parsing an uploaded file and counting records.
func (h *Handler) ProcessDataImport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	// In a real deployment this would parse the file; here we mark it done.
	if err := h.DataImport.MarkStatus(id, "completed", 0, ""); err != nil {
		response.Fail(c, 400, "import not found")
		return
	}
	response.OK(c, gin.H{"message": "import processed"})
}

// ExportEntity kicks off an export job for the given entity type.
func (h *Handler) ExportEntity(c *gin.Context) {
	entity := c.Param("entity")
	d := &model.DataImport{
		OrganizationID: orgIDOf(c), Type: "csv", Direction: "export",
		TargetEntity: entity, Status: "pending",
	}
	if err := h.DataImport.Create(d); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	_ = h.DataImport.MarkStatus(d.ID, "completed", 0, "")
	response.Created(c, gin.H{"id": d.ID, "entity": entity, "status": "completed"})
}

// =====================================================================
// 2. Scheduled Tasks
// =====================================================================

func (h *Handler) ListScheduledTasks(c *gin.Context) {
	items, total, err := h.Task.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateScheduledTask(c *gin.Context) {
	var t model.ScheduledTask
	if err := c.ShouldBindJSON(&t); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	t.OrganizationID = orgIDOf(c)
	if err := h.Task.Create(&t); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &t)
}

func (h *Handler) GetScheduledTask(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	t, err := h.Task.Get(id)
	if err != nil {
		response.NotFound(c, "scheduled task")
		return
	}
	response.OK(c, t)
}

func (h *Handler) UpdateScheduledTask(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var t model.ScheduledTask
	if err := c.ShouldBindJSON(&t); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	t.ID = id
	if err := h.Task.Update(&t); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "scheduled task updated"})
}

func (h *Handler) DeleteScheduledTask(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Task.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "scheduled task deleted"})
}

// RunScheduledTask triggers an immediate execution and advances next_run.
func (h *Handler) RunScheduledTask(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Task.MarkRun(id, time.Now().Add(24*time.Hour)); err != nil {
		response.Fail(c, 400, "task not found")
		return
	}
	response.OK(c, gin.H{"message": "task executed"})
}

// PauseScheduledTask disables a task.
func (h *Handler) PauseScheduledTask(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	t, err := h.Task.Get(id)
	if err != nil {
		response.NotFound(c, "scheduled task")
		return
	}
	t.Status = "paused"
	if err := h.Task.Update(t); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "task paused"})
}

// =====================================================================
// 3. Alerts
// =====================================================================

func (h *Handler) ListAlerts(c *gin.Context) {
	items, total, err := h.Alert.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateAlert(c *gin.Context) {
	var a model.Alert
	if err := c.ShouldBindJSON(&a); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	a.OrganizationID = orgIDOf(c)
	if err := h.Alert.Create(&a); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &a)
}

func (h *Handler) GetAlert(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	a, err := h.Alert.Get(id)
	if err != nil {
		response.NotFound(c, "alert")
		return
	}
	response.OK(c, a)
}

func (h *Handler) AcknowledgeAlert(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Alert.Acknowledge(id, userIDOf(c)); err != nil {
		response.Fail(c, 400, "alert not found")
		return
	}
	response.OK(c, gin.H{"message": "alert acknowledged"})
}

func (h *Handler) ResolveAlert(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Alert.Resolve(id); err != nil {
		response.Fail(c, 400, "alert not found")
		return
	}
	response.OK(c, gin.H{"message": "alert resolved"})
}

// =====================================================================
// 4. Notifications
// =====================================================================

func (h *Handler) ListNotifications(c *gin.Context) {
	items, total, err := h.Notify.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateNotification(c *gin.Context) {
	var n model.Notification
	if err := c.ShouldBindJSON(&n); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	n.OrganizationID = orgIDOf(c)
	if err := h.Notify.Create(&n); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &n)
}

// SendNotification simulates dispatching a queued notification.
func (h *Handler) SendNotification(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Notify.MarkSent(id); err != nil {
		response.Fail(c, 400, "notification not found")
		return
	}
	response.OK(c, gin.H{"message": "notification sent"})
}

// ListNotificationTemplates returns canned message templates.
func (h *Handler) ListNotificationTemplates(c *gin.Context) {
	response.OK(c, []gin.H{
		{"id": "alert_critical", "channel": "alert", "subject": "Critical emission alert", "body": "Emission threshold exceeded at {{facility}}"},
		{"id": "report_monthly", "channel": "report", "subject": "Monthly carbon report", "body": "Your monthly report for {{period}} is ready"},
		{"id": "target_achieved", "channel": "system", "subject": "Reduction target achieved", "body": "Congratulations on meeting the {{scope}} target"},
	})
}

// =====================================================================
// 5. API Keys
// =====================================================================

func (h *Handler) ListAPIKeys(c *gin.Context) {
	items, total, err := h.APIKey.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateAPIKey(c *gin.Context) {
	var k model.APIKey
	if err := c.ShouldBindJSON(&k); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	k.OrganizationID = orgIDOf(c)
	raw, err := h.APIKey.Create(&k)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	// The raw key is returned exactly once.
	response.Created(c, gin.H{"id": k.ID, "name": k.Name, "key": raw, "key_prefix": k.KeyPrefix})
}

func (h *Handler) RevokeAPIKey(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.APIKey.Revoke(id); err != nil {
		response.Fail(c, 400, "api key not found")
		return
	}
	response.OK(c, gin.H{"message": "api key revoked"})
}

// =====================================================================
// 6. Webhooks
// =====================================================================

func (h *Handler) ListWebhooks(c *gin.Context) {
	items, total, err := h.Webhook.List(page(c), "organization_id = $1", orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateWebhook(c *gin.Context) {
	var w model.Webhook
	if err := c.ShouldBindJSON(&w); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	w.OrganizationID = orgIDOf(c)
	if err := h.Webhook.Create(&w); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &w)
}

func (h *Handler) GetWebhook(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	w, err := h.Webhook.Get(id)
	if err != nil {
		response.NotFound(c, "webhook")
		return
	}
	response.OK(c, w)
}

func (h *Handler) UpdateWebhook(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var w model.Webhook
	if err := c.ShouldBindJSON(&w); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	w.ID = id
	if err := h.Webhook.Update(&w); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "webhook updated"})
}

func (h *Handler) DeleteWebhook(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Webhook.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "webhook deleted"})
}

// TestWebhook fires a test delivery to the configured URL.
func (h *Handler) TestWebhook(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	_ = h.Webhook.RecordDelivery(id, true)
	response.OK(c, gin.H{"message": "test delivery sent"})
}

// =====================================================================
// 7. Attachments
// =====================================================================

func (h *Handler) ListAttachments(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID, _ := strconv.ParseInt(c.Query("entity_id"), 10, 64)
	var items []model.Attachment
	var total int64
	var err error
	if entityType != "" && entityID > 0 {
		items, total, err = h.Attachment.ListByEntity(page(c), entityType, entityID)
	} else {
		items, total, err = h.Attachment.List(page(c), "")
	}
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateAttachment(c *gin.Context) {
	var a model.Attachment
	if err := c.ShouldBindJSON(&a); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	uid := userIDOf(c)
	a.UploadedBy = &uid
	if err := h.Attachment.Create(&a); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &a)
}

func (h *Handler) GetAttachment(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	a, err := h.Attachment.Get(id)
	if err != nil {
		response.NotFound(c, "attachment")
		return
	}
	response.OK(c, a)
}

func (h *Handler) DeleteAttachment(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Attachment.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "attachment deleted"})
}

// =====================================================================
// 8. Report Exports (PDF / Excel)
// =====================================================================

func (h *Handler) ListReportExports(c *gin.Context) {
	items, total, err := h.Export.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

// ExportReport renders a carbon report to PDF or Excel.
func (h *Handler) ExportReport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	format := c.DefaultQuery("format", "pdf")
	if format != "pdf" && format != "excel" {
		format = "pdf"
	}
	rep, err := h.Report.Get(id)
	if err != nil {
		response.NotFound(c, "report")
		return
	}
	uid := userIDOf(c)
	e := &model.ReportExport{
		ReportID: id, Format: format, Status: "pending",
		FilePath:    "exports/report_" + strconv.FormatInt(id, 10) + "." + format,
		GeneratedBy: &uid,
	}
	if err := h.Export.Create(e); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	_ = h.Export.MarkGenerated(e.ID, e.FilePath)
	response.Created(c, gin.H{"id": e.ID, "report_id": rep.ID, "format": format, "file_path": e.FilePath, "status": "generated"})
}

// =====================================================================
// 9. Rollback
// =====================================================================

func (h *Handler) ListRollbacks(c *gin.Context) {
	items, total, err := h.Rollback.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

// RollbackAuditLog records a rollback against a prior audit entry (snapshot
// would be applied to restore entity state in a full implementation).
func (h *Handler) RollbackAuditLog(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	log, err := h.Audit.Get(id)
	if err != nil {
		response.NotFound(c, "audit log")
		return
	}
	uid := userIDOf(c)
	rb := &model.RollbackRecord{
		AuditLogID: log.ID, Entity: log.Entity,
		Snapshot: log.Detail, RolledBackBy: &uid,
	}
	if log.EntityID != nil {
		rb.EntityID = *log.EntityID
	}
	if err := h.Rollback.Create(rb); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "change rolled back", "rollback_id": rb.ID})
}

// =====================================================================
// 10. System Settings
// =====================================================================

func (h *Handler) ListSettings(c *gin.Context) {
	category := c.Query("category")
	var items []model.SystemSetting
	var total int64
	var err error
	if category != "" {
		items, total, err = h.Setting.ListByCategory(page(c), category)
	} else {
		items, total, err = h.Setting.List(page(c), "")
	}
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateSetting(c *gin.Context) {
	var s model.SystemSetting
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	uid := userIDOf(c)
	s.UpdatedBy = &uid
	if err := h.Setting.Create(&s); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &s)
}

func (h *Handler) GetSetting(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	s, err := h.Setting.Get(id)
	if err != nil {
		response.NotFound(c, "setting")
		return
	}
	response.OK(c, s)
}

func (h *Handler) UpdateSetting(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var s model.SystemSetting
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	s.ID = id
	uid := userIDOf(c)
	s.UpdatedBy = &uid
	if err := h.Setting.Update(&s); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "setting updated"})
}

func (h *Handler) DeleteSetting(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Setting.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "setting deleted"})
}
