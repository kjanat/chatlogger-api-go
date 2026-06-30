-- Migration for exports table to support async export functionality

-- Create exports table
CREATE TABLE IF NOT EXISTS exports (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    format VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    file_path TEXT,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,

    CONSTRAINT fk_exports_organization FOREIGN KEY (organization_id)
        REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_exports_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

-- Add indexes for quicker lookup
CREATE INDEX idx_exports_organization_id ON exports(organization_id);
CREATE INDEX idx_exports_user_id ON exports(user_id);
CREATE INDEX idx_exports_status ON exports(status);
CREATE INDEX idx_exports_created_at ON exports(created_at);

-- Create trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_exports_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_exports_updated_at
BEFORE UPDATE ON exports
FOR EACH ROW
EXECUTE FUNCTION update_exports_updated_at();

COMMENT ON TABLE exports IS 'Stores asynchronous export jobs with status tracking';
COMMENT ON COLUMN exports.format IS 'Export format: json, csv';
COMMENT ON COLUMN exports.type IS 'Export type: chats, messages, all';
COMMENT ON COLUMN exports.status IS 'Current status: pending, processing, completed, failed';
