CREATE TABLE IF NOT EXISTS tasks (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);

CREATE TABLE IF NOT EXISTS recurrences (
	id BIGSERIAL PRIMARY KEY,

	type TEXT NOT NULL,
	interval_days INT,
	interval_months INT,
	specific_day DATE[],
	even_odd_days BOOL,

	start_date TIMESTAMPTZ,
	end_date TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_recurrence_type ON recurrences (type);

CREATE TABLE IF NOT EXISTS recurrence_on_task (
	task_id BIGINT NOT NULL REFERENCES tasks(id),
	recurrence_id BIGINT NOT NULL REFERENCES recurrences(id)
);

CREATE INDEX IF NOT EXISTS idx_recurrence_on_task_task_id ON recurrence_on_task (task_id, recurrence_id);	