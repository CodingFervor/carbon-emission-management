package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type ReductionTargetRepo struct {
	*GenericRepo[model.ReductionTarget]
	db *sql.DB
}

func NewReductionTargetRepo(db *sql.DB) *ReductionTargetRepo {
	return &ReductionTargetRepo{
		GenericRepo: NewGenericRepo[model.ReductionTarget](db, "reduction_targets", func() *model.ReductionTarget { return &model.ReductionTarget{} }),
		db:          db,
	}
}

func (r *ReductionTargetRepo) Create(t *model.ReductionTarget) error {
	q := `INSERT INTO reduction_targets (organization_id, scope, baseline_year, target_year, target_pct, baseline_co2_t, current_co2_t, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q,
		t.OrganizationID, t.Scope, t.BaselineYear, t.TargetYear, t.TargetPct,
		t.BaselineCO2T, t.CurrentCO2T, defaultStr(t.Status, "on_track"),
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *ReductionTargetRepo) Update(t *model.ReductionTarget) error {
	q := `UPDATE reduction_targets SET scope=$1, baseline_year=$2, target_year=$3, target_pct=$4,
	      baseline_co2_t=$5, current_co2_t=$6, status=$7, updated_at=CURRENT_TIMESTAMP WHERE id=$8`
	_, err := r.db.Exec(q, t.Scope, t.BaselineYear, t.TargetYear, t.TargetPct, t.BaselineCO2T, t.CurrentCO2T, defaultStr(t.Status, "on_track"), t.ID)
	return err
}
