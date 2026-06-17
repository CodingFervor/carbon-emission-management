package model

import "time"

// User represents an authenticated account within an organization.
type User struct {
	ID             int64     `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Password       string    `json:"-" db:"password"`
	Email          string    `json:"email" db:"email"`
	Role           string    `json:"role" db:"role"` // admin, manager, analyst, viewer
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Status         string    `json:"status" db:"status"` // active, disabled
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Organization is a reporting entity (a company or business unit).
type Organization struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Industry      string    `json:"industry" db:"industry"` // manufacturing, logistics, energy, services, it
	Country       string    `json:"country" db:"country"`
	ReportingYear int       `json:"reporting_year" db:"reporting_year"`
	BaseYear      int       `json:"base_year" db:"base_year"`
	Status        string    `json:"status" db:"status"` // active, archived
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Facility is a physical site where emissions originate.
type Facility struct {
	ID             int64     `json:"id" db:"id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Address        string    `json:"address" db:"address"`
	Latitude       float64   `json:"latitude" db:"latitude"`
	Longitude      float64   `json:"longitude" db:"longitude"`
	Type           string    `json:"type" db:"type"` // factory, office, warehouse, data_center, vehicle_fleet
	Country        string    `json:"country" db:"country"`
	Status         string    `json:"status" db:"status"` // active, inactive
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// EmissionSource is a categorized origin of greenhouse gas emissions.
type EmissionSource struct {
	ID         int64     `json:"id" db:"id"`
	FacilityID int64     `json:"facility_id" db:"facility_id"`
	Name       string    `json:"name" db:"name"`
	Scope      int       `json:"scope" db:"scope"`         // 1, 2, 3
	Category   string    `json:"category" db:"category"`   // stationary_combustion, mobile_combustion, electricity, heat, business_travel, employee_commuting, purchased_goods, waste, upstream_transport
	FuelType   string    `json:"fuel_type" db:"fuel_type"` // natural_gas, diesel, gasoline, coal, electricity, none
	Unit       string    `json:"unit" db:"unit"`           // m3, L, kg, kWh, tkm
	Status     string    `json:"status" db:"status"`       // active, inactive
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// EmissionFactor maps an activity quantity to CO2-equivalent emissions.
type EmissionFactor struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	ActivityUnit string     `json:"activity_unit" db:"activity_unit"` // m3, L, kg, kWh, tkm
	FactorValue  float64    `json:"factor_value" db:"factor_value"`   // kg CO2e per activity unit
	CO2Unit      string     `json:"co2_unit" db:"co2_unit"`           // kg, t
	Scope        int        `json:"scope" db:"scope"`                 // 1, 2, 3
	SourceRef    string     `json:"source_ref" db:"source_ref"`       // IPCC, DEFRA, EPA, custom
	ValidFrom    *time.Time `json:"valid_from" db:"valid_from"`
	ValidTo      *time.Time `json:"valid_to" db:"valid_to"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// EmissionRecord is a single measured or calculated emission data point.
type EmissionRecord struct {
	ID            int64     `json:"id" db:"id"`
	SourceID      int64     `json:"source_id" db:"source_id"`
	FacilityID    int64     `json:"facility_id" db:"facility_id"`
	Period        string    `json:"period" db:"period"` // monthly, annual
	PeriodStart   time.Time `json:"period_start" db:"period_start"`
	PeriodEnd     time.Time `json:"period_end" db:"period_end"`
	ActivityValue float64   `json:"activity_value" db:"activity_value"`
	FactorID      *int64    `json:"factor_id" db:"factor_id"`
	CO2Kg         float64   `json:"co2_kg" db:"co2_kg"`
	Notes         string    `json:"notes" db:"notes"`
	Status        string    `json:"status" db:"status"` // draft, submitted, verified
	CalculatedAt  time.Time `json:"calculated_at" db:"calculated_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CarbonCredit is an offset unit that can be retired against emissions.
type CarbonCredit struct {
	ID             int64      `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Type           string     `json:"type" db:"type"` // cer, eru, ver, rmuj
	Project        string     `json:"project" db:"project"`
	VintageYear    int        `json:"vintage_year" db:"vintage_year"`
	AmountTons     float64    `json:"amount_tons" db:"amount_tons"`
	PricePerTon    float64    `json:"price_per_ton" db:"price_per_ton"`
	Status         string     `json:"status" db:"status"` // available, reserved, retired
	RetirementDate *time.Time `json:"retirement_date" db:"retirement_date"`
	RegistryRef    string     `json:"registry_ref" db:"registry_ref"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ReductionTarget defines a goal to cut emissions versus a baseline.
type ReductionTarget struct {
	ID             int64     `json:"id" db:"id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Scope          int       `json:"scope" db:"scope"` // 1, 2, 3 (0 = all scopes)
	BaselineYear   int       `json:"baseline_year" db:"baseline_year"`
	TargetYear     int       `json:"target_year" db:"target_year"`
	TargetPct      float64   `json:"target_pct" db:"target_pct"`
	BaselineCO2T   float64   `json:"baseline_co2_t" db:"baseline_co2_t"`
	CurrentCO2T    float64   `json:"current_co2_t" db:"current_co2_t"`
	Status         string    `json:"status" db:"status"` // on_track, at_risk, achieved, missed
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CarbonReport is a compiled ESG / carbon disclosure for a period.
type CarbonReport struct {
	ID             int64     `json:"id" db:"id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Period         string    `json:"period" db:"period"` // 2025, Q1-2025
	PeriodStart    time.Time `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time `json:"period_end" db:"period_end"`
	TotalCO2T      float64   `json:"total_co2_t" db:"total_co2_t"`
	Scope1CO2T     float64   `json:"scope1_co2_t" db:"scope1_co2_t"`
	Scope2CO2T     float64   `json:"scope2_co2_t" db:"scope2_co2_t"`
	Scope3CO2T     float64   `json:"scope3_co2_t" db:"scope3_co2_t"`
	OffsetsT       float64   `json:"offsets_t" db:"offsets_t"`
	NetCO2T        float64   `json:"net_co2_t" db:"net_co2_t"`
	Status         string    `json:"status" db:"status"`     // draft, generated, submitted, approved
	Standard       string    `json:"standard" db:"standard"` // GHGP, ISO14064, CDP, CSRD
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// AuditLog records a user action against a domain entity.
type AuditLog struct {
	ID        int64     `json:"id" db:"id"`
	UserID    *int64    `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"` // create, update, delete, login, generate, retire
	Entity    string    `json:"entity" db:"entity"` // organization, facility, emission_record, carbon_report, ...
	EntityID  *int64    `json:"entity_id" db:"entity_id"`
	Detail    string    `json:"detail" db:"detail"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// =========================================================
// Ops & management features (10 new features)
// =========================================================

// DataImport tracks a bulk CSV/Excel import or export job.
type DataImport struct {
	ID             int64     `json:"id" db:"id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Type           string    `json:"type" db:"type"`                   // csv, excel
	Direction      string    `json:"direction" db:"direction"`         // import, export
	TargetEntity   string    `json:"target_entity" db:"target_entity"` // facilities, emission_records, ...
	FilePath       string    `json:"file_path" db:"file_path"`
	Status         string    `json:"status" db:"status"` // pending, processing, completed, failed
	RecordsCount   int       `json:"records_count" db:"records_count"`
	ErrorMsg       string    `json:"error_msg" db:"error_msg"`
	StartedBy      *int64    `json:"started_by" db:"started_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ScheduledTask defines a recurring job driven by a cron expression.
type ScheduledTask struct {
	ID             int64      `json:"id" db:"id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	Name           string     `json:"name" db:"name"`
	CronExpr       string     `json:"cron_expr" db:"cron_expr"`
	TargetEndpoint string     `json:"target_endpoint" db:"target_endpoint"`
	Payload        string     `json:"payload" db:"payload"`
	Status         string     `json:"status" db:"status"` // active, paused
	LastRun        *time.Time `json:"last_run" db:"last_run"`
	NextRun        *time.Time `json:"next_run" db:"next_run"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Alert is a triggered threshold, anomaly, or trend notification.
type Alert struct {
	ID             int64     `json:"id" db:"id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	FacilityID     *int64    `json:"facility_id" db:"facility_id"`
	Type           string    `json:"type" db:"type"`         // threshold, anomaly, trend
	Severity       string    `json:"severity" db:"severity"` // info, warning, critical
	Message        string    `json:"message" db:"message"`
	TriggerValue   float64   `json:"trigger_value" db:"trigger_value"`
	Status         string    `json:"status" db:"status"` // active, acknowledged, resolved
	AckedBy        *int64    `json:"acked_by" db:"acked_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Notification is an outbound email/sms/webhook message.
type Notification struct {
	ID             int64      `json:"id" db:"id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	Type           string     `json:"type" db:"type"` // email, sms, webhook
	Recipient      string     `json:"recipient" db:"recipient"`
	Subject        string     `json:"subject" db:"subject"`
	Body           string     `json:"body" db:"body"`
	Channel        string     `json:"channel" db:"channel"` // system, alert, report
	Status         string     `json:"status" db:"status"`   // pending, sent, failed
	Retries        int        `json:"retries" db:"retries"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// APIKey grants programmatic access scoped to an organization.
type APIKey struct {
	ID             int64      `json:"id" db:"id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	Name           string     `json:"name" db:"name"`
	KeyHash        string     `json:"-" db:"key_hash"`
	KeyPrefix      string     `json:"key_prefix" db:"key_prefix"`
	Scopes         string     `json:"scopes" db:"scopes"` // comma-separated, e.g. read,write
	ExpiresAt      *time.Time `json:"expires_at" db:"expires_at"`
	LastUsed       *time.Time `json:"last_used" db:"last_used"`
	Status         string     `json:"status" db:"status"` // active, revoked
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Webhook is a subscription that receives event payloads at an external URL.
type Webhook struct {
	ID             int64      `json:"id" db:"id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	URL            string     `json:"url" db:"url"`
	Events         string     `json:"events" db:"events"` // emission.created,report.generated,...
	Secret         string     `json:"-" db:"secret"`
	Status         string     `json:"status" db:"status"` // active, disabled
	FailureCount   int        `json:"failure_count" db:"failure_count"`
	LastFired      *time.Time `json:"last_fired" db:"last_fired"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Attachment is a file linked to any domain entity (polymorphic).
type Attachment struct {
	ID         int64     `json:"id" db:"id"`
	EntityType string    `json:"entity_type" db:"entity_type"` // report, emission_record, ...
	EntityID   int64     `json:"entity_id" db:"entity_id"`
	Filename   string    `json:"filename" db:"filename"`
	FilePath   string    `json:"file_path" db:"file_path"`
	FileSize   int64     `json:"file_size" db:"file_size"`
	MimeType   string    `json:"mime_type" db:"mime_type"`
	UploadedBy *int64    `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ReportExport is a rendered (PDF/Excel) export of a carbon report.
type ReportExport struct {
	ID          int64     `json:"id" db:"id"`
	ReportID    int64     `json:"report_id" db:"report_id"`
	Format      string    `json:"format" db:"format"` // pdf, excel
	FilePath    string    `json:"file_path" db:"file_path"`
	Options     string    `json:"options" db:"options"` // JSON string of render options
	Status      string    `json:"status" db:"status"`   // pending, generated, failed
	GeneratedBy *int64    `json:"generated_by" db:"generated_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// RollbackRecord captures an entity snapshot so a destructive change can be undone.
type RollbackRecord struct {
	ID           int64     `json:"id" db:"id"`
	AuditLogID   int64     `json:"audit_log_id" db:"audit_log_id"`
	Entity       string    `json:"entity" db:"entity"`
	EntityID     int64     `json:"entity_id" db:"entity_id"`
	Snapshot     string    `json:"snapshot" db:"snapshot"` // JSON of pre-change state
	RolledBackBy *int64    `json:"rolled_back_by" db:"rolled_back_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// SystemSetting is a key/value configuration entry grouped by category.
type SystemSetting struct {
	ID          int64     `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Category    string    `json:"category" db:"category"` // general, notification, factor, ...
	Description string    `json:"description" db:"description"`
	DataType    string    `json:"data_type" db:"data_type"` // string, int, float, bool, json
	UpdatedBy   *int64    `json:"updated_by" db:"updated_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
