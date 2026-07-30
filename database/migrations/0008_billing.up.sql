-- Billing & Expenses
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    amount NUMERIC NOT NULL,
    tax NUMERIC,
    total NUMERIC NOT NULL,
    status TEXT CHECK (status IN ('draft', 'sent', 'paid', 'overdue', 'cancelled')),
    due_date TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID REFERENCES invoices(id),
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    rate NUMERIC NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    description TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    category TEXT,
    receipt_path TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID REFERENCES invoices(id),
    amount NUMERIC NOT NULL,
    payment_method TEXT,
    transaction_id TEXT,
    date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);