-- Documents & OCR
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    mime_type TEXT,
    size BIGINT,
    hash TEXT,
    ocr_status TEXT CHECK (ocr_status IN ('pending', 'processing', 'completed', 'failed')),
    ocr_text TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID REFERENCES documents(id),
    version INTEGER NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID REFERENCES documents(id),
    embedding VECTOR(384),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE document_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID REFERENCES documents(id),
    tag TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE ocr_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID REFERENCES documents(id),
    status TEXT CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    error TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);