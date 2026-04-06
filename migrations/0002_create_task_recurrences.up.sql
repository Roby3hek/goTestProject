CREATE TABLE IF NOT EXISTS task_recurrences (
                                                task_id BIGINT PRIMARY KEY
                                                REFERENCES tasks(id) ON DELETE CASCADE,

    recurrence_type TEXT NOT NULL DEFAULT 'none',

    recurrence_interval INTEGER NOT NULL DEFAULT 1
    CHECK (recurrence_interval >= 0),

    recurrence_day_of_month INTEGER
    CHECK (recurrence_day_of_month BETWEEN 1 AND 31),

    specific_dates TEXT,

    start_date TIMESTAMPTZ NOT NULL,
    end_date   TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_valid_recurrence_type
    CHECK (recurrence_type IN ('none', 'daily', 'monthly', 'even', 'odd', 'specific')),

    CONSTRAINT chk_monthly_requires_day
    CHECK (recurrence_type != 'monthly' OR recurrence_day_of_month IS NOT NULL),

    CONSTRAINT chk_specific_requires_dates
    CHECK (recurrence_type != 'specific' OR (specific_dates IS NOT NULL AND specific_dates != ''))
);

CREATE INDEX IF NOT EXISTS idx_task_recurrences_type ON task_recurrences (recurrence_type);
CREATE INDEX IF NOT EXISTS idx_task_recurrences_start_date ON task_recurrences (start_date);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_task_recurrences_updated_at
    BEFORE UPDATE ON task_recurrences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();