package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type AuditLogRepo struct {
	*GenericRepo[model.AuditLog]
	db *sql.DB
}

func NewAuditLogRepo(db *sql.DB) *AuditLogRepo {
	return &AuditLogRepo{
		GenericRepo: NewGenericRepo[model.AuditLog](db, "audit_logs", func() *model.AuditLog { return &model.AuditLog{} }),
		db:          db,
	}
}

// Record writes a single audit entry.
func (r *AuditLogRepo) Record(l *model.AuditLog) error {
	q := `INSERT INTO audit_logs (user_id, action, entity, entity_id, detail, ip_address)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`
	return r.db.QueryRow(q, l.UserID, l.Action, l.Entity, l.EntityID, l.Detail, l.IPAddress).Scan(&l.ID, &l.CreatedAt)
}
