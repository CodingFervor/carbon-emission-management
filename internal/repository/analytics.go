package repository

import "database/sql"

// AnalyticsRepo provides aggregate read queries backing the analytics endpoints.
type AnalyticsRepo struct {
	db *sql.DB
}

func NewAnalyticsRepo(db *sql.DB) *AnalyticsRepo { return &AnalyticsRepo{db: db} }

// Dashboard summarizes total/scope/offset/net emissions for an organization.
type Dashboard struct {
	TotalCO2T    float64 `json:"total_co2_t"`
	Scope1CO2T   float64 `json:"scope1_co2_t"`
	Scope2CO2T   float64 `json:"scope2_co2_t"`
	Scope3CO2T   float64 `json:"scope3_co2_t"`
	OffsetsT     float64 `json:"offsets_t"`
	NetCO2T      float64 `json:"net_co2_t"`
	YoYChangePct float64 `json:"yoy_change_pct"`
}

func (r *AnalyticsRepo) Dashboard(orgID int64) (*Dashboard, error) {
	d := &Dashboard{}
	q := `SELECT
		  COALESCE(SUM(CASE WHEN es.scope=1 THEN er.co2_kg END),0)/1000.0,
		  COALESCE(SUM(CASE WHEN es.scope=2 THEN er.co2_kg END),0)/1000.0,
		  COALESCE(SUM(CASE WHEN es.scope=3 THEN er.co2_kg END),0)/1000.0
		FROM emission_records er
		JOIN emission_sources es ON es.id = er.source_id
		JOIN facilities f ON f.id = er.facility_id
		WHERE f.organization_id = $1`
	if err := r.db.QueryRow(q, orgID).Scan(&d.Scope1CO2T, &d.Scope2CO2T, &d.Scope3CO2T); err != nil {
		return nil, err
	}
	d.TotalCO2T = d.Scope1CO2T + d.Scope2CO2T + d.Scope3CO2T
	// Offsets = retired credits in tons.
	_ = r.db.QueryRow("SELECT COALESCE(SUM(amount_tons),0) FROM carbon_credits WHERE status='retired'").Scan(&d.OffsetsT)
	d.NetCO2T = d.TotalCO2T - d.OffsetsT
	return d, nil
}

// ScopeBucket holds one scope's total emissions (tons).
type ScopeBucket struct {
	Scope int     `json:"scope"`
	CO2T  float64 `json:"co2_t"`
}

func (r *AnalyticsRepo) ByScope(orgID int64) ([]ScopeBucket, error) {
	q := `SELECT es.scope, COALESCE(SUM(er.co2_kg),0)/1000.0
		FROM emission_records er
		JOIN emission_sources es ON es.id = er.source_id
		JOIN facilities f ON f.id = er.facility_id
		WHERE f.organization_id = $1
		GROUP BY es.scope ORDER BY es.scope`
	rows, err := r.db.Query(q, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ScopeBucket{}
	for rows.Next() {
		var b ScopeBucket
		if err := rows.Scan(&b.Scope, &b.CO2T); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}

// TrendPoint is one period's total emissions (tons).
type TrendPoint struct {
	Period string  `json:"period"`
	CO2T   float64 `json:"co2_t"`
}

// Trend aggregates emissions grouped by month for an organization.
func (r *AnalyticsRepo) Trend(orgID int64) ([]TrendPoint, error) {
	q := `SELECT to_char(date_trunc('month', er.period_start),'YYYY-MM') AS period,
		  COALESCE(SUM(er.co2_kg),0)/1000.0
		FROM emission_records er
		JOIN facilities f ON f.id = er.facility_id
		WHERE f.organization_id = $1
		GROUP BY period ORDER BY period`
	rows, err := r.db.Query(q, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []TrendPoint{}
	for rows.Next() {
		var p TrendPoint
		if err := rows.Scan(&p.Period, &p.CO2T); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}

// Comparison reports baseline vs current totals and the reduction percentage.
type Comparison struct {
	Baseline     float64 `json:"baseline"`
	Current      float64 `json:"current"`
	ReductionPct float64 `json:"reduction_pct"`
}

func (r *AnalyticsRepo) Comparison(orgID int64) (*Comparison, error) {
	c := &Comparison{}
	// Baseline from the organization's base year, current from the latest year.
	q := `SELECT
		  COALESCE(SUM(CASE WHEN EXTRACT(YEAR FROM er.period_start)=o.base_year THEN er.co2_kg END),0)/1000.0,
		  COALESCE(SUM(CASE WHEN EXTRACT(YEAR FROM er.period_start)=o.reporting_year THEN er.co2_kg END),0)/1000.0
		FROM emission_records er
		JOIN facilities f ON f.id = er.facility_id
		JOIN organizations o ON o.id = f.organization_id
		WHERE f.organization_id = $1`
	if err := r.db.QueryRow(q, orgID).Scan(&c.Baseline, &c.Current); err != nil {
		return nil, err
	}
	if c.Baseline > 0 {
		c.ReductionPct = (c.Baseline - c.Current) / c.Baseline * 100
	}
	return c, nil
}

// FacilityBucket holds one facility's total emissions (tons).
type FacilityBucket struct {
	FacilityID   int64   `json:"facility_id"`
	FacilityName string  `json:"facility_name"`
	CO2T         float64 `json:"co2_t"`
}

func (r *AnalyticsRepo) FacilityBreakdown(orgID int64) ([]FacilityBucket, error) {
	q := `SELECT f.id, f.name, COALESCE(SUM(er.co2_kg),0)/1000.0
		FROM facilities f
		LEFT JOIN emission_records er ON er.facility_id = f.id
		WHERE f.organization_id = $1
		GROUP BY f.id, f.name ORDER BY co2_t DESC`
	rows, err := r.db.Query(q, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []FacilityBucket{}
	for rows.Next() {
		var b FacilityBucket
		if err := rows.Scan(&b.FacilityID, &b.FacilityName, &b.CO2T); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}
