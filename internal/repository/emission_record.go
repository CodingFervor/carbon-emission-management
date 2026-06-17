package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type EmissionRecordRepo struct {
	*GenericRepo[model.EmissionRecord]
	db *sql.DB
}

func NewEmissionRecordRepo(db *sql.DB) *EmissionRecordRepo {
	return &EmissionRecordRepo{
		GenericRepo: NewGenericRepo[model.EmissionRecord](db, "emission_records", func() *model.EmissionRecord { return &model.EmissionRecord{} }),
		db:          db,
	}
}

// ListByFacility returns records for a facility.
func (r *EmissionRecordRepo) ListByFacility(p Page, facilityID int64) ([]model.EmissionRecord, int64, error) {
	return r.List(p, "facility_id = $1", facilityID)
}

func (r *EmissionRecordRepo) Create(rec *model.EmissionRecord) error {
	q := `INSERT INTO emission_records
	      (source_id, facility_id, period, period_start, period_end, activity_value, factor_id, co2_kg, notes, status)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	      RETURNING id, calculated_at, created_at, updated_at`
	return r.db.QueryRow(q,
		rec.SourceID, rec.FacilityID, defaultStr(rec.Period, "monthly"),
		rec.PeriodStart, rec.PeriodEnd, rec.ActivityValue, rec.FactorID, rec.CO2Kg, rec.Notes,
		defaultStr(rec.Status, "draft"),
	).Scan(&rec.ID, &rec.CalculatedAt, &rec.CreatedAt, &rec.UpdatedAt)
}

func (r *EmissionRecordRepo) Update(rec *model.EmissionRecord) error {
	q := `UPDATE emission_records SET activity_value=$1, co2_kg=$2, notes=$3, status=$4,
	      updated_at=CURRENT_TIMESTAMP WHERE id=$5`
	_, err := r.db.Exec(q, rec.ActivityValue, rec.CO2Kg, rec.Notes, defaultStr(rec.Status, "draft"), rec.ID)
	return err
}

// CalcRequest carries the inputs for an on-the-fly CO2e calculation.
type CalcRequest struct {
	ActivityValue float64 `json:"activity_value"`
	FactorID      int64   `json:"factor_id"`
}

// Calculate computes CO2e (kg) = activity_value * factor_value, adjusting for CO2 unit.
func (r *EmissionRecordRepo) Calculate(req CalcRequest) (float64, error) {
	var factorValue float64
	var co2Unit string
	err := r.db.QueryRow(
		"SELECT factor_value, co2_unit FROM emission_factors WHERE id = $1", req.FactorID,
	).Scan(&factorValue, &co2Unit)
	if err != nil {
		return 0, err
	}
	co2kg := req.ActivityValue * factorValue
	if co2Unit == "t" {
		co2kg *= 1000 // tons -> kg
	}
	return co2kg, nil
}
