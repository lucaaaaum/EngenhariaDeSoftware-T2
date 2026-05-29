CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    assigned_to UUID,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status INT NOT NULL
);
