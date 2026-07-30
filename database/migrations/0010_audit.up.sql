-- Audit & Compliance
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE backup_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    backup_path TEXT NOT NULL,
    status TEXT CHECK (status IN ('pending', 'completed', 'failed')),
    size BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);