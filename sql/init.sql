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
-- Seed data
-- =========================================================
INSERT INTO organizations (name, industry, country, reporting_year, base_year)
VALUES ('GreenTech Industries', 'manufacturing', 'China', 2025, 2020);

INSERT INTO users (username, password, email, role, organization_id)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@carbon.io', 'admin', 1);

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
