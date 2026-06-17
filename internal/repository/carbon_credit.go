package repository

import (
	"database/sql"
	"time"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type CarbonCreditRepo struct {
	*GenericRepo[model.CarbonCredit]
	db *sql.DB
}

func NewCarbonCreditRepo(db *sql.DB) *CarbonCreditRepo {
	return &CarbonCreditRepo{
		GenericRepo: NewGenericRepo[model.CarbonCredit](db, "carbon_credits", func() *model.CarbonCredit { return &model.CarbonCredit{} }),
		db:          db,
	}
}

func (r *CarbonCreditRepo) Create(c *model.CarbonCredit) error {
	q := `INSERT INTO carbon_credits (name, type, project, vintage_year, amount_tons, price_per_ton, status, registry_ref)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q,
		c.Name, c.Type, c.Project, c.VintageYear, c.AmountTons, c.PricePerTon,
		defaultStr(c.Status, "available"), c.RegistryRef,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CarbonCreditRepo) Update(c *model.CarbonCredit) error {
	q := `UPDATE carbon_credits SET name=$1, type=$2, project=$3, vintage_year=$4, amount_tons=$5,
	      price_per_ton=$6, registry_ref=$7, updated_at=CURRENT_TIMESTAMP WHERE id=$8`
	_, err := r.db.Exec(q, c.Name, c.Type, c.Project, c.VintageYear, c.AmountTons, c.PricePerTon, c.RegistryRef, c.ID)
	return err
}

// Retire marks a credit as retired, recording the retirement date.
func (r *CarbonCreditRepo) Retire(id int64) error {
	_, err := r.db.Exec(
		"UPDATE carbon_credits SET status='retired', retirement_date=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2 AND status='available'",
		time.Now(), id,
	)
	return err
}
