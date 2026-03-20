ALTER TABLE scan_jobs DROP CONSTRAINT IF EXISTS scan_jobs_scan_type_check;
ALTER TABLE scan_jobs
ADD CONSTRAINT scan_jobs_scan_type_check
CHECK (scan_type IN ('subdomain','dns','whois','port','asn','ssl','all','ip','tech'));

CREATE TABLE IF NOT EXISTS ip_scan_results (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    ip_address VARCHAR(64) NOT NULL,
    geolocation_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    asn_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    reverse_dns TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ip_scan_results_asset_id ON ip_scan_results(asset_id);
CREATE INDEX IF NOT EXISTS idx_ip_scan_results_scan_job_id ON ip_scan_results(scan_job_id);

CREATE TABLE IF NOT EXISTS port_scan_results (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    ip_address VARCHAR(64) NOT NULL,
    open_ports_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    closed_ports INTEGER NOT NULL DEFAULT 0,
    total_scanned INTEGER NOT NULL DEFAULT 0,
    scan_duration_ms INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_port_scan_results_asset_id ON port_scan_results(asset_id);
CREATE INDEX IF NOT EXISTS idx_port_scan_results_scan_job_id ON port_scan_results(scan_job_id);