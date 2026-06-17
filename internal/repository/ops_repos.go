package repository

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

// ----- 1. DataImport (CSV/Excel import & export) -----

type DataImportRepo struct {
	*GenericRepo[model.DataImport]
	db *sql.DB
}

func NewDataImportRepo(db *sql.DB) *DataImportRepo {
	return &DataImportRepo{
		GenericRepo: NewGenericRepo[model.DataImport](db, "data_imports", func() *model.DataImport { return &model.DataImport{} }),
		db:          db,
	}
}

func (r *DataImportRepo) Create(d *model.DataImport) error {
	q := `INSERT INTO data_imports (organization_id, type, direction, target_entity, file_path, status, started_by)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, d.OrganizationID, d.Type, defaultStr(d.Direction, "import"),
		d.TargetEntity, d.FilePath, defaultStr(d.Status, "pending"), d.StartedBy,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

// MarkStatus updates a job's status, record count, and error message.
func (r *DataImportRepo) MarkStatus(id int64, status string, count int, errMsg string) error {
	_, err := r.db.Exec(
		"UPDATE data_imports SET status=$1, records_count=$2, error_msg=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4",
		status, count, errMsg, id,
	)
	return err
}

// ----- 2. ScheduledTask -----

type ScheduledTaskRepo struct {
	*GenericRepo[model.ScheduledTask]
	db *sql.DB
}

func NewScheduledTaskRepo(db *sql.DB) *ScheduledTaskRepo {
	return &ScheduledTaskRepo{
		GenericRepo: NewGenericRepo[model.ScheduledTask](db, "scheduled_tasks", func() *model.ScheduledTask { return &model.ScheduledTask{} }),
		db:          db,
	}
}

func (r *ScheduledTaskRepo) Create(t *model.ScheduledTask) error {
	q := `INSERT INTO scheduled_tasks (organization_id, name, cron_expr, target_endpoint, payload, status)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, t.OrganizationID, t.Name, t.CronExpr, t.TargetEndpoint, t.Payload,
		defaultStr(t.Status, "active")).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *ScheduledTaskRepo) Update(t *model.ScheduledTask) error {
	q := `UPDATE scheduled_tasks SET name=$1, cron_expr=$2, target_endpoint=$3, payload=$4,
	      status=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$6`
	_, err := r.db.Exec(q, t.Name, t.CronExpr, t.TargetEndpoint, t.Payload, defaultStr(t.Status, "active"), t.ID)
	return err
}

// MarkRun records that the task just ran and sets its next run time.
func (r *ScheduledTaskRepo) MarkRun(id int64, next time.Time) error {
	_, err := r.db.Exec(
		"UPDATE scheduled_tasks SET last_run=CURRENT_TIMESTAMP, next_run=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2",
		next, id)
	return err
}

// ----- 3. Alert -----

type AlertRepo struct {
	*GenericRepo[model.Alert]
	db *sql.DB
}

func NewAlertRepo(db *sql.DB) *AlertRepo {
	return &AlertRepo{
		GenericRepo: NewGenericRepo[model.Alert](db, "alerts", func() *model.Alert { return &model.Alert{} }),
		db:          db,
	}
}

