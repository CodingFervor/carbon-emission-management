package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type EmissionFactorRepo struct {
	*GenericRepo[model.EmissionFactor]
}

func NewEmissionFactorRepo(db *sql.DB) *EmissionFactorRepo {
	return &EmissionFactorRepo{GenericRepo: NewGenericRepo[model.EmissionFactor](db, "emission_factors", func() *model.EmissionFactor { return &model.EmissionFactor{} })}
}

func (r *EmissionFactorRepo) Create(f *model.EmissionFactor) error {
	q := `INSERT INTO emission_factors (name, activity_unit, factor_value, co2_unit, scope, source_ref, valid_from, valid_to)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`
	return r.DB().QueryRow(q,
		f.Name, f.ActivityUnit, f.FactorValue, f.CO2Unit, f.Scope, f.SourceRef, f.ValidFrom, f.ValidTo,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func (r *EmissionFactorRepo) Update(f *model.EmissionFactor) error {
	q := `UPDATE emission_factors SET name=$1, activity_unit=$2, factor_value=$3, co2_unit=$4, scope=$5,
	      source_ref=$6, valid_from=$7, valid_to=$8, updated_at=CURRENT_TIMESTAMP WHERE id=$9`
	_, err := r.DB().Exec(q, f.Name, f.ActivityUnit, f.FactorValue, f.CO2Unit, f.Scope, f.SourceRef, f.ValidFrom, f.ValidTo, f.ID)
	return err
}
