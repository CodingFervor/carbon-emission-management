package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type OrganizationRepo struct {
	*GenericRepo[model.Organization]
}

func NewOrganizationRepo(db *sql.DB) *OrganizationRepo {
	return &OrganizationRepo{GenericRepo: NewGenericRepo[model.Organization](db, "organizations", func() *model.Organization { return &model.Organization{} })}
}

func (r *OrganizationRepo) Create(o *model.Organization) error {
	q := `INSERT INTO organizations (name, industry, country, reporting_year, base_year, status)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	return r.DB().QueryRow(q,
		o.Name, o.Industry, o.Country, o.ReportingYear, o.BaseYear, defaultStr(o.Status, "active"),
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrganizationRepo) Update(o *model.Organization) error {
	q := `UPDATE organizations SET name=$1, industry=$2, country=$3, reporting_year=$4,
	      base_year=$5, status=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`
	_, err := r.DB().Exec(q, o.Name, o.Industry, o.Country, o.ReportingYear, o.BaseYear, defaultStr(o.Status, "active"), o.ID)
	return err
}