func (r *AlertRepo) Create(a *model.Alert) error {
	q := `INSERT INTO alerts (organization_id, facility_id, type, severity, message, trigger_value, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, a.OrganizationID, a.FacilityID, a.Type, a.Severity, a.Message,
		a.TriggerValue, defaultStr(a.Status, "active")).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AlertRepo) Update(a *model.Alert) error {
	q := `UPDATE alerts SET type=$1, severity=$2, message=$3, trigger_value=$4, status=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$6`
	_, err := r.db.Exec(q, a.Type, a.Severity, a.Message, a.TriggerValue, defaultStr(a.Status, "active"), a.ID)
	return err
}

func (r *AlertRepo) Acknowledge(id int64, userID int64) error {
	_, err := r.db.Exec(
		"UPDATE alerts SET status='acknowledged', acked_by=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2",
		userID, id)
	return err
}

func (r *AlertRepo) Resolve(id int64) error {
	_, err := r.db.Exec("UPDATE alerts SET status='resolved', updated_at=CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

// ----- 4. Notification -----

type NotificationRepo struct {
	*GenericRepo[model.Notification]
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{
		GenericRepo: NewGenericRepo[model.Notification](db, "notifications", func() *model.Notification { return &model.Notification{} }),
		db:          db,
	}
}

func (r *NotificationRepo) Create(n *model.Notification) error {
	q := `INSERT INTO notifications (organization_id, type, recipient, subject, body, channel, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, n.OrganizationID, n.Type, n.Recipient, n.Subject, n.Body,
		n.Channel, defaultStr(n.Status, "pending")).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

// MarkSent records a successful delivery.
func (r *NotificationRepo) MarkSent(id int64) error {
	_, err := r.db.Exec(
		"UPDATE notifications SET status='sent', sent_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

// ----- 5. APIKey -----

type APIKeyRepo struct {
	*GenericRepo[model.APIKey]
	db *sql.DB
}

func NewAPIKeyRepo(db *sql.DB) *APIKeyRepo {
	return &APIKeyRepo{
		GenericRepo: NewGenericRepo[model.APIKey](db, "api_keys", func() *model.APIKey { return &model.APIKey{} }),
		db:          db,
	}
}

// Create generates a new raw key, stores only its hash, and returns the raw key once.
func (r *APIKeyRepo) Create(k *model.APIKey) (string, error) {
	raw, err := generateAPIKey()
	if err != nil {
		return "", err
	}
	k.KeyHash = hashKey(raw)
	k.KeyPrefix = raw[:10]
	k.Status = defaultStr(k.Status, "active")
	q := `INSERT INTO api_keys (organization_id, name, key_hash, key_prefix, scopes, expires_at, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`
	if err := r.db.QueryRow(q, k.OrganizationID, k.Name, k.KeyHash, k.KeyPrefix, k.Scopes,
		k.ExpiresAt, k.Status).Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt); err != nil {
		return "", err
	}
	return raw, nil
}

// FindByRawKey looks up an API key by its raw value (hashing first).
func (r *APIKeyRepo) FindByRawKey(raw string) (*model.APIKey, error) {
	q := "SELECT id, organization_id, name, key_hash, key_prefix, scopes, expires_at, last_used, status, created_at, updated_at FROM api_keys WHERE key_hash=$1"
	k := &model.APIKey{}
	if err := r.db.QueryRow(q, hashKey(raw)).Scan(
		&k.ID, &k.OrganizationID, &k.Name, &k.KeyHash, &k.KeyPrefix, &k.Scopes,
		&k.ExpiresAt, &k.LastUsed, &k.Status, &k.CreatedAt, &k.UpdatedAt); err != nil {
		return nil, err
	}
	return k, nil
}

// TouchLastUsed records the most recent use of a key.
func (r *APIKeyRepo) TouchLastUsed(id int64) error {
	_, err := r.db.Exec("UPDATE api_keys SET last_used=CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

// Revoke disables a key permanently.
func (r *APIKeyRepo) Revoke(id int64) error {
	_, err := r.db.Exec("UPDATE api_keys SET status='revoked', updated_at=CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "cem_" + hex.EncodeToString(b), nil
}

func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// ----- 6. Webhook -----

type WebhookRepo struct {
	*GenericRepo[model.Webhook]
	db *sql.DB
}

func NewWebhookRepo(db *sql.DB) *WebhookRepo {
	return &WebhookRepo{
		GenericRepo: NewGenericRepo[model.Webhook](db, "webhooks", func() *model.Webhook { return &model.Webhook{} }),
		db:          db,
	}
}

func (r *WebhookRepo) Create(w *model.Webhook) error {
	q := `INSERT INTO webhooks (organization_id, url, events, secret, status)
	      VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, w.OrganizationID, w.URL, w.Events, w.Secret,
		defaultStr(w.Status, "active")).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *WebhookRepo) Update(w *model.Webhook) error {
	q := `UPDATE webhooks SET url=$1, events=$2, secret=$3, status=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`
	_, err := r.db.Exec(q, w.URL, w.Events, w.Secret, defaultStr(w.Status, "active"), w.ID)
	return err
}

// RecordDelivery logs the outcome of a webhook delivery attempt.
func (r *WebhookRepo) RecordDelivery(id int64, ok bool) error {
	if ok {
		_, err := r.db.Exec(
			"UPDATE webhooks SET failure_count=0, last_fired=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1", id)
		return err
	}
	_, err := r.db.Exec(
		"UPDATE webhooks SET failure_count=failure_count+1, last_fired=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

// ----- 7. Attachment -----

type AttachmentRepo struct {
	*GenericRepo[model.Attachment]
	db *sql.DB
}

func NewAttachmentRepo(db *sql.DB) *AttachmentRepo {
	return &AttachmentRepo{
		GenericRepo: NewGenericRepo[model.Attachment](db, "attachments", func() *model.Attachment { return &model.Attachment{} }),
		db:          db,
	}
}

func (r *AttachmentRepo) Create(a *model.Attachment) error {
	q := `INSERT INTO attachments (entity_type, entity_id, filename, file_path, file_size, mime_type, uploaded_by)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at`
	return r.db.QueryRow(q, a.EntityType, a.EntityID, a.Filename, a.FilePath, a.FileSize,
		a.MimeType, a.UploadedBy).Scan(&a.ID, &a.CreatedAt)
}

// ListByEntity returns attachments for a given entity.
func (r *AttachmentRepo) ListByEntity(p Page, entityType string, entityID int64) ([]model.Attachment, int64, error) {
	return r.List(p, "entity_type = $1 AND entity_id = $2", entityType, entityID)
}

// ----- 8. ReportExport -----

type ReportExportRepo struct {
	*GenericRepo[model.ReportExport]
	db *sql.DB
}

func NewReportExportRepo(db *sql.DB) *ReportExportRepo {
	return &ReportExportRepo{
		GenericRepo: NewGenericRepo[model.ReportExport](db, "report_exports", func() *model.ReportExport { return &model.ReportExport{} }),
		db:          db,
	}
}

func (r *ReportExportRepo) Create(e *model.ReportExport) error {
	q := `INSERT INTO report_exports (report_id, format, file_path, options, status, generated_by)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`
	return r.db.QueryRow(q, e.ReportID, e.Format, e.FilePath, e.Options,
		defaultStr(e.Status, "pending"), e.GeneratedBy).Scan(&e.ID, &e.CreatedAt)
}

func (r *ReportExportRepo) MarkGenerated(id int64, filePath string) error {
	_, err := r.db.Exec(
		"UPDATE report_exports SET status='generated', file_path=$1 WHERE id=$2", filePath, id)
	return err
}

// ----- 9. RollbackRecord -----

type RollbackRepo struct {
	*GenericRepo[model.RollbackRecord]
	db *sql.DB
}

func NewRollbackRepo(db *sql.DB) *RollbackRepo {
	return &RollbackRepo{
		GenericRepo: NewGenericRepo[model.RollbackRecord](db, "rollback_records", func() *model.RollbackRecord { return &model.RollbackRecord{} }),
		db:          db,
	}
}

func (r *RollbackRepo) Create(rb *model.RollbackRecord) error {
	q := `INSERT INTO rollback_records (audit_log_id, entity, entity_id, snapshot, rolled_back_by)
	      VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`
	return r.db.QueryRow(q, rb.AuditLogID, rb.Entity, rb.EntityID, rb.Snapshot, rb.RolledBackBy).Scan(&rb.ID, &rb.CreatedAt)
}

// ----- 10. SystemSetting -----

type SystemSettingRepo struct {
	*GenericRepo[model.SystemSetting]
	db *sql.DB
}

func NewSystemSettingRepo(db *sql.DB) *SystemSettingRepo {
	return &SystemSettingRepo{
		GenericRepo: NewGenericRepo[model.SystemSetting](db, "system_settings", func() *model.SystemSetting { return &model.SystemSetting{} }),
		db:          db,
	}
}

func (r *SystemSettingRepo) Create(s *model.SystemSetting) error {
	q := `INSERT INTO system_settings (key, value, category, description, data_type, updated_by)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, s.Key, s.Value, s.Category, s.Description, defaultStr(s.DataType, "string"), s.UpdatedBy).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SystemSettingRepo) Update(s *model.SystemSetting) error {
	q := `UPDATE system_settings SET value=$1, category=$2, description=$3, data_type=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`
	_, err := r.db.Exec(q, s.Value, s.Category, s.Description, defaultStr(s.DataType, "string"), s.ID)
	return err
}

// ListByCategory returns settings in a category.
func (r *SystemSettingRepo) ListByCategory(p Page, category string) ([]model.SystemSetting, int64, error) {
	return r.List(p, "category = $1", category)
}
