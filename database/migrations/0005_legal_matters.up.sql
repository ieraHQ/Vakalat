-- Matters & Hearings
CREATE TABLE matters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    client_id UUID REFERENCES clients(id),
    court_id UUID REFERENCES courts(id),
    judge_id UUID REFERENCES judges(id),
    advocate_id UUID REFERENCES advocates(id),
    case_number TEXT,
    case_type TEXT,
    stage TEXT,
    status TEXT,
    priority TEXT,
    limitation_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE matter_parties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    party_type TEXT CHECK (party_type IN ('petitioner', 'respondent', 'intervenor')),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE hearings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    date TIMESTAMPTZ NOT NULL,
    notes TEXT,
    outcome TEXT,
    next_hearing_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    hearing_id UUID REFERENCES hearings(id),
    title TEXT NOT NULL,
    description TEXT,
    date TIMESTAMPTZ NOT NULL,
    document_id UUID REFERENCES documents(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);