package repository

import (
	"database/sql"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type CarbonReportRepo struct {
	*GenericRepo[model.CarbonReport]
	db *sql.DB
}

func NewCarbonReportRepo(db *sql.DB) *CarbonReportRepo {
	return &CarbonReportRepo{
		GenericRepo: NewGenericRepo[model.CarbonReport](db, "carbon_reports", func() *model.CarbonReport { return &model.CarbonReport{} }),
		db:          db,
	}
}

func (r *CarbonReportRepo) Create(rep *model.CarbonReport) error {
	q := `INSERT INTO carbon_reports (organization_id, period, period_start, period_end, total_co2_t,
	      scope1_co2_t, scope2_co2_t, scope3_co2_t, offsets_t, net_co2_t, status, standard)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	      RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q,
		rep.OrganizationID, rep.Period, rep.PeriodStart, rep.PeriodEnd,
		rep.TotalCO2T, rep.Scope1CO2T, rep.Scope2CO2T, rep.Scope3CO2T, rep.OffsetsT, rep.NetCO2T,
		defaultStr(rep.Status, "draft"), defaultStr(rep.Standard, "GHGP"),
	).Scan(&rep.ID, &rep.CreatedAt, &rep.UpdatedAt)
}

// Generate compiles a report by aggregating emission records in the period
// joined to their (scope-bearing) emission sources, then persists the totals.
func (r *CarbonReportRepo) Generate(rep *model.CarbonReport) error {
	q := `SELECT
		  COALESCE(SUM(CASE WHEN es.scope=1 THEN er.co2_kg END),0)/1000.0,
		  COALESCE(SUM(CASE WHEN es.scope=2 THEN er.co2_kg END),0)/1000.0,
		  COALESCE(SUM(CASE WHEN es.scope=3 THEN er.co2_kg END),0)/1000.0
		FROM emission_records er
		JOIN emission_sources es ON es.id = er.source_id
		JOIN facilities f ON f.id = er.facility_id
		WHERE f.organization_id = $1
		  AND er.period_start >= $2 AND er.period_end <= $3
		  AND er.status = 'verified'`
	var s1, s2, s3 float64
	err := r.db.QueryRow(q, rep.OrganizationID, rep.PeriodStart, rep.PeriodEnd).Scan(&s1, &s2, &s3)
	if err != nil {
		return err
	}
	rep.Scope1CO2T = s1
	rep.Scope2CO2T = s2
	rep.Scope3CO2T = s3
	rep.TotalCO2T = s1 + s2 + s3
	rep.NetCO2T = rep.TotalCO2T - rep.OffsetsT
	rep.Status = "generated"
	return r.Create(rep)
}
