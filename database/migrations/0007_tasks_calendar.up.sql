-- Tasks & Calendar
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    title TEXT NOT NULL,
    description TEXT,
    due_date TIMESTAMPTZ,
    priority TEXT,
    status TEXT CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    assigned_to UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE calendar_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    matter_id UUID REFERENCES matters(id),
    title TEXT NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    location TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE reminders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    task_id UUID REFERENCES tasks(id),
    event_id UUID REFERENCES calendar_events(id),
    reminder_time TIMESTAMPTZ NOT NULL,
    notified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);