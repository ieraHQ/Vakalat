-- AI & Search
CREATE TABLE ai_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    prompt TEXT NOT NULL,
    response TEXT,
    model TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE ai_summaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    summary TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE search_index (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    content TEXT,
    embedding VECTOR(384),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    query TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);