package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type EmissionSourceRepo struct {
	*GenericRepo[model.EmissionSource]
}

func NewEmissionSourceRepo(db *sql.DB) *EmissionSourceRepo {
	return &EmissionSourceRepo{GenericRepo: NewGenericRepo[model.EmissionSource](db, "emission_sources", func() *model.EmissionSource { return &model.EmissionSource{} })}
}

func (r *EmissionSourceRepo) Create(s *model.EmissionSource) error {
	q := `INSERT INTO emission_sources (facility_id, name, scope, category, fuel_type, unit, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`
	return r.DB().QueryRow(q,
		s.FacilityID, s.Name, s.Scope, s.Category, s.FuelType, s.Unit, defaultStr(s.Status, "active"),
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *EmissionSourceRepo) Update(s *model.EmissionSource) error {
	q := `UPDATE emission_sources SET name=$1, scope=$2, category=$3, fuel_type=$4, unit=$5,
	      status=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`
	_, err := r.DB().Exec(q, s.Name, s.Scope, s.Category, s.FuelType, s.Unit, defaultStr(s.Status, "active"), s.ID)
	return err
}

// ListByFacility returns sources attached to a facility.
func (r *EmissionSourceRepo) ListByFacility(p Page, facilityID int64) ([]model.EmissionSource, int64, error) {
	return r.List(p, "facility_id = $1", facilityID)
}
