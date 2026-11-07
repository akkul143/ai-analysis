CREATE TABLE runs (
    id UUID PRIMARY KEY,
    mode TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('running', 'completed', 'failed')),
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX runs_started_at_idx ON runs (started_at DESC);
CREATE INDEX runs_status_idx ON runs (status);
