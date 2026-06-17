package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type FacilityRepo struct {
	*GenericRepo[model.Facility]
}

func NewFacilityRepo(db *sql.DB) *FacilityRepo {
	return &FacilityRepo{GenericRepo: NewGenericRepo[model.Facility](db, "facilities", func() *model.Facility { return &model.Facility{} })}
}

// ListByOrganization returns facilities belonging to an organization.
func (r *FacilityRepo) ListByOrganization(p Page, orgID int64) ([]model.Facility, int64, error) {
	return r.List(p, "organization_id = $1", orgID)
}

func (r *FacilityRepo) Create(f *model.Facility) error {
	q := `INSERT INTO facilities (organization_id, name, address, latitude, longitude, type, country, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`
	return r.DB().QueryRow(q,
		f.OrganizationID, f.Name, f.Address, f.Latitude, f.Longitude, f.Type, f.Country, defaultStr(f.Status, "active"),
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func (r *FacilityRepo) Update(f *model.Facility) error {
	q := `UPDATE facilities SET name=$1, address=$2, latitude=$3, longitude=$4, type=$5,
	      country=$6, status=$7, updated_at=CURRENT_TIMESTAMP WHERE id=$8`
	_, err := r.DB().Exec(q, f.Name, f.Address, f.Latitude, f.Longitude, f.Type, f.Country, defaultStr(f.Status, "active"), f.ID)
	return err
}
