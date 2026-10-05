CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    date_of_birth DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    muscle_group muscle_group NOT NULL,
    equipment equipment,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS schedule_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL,
    is_rest_day BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS schedule_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_day_id UUID NOT NULL REFERENCES schedule_days(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    sets INTEGER,
    reps INTEGER,
    weight DOUBLE PRECISION,
    rest INTEGER,
    duration INTEGER,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_exercises_deleted_at ON exercises(deleted_at);
CREATE INDEX IF NOT EXISTS idx_schedules_user_id ON schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_schedule_days_schedule_id ON schedule_days(schedule_id);
CREATE INDEX IF NOT EXISTS idx_schedule_exercises_day_id ON schedule_exercises(schedule_day_id);