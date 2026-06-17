-- Carbon Emission Management System - initial schema
-- PostgreSQL 16

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================================================
-- Organizations
-- =========================================================
CREATE TABLE organizations (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(200) NOT NULL,
    industry       VARCHAR(50) NOT NULL CHECK (industry IN ('manufacturing','logistics','energy','services','it','agriculture','construction','retail','other')),
    country        VARCHAR(100) NOT NULL,
    reporting_year INT NOT NULL,
    base_year      INT NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','archived')),
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================
-- Users
-- =========================================================
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(100) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    email           VARCHAR(200) NOT NULL UNIQUE,
    role            VARCHAR(20) NOT NULL DEFAULT 'viewer' CHECK (role IN ('admin','manager','analyst','viewer')),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_organization ON users(organization_id);

-- =========================================================
-- Facilities
-- =========================================================
CREATE TABLE facilities (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(200) NOT NULL,
    address         VARCHAR(500),
    latitude        DECIMAL(10,7),
    longitude       DECIMAL(10,7),
    type            VARCHAR(30) NOT NULL CHECK (type IN ('factory','office','warehouse','data_center','vehicle_fleet','retail','other')),
    country         VARCHAR(100),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_facilities_org ON facilities(organization_id);
CREATE INDEX idx_facilities_status ON facilities(status);

-- =========================================================
-- Emission Sources
-- =========================================================
CREATE TABLE emission_sources (
    id          BIGSERIAL PRIMARY KEY,
    facility_id BIGINT NOT NULL REFERENCES facilities(id),
    name        VARCHAR(200) NOT NULL,
    scope       INT NOT NULL CHECK (scope IN (1, 2, 3)),
    category    VARCHAR(50) NOT NULL CHECK (category IN ('stationary_combustion','mobile_combustion','electricity','heat','steam','business_travel','employee_commuting','purchased_goods','waste','upstream_transport','refrigerants')),
    fuel_type   VARCHAR(30) NOT NULL DEFAULT 'none' CHECK (fuel_type IN ('natural_gas','diesel','gasoline','coal','lpg','electricity','steam','refrigerant','none')),
    unit        VARCHAR(20) NOT NULL DEFAULT 'kwh' CHECK (unit IN ('m3','L','kg','kwh','tkm','GJ','t')),
    status      VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sources_facility ON emission_sources(facility_id);
CREATE INDEX idx_sources_scope ON emission_sources(scope);

-- =========================================================
-- Emission Factors
-- =========================================================
CREATE TABLE emission_factors (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    activity_unit VARCHAR(20) NOT NULL CHECK (activity_unit IN ('m3','L','kg','kwh','tkm','GJ','t')),
    factor_value  DECIMAL(14,6) NOT NULL,
    co2_unit      VARCHAR(5) NOT NULL DEFAULT 'kg' CHECK (co2_unit IN ('kg','t')),
    scope         INT NOT NULL CHECK (scope IN (1, 2, 3)),
    source_ref    VARCHAR(100),
    valid_from    TIMESTAMP,
    valid_to      TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_factors_scope ON emission_factors(scope);
CREATE INDEX idx_factors_source ON emission_factors(source_ref);

-- =========================================================
-- Emission Records
-- =========================================================
CREATE TABLE emission_records (
    id             BIGSERIAL PRIMARY KEY,
    source_id      BIGINT NOT NULL REFERENCES emission_sources(id),
    facility_id    BIGINT NOT NULL REFERENCES facilities(id),
    period         VARCHAR(20) NOT NULL DEFAULT 'monthly' CHECK (period IN ('monthly','quarterly','annual')),
    period_start   DATE NOT NULL,
    period_end     DATE NOT NULL,
    activity_value DECIMAL(16,4) NOT NULL,
    factor_id      BIGINT REFERENCES emission_factors(id),
    co2_kg         DECIMAL(16,4) NOT NULL DEFAULT 0,
    notes          TEXT,
    status         VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','submitted','verified','rejected')),
    calculated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_records_source ON emission_records(source_id);
CREATE INDEX idx_records_facility ON emission_records(facility_id);
CREATE INDEX idx_records_period ON emission_records(period_start, period_end);
CREATE INDEX idx_records_status ON emission_records(status);

-- =========================================================
-- Carbon Credits (offsets)
-- =========================================================
CREATE TABLE carbon_credits (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    type            VARCHAR(20) NOT NULL CHECK (type IN ('cer','eru','ver','rmuj')),
    project         VARCHAR(200),
    vintage_year    INT NOT NULL,
    amount_tons     DECIMAL(14,4) NOT NULL,
    price_per_ton   DECIMAL(12,2),
    status          VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available','reserved','retired')),
    retirement_date TIMESTAMP,
    registry_ref    VARCHAR(200),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_credits_status ON carbon_credits(status);
CREATE INDEX idx_credits_type ON carbon_credits(type);

-- =========================================================
-- Reduction Targets
-- =========================================================
CREATE TABLE reduction_targets (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    scope           INT NOT NULL DEFAULT 0 CHECK (scope IN (0, 1, 2, 3)),
    baseline_year   INT NOT NULL,
    target_year     INT NOT NULL,
    target_pct      DECIMAL(6,2) NOT NULL,
    baseline_co2_t  DECIMAL(14,4) NOT NULL DEFAULT 0,
    current_co2_t   DECIMAL(14,4) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'on_track' CHECK (status IN ('on_track','at_risk','achieved','missed')),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_targets_org ON reduction_targets(organization_id);

-- =========================================================
-- Carbon Reports
-- =========================================================
CREATE TABLE carbon_reports (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    period          VARCHAR(30) NOT NULL,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,
    total_co2_t     DECIMAL(16,4) NOT NULL DEFAULT 0,
    scope1_co2_t    DECIMAL(16,4) NOT NULL DEFAULT 0,
    scope2_co2_t    DECIMAL(16,4) NOT NULL DEFAULT 0,
    scope3_co2_t    DECIMAL(16,4) NOT NULL DEFAULT 0,
    offsets_t       DECIMAL(16,4) NOT NULL DEFAULT 0,
    net_co2_t       DECIMAL(16,4) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','generated','submitted','approved')),
    standard        VARCHAR(20) NOT NULL DEFAULT 'GHGP' CHECK (standard IN ('GHGP','ISO14064','CDP','CSRD')),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_reports_org ON carbon_reports(organization_id);
CREATE INDEX idx_reports_period ON carbon_reports(period_start, period_end);
CREATE UNIQUE INDEX idx_reports_org_period ON carbon_reports(organization_id, period);

-- =========================================================
-- Audit Logs
-- =========================================================
CREATE TABLE audit_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id),
    action     VARCHAR(30) NOT NULL CHECK (action IN ('create','update','delete','login','logout','generate','retire','verify','export')),
    entity     VARCHAR(50) NOT NULL,
    entity_id  BIGINT,
    detail     TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at);

-- =========================================================
-- Ops & management features (10 new tables)
-- =========================================================

-- 1. Data imports / exports
CREATE TABLE data_imports (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    type            VARCHAR(10) NOT NULL CHECK (type IN ('csv','excel')),
    direction       VARCHAR(10) NOT NULL DEFAULT 'import' CHECK (direction IN ('import','export')),
    target_entity   VARCHAR(50) NOT NULL,
    file_path       VARCHAR(500),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','completed','failed')),
    records_count   INT NOT NULL DEFAULT 0,
    error_msg       TEXT,
    started_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_data_imports_org ON data_imports(organization_id);
CREATE INDEX idx_data_imports_status ON data_imports(status);

-- 2. Scheduled tasks
CREATE TABLE scheduled_tasks (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(200) NOT NULL,
    cron_expr       VARCHAR(50) NOT NULL,
    target_endpoint VARCHAR(500) NOT NULL,
    payload         TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','paused')),
    last_run        TIMESTAMP,
    next_run        TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_tasks_org ON scheduled_tasks(organization_id);
CREATE INDEX idx_tasks_status ON scheduled_tasks(status);

-- 3. Alerts
CREATE TABLE alerts (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    facility_id     BIGINT REFERENCES facilities(id),
    type            VARCHAR(20) NOT NULL CHECK (type IN ('threshold','anomaly','trend')),
    severity        VARCHAR(10) NOT NULL DEFAULT 'warning' CHECK (severity IN ('info','warning','critical')),
    message         TEXT NOT NULL,
    trigger_value   DECIMAL(16,4) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','acknowledged','resolved')),
    acked_by        BIGINT REFERENCES users(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_alerts_org ON alerts(organization_id);
CREATE INDEX idx_alerts_status ON alerts(status);

-- 4. Notifications
CREATE TABLE notifications (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    type            VARCHAR(10) NOT NULL CHECK (type IN ('email','sms','webhook')),
    recipient       VARCHAR(200) NOT NULL,
    subject         VARCHAR(200),
    body            TEXT,
    channel         VARCHAR(20) NOT NULL DEFAULT 'system' CHECK (channel IN ('system','alert','report')),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sent','failed')),
    retries         INT NOT NULL DEFAULT 0,
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_notifications_org ON notifications(organization_id);
CREATE INDEX idx_notifications_status ON notifications(status);

-- 5. API keys
CREATE TABLE api_keys (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(200) NOT NULL,
    key_hash        VARCHAR(64) NOT NULL UNIQUE,
    key_prefix      VARCHAR(20) NOT NULL,
    scopes          VARCHAR(200) NOT NULL DEFAULT 'read',
    expires_at      TIMESTAMP,
    last_used       TIMESTAMP,
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','revoked')),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_api_keys_org ON api_keys(organization_id);

-- 6. Webhooks
CREATE TABLE webhooks (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    url             VARCHAR(500) NOT NULL,
    events          VARCHAR(500) NOT NULL,
    secret          VARCHAR(200),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
    failure_count   INT NOT NULL DEFAULT 0,
    last_fired      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_webhooks_org ON webhooks(organization_id);

-- 7. Attachments
CREATE TABLE attachments (
    id          BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   BIGINT NOT NULL,
    filename    VARCHAR(255) NOT NULL,
    file_path   VARCHAR(500) NOT NULL,
    file_size   BIGINT NOT NULL DEFAULT 0,
    mime_type   VARCHAR(100),
    uploaded_by BIGINT REFERENCES users(id),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_attachments_entity ON attachments(entity_type, entity_id);

-- 8. Report exports
CREATE TABLE report_exports (
    id           BIGSERIAL PRIMARY KEY,
    report_id    BIGINT NOT NULL REFERENCES carbon_reports(id),
    format       VARCHAR(10) NOT NULL CHECK (format IN ('pdf','excel')),
    file_path    VARCHAR(500),
    options      TEXT,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','generated','failed')),
    generated_by BIGINT REFERENCES users(id),
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_report_exports_report ON report_exports(report_id);

-- 9. Rollback records
CREATE TABLE rollback_records (
    id            BIGSERIAL PRIMARY KEY,
    audit_log_id  BIGINT NOT NULL REFERENCES audit_logs(id),
    entity        VARCHAR(50) NOT NULL,
    entity_id     BIGINT NOT NULL,
    snapshot      TEXT NOT NULL,
    rolled_back_by BIGINT REFERENCES users(id),
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_rollbacks_entity ON rollback_records(entity, entity_id);

-- 10. System settings
CREATE TABLE system_settings (
    id          BIGSERIAL PRIMARY KEY,
    key         VARCHAR(100) NOT NULL UNIQUE,
    value       TEXT NOT NULL,
    category    VARCHAR(50) NOT NULL DEFAULT 'general',
    description VARCHAR(500),
    data_type   VARCHAR(10) NOT NULL DEFAULT 'string' CHECK (data_type IN ('string','int','float','bool','json')),
    updated_by  BIGINT REFERENCES users(id),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_settings_category ON system_settings(category);

-- =========================================================
-- Seed data
-- =========================================================
INSERT INTO organizations (name, industry, country, reporting_year, base_year)
VALUES ('GreenTech Industries', 'manufacturing', 'China', 2025, 2020);

INSERT INTO users (username, password, email, role, organization_id)
VALUES ('admin', '$2a$10$wW8.V3DjX8m9oNqZ4pK1kueQ0sR5tYvL7mZ2kJ.HxYcBnW3pKqL.C', 'admin@carbon.io', 'admin', 1);

INSERT INTO facilities (organization_id, name, address, latitude, longitude, type, country)
VALUES (1, 'Headquarters Plant', '88 Industrial Road, Shanghai', 31.2304000, 121.4737000, 'factory', 'China');

INSERT INTO emission_factors (name, activity_unit, factor_value, co2_unit, scope, source_ref)
VALUES
    ('Natural Gas Combustion',  'm3',   2.02110, 'kg', 1, 'IPCC 2006'),
    ('Diesel Combustion',       'L',    2.68200, 'kg', 1, 'IPCC 2006'),
    ('Grid Electricity (China)','kwh',  0.58100, 'kg', 2, 'MEE 2024'),
    ('Business Travel - Air',   'tkm',  0.18000, 'kg', 3, 'DEFRA 2024');

INSERT INTO emission_sources (facility_id, name, scope, category, fuel_type, unit)
VALUES
    (1, 'Boiler A - Natural Gas', 1, 'stationary_combustion', 'natural_gas', 'm3'),
    (1, 'Forklift Fleet - Diesel', 1, 'mobile_combustion', 'diesel', 'L'),
    (1, 'Purchased Electricity',   2, 'electricity', 'electricity', 'kwh');

-- System settings seed
INSERT INTO system_settings (key, value, category, description, data_type) VALUES
    ('default_reporting_standard', 'GHGP', 'general',     'Default ESG disclosure standard',  'string'),
    ('alert_threshold_pct',        '10',   'alert',       'YoY emission increase alert (%)',  'int'),
    ('notification_enabled',       'true', 'notification','Master switch for notifications',  'bool'),
    ('auto_generate_reports',      'false','report',      'Auto-generate monthly reports',    'bool');

-- Sample scheduled task
INSERT INTO scheduled_tasks (organization_id, name, cron_expr, target_endpoint, status)
VALUES (1, 'Monthly Emission Rollup', '0 0 1 * *', '/api/v1/internal/rollup', 'active');

-- Sample alert
INSERT INTO alerts (organization_id, facility_id, type, severity, message, trigger_value)
VALUES (1, 1, 'threshold', 'warning', 'Natural gas consumption exceeded monthly budget by 12%', 12.0);
