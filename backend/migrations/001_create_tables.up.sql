-- ============================================================
-- Migration 001: Create all core tables
-- ============================================================

-- 1. Custom enums (Pg не поддерживает IF NOT EXISTS для типов, используем DO)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('judge', 'secretary', 'admin');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_status') THEN
        CREATE TYPE document_status AS ENUM ('pending', 'approved', 'rejected', 'need_revision');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'race_group_status') THEN
        CREATE TYPE race_group_status AS ENUM ('scheduled', 'running', 'finished', 'cancelled');
    END IF;
END $$;

-- 2. User (judge, secretary, admin)
CREATE TABLE IF NOT EXISTS "user" (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          user_role   NOT NULL,
    last_name     TEXT        NOT NULL,
    first_name    TEXT        NOT NULL,
    middle_name   TEXT,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3. Regatta (competition event)
CREATE TABLE IF NOT EXISTS regatta (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT         NOT NULL,
    start_date  DATE         NOT NULL,
    end_date    DATE         NOT NULL,
    location    TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CONSTRAINT chk_regatta_dates CHECK (end_date >= start_date)
);

-- 4. Category (scoring division: age + gender + boat class)
CREATE TABLE IF NOT EXISTS category (
    id          UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT    NOT NULL,
    gender      TEXT    NOT NULL CHECK (gender IN ('male', 'female', 'mixed')),
    age_from    INT     NOT NULL CHECK (age_from >= 0),
    age_to      INT     NOT NULL CHECK (age_to >= age_from),
    boat_class  TEXT    NOT NULL,
    UNIQUE (name, boat_class)
);

-- 5. Sailor (participant / athlete)
CREATE TABLE IF NOT EXISTS sailor (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    last_name        TEXT        NOT NULL,
    first_name       TEXT        NOT NULL,
    middle_name      TEXT,
    gender           TEXT        NOT NULL CHECK (gender IN ('male', 'female')),
    birth_date       DATE        NOT NULL,
    boat_class       TEXT        NOT NULL,
    sail_number      TEXT        NOT NULL,
    regatta_id       UUID        NOT NULL REFERENCES regatta(id) ON DELETE CASCADE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (regatta_id, boat_class, sail_number)
);

CREATE INDEX IF NOT EXISTS idx_sailor_regatta_class_number
    ON sailor (regatta_id, boat_class, sail_number);

-- 6. Document (sailor's paperwork)
CREATE TABLE IF NOT EXISTS document (
    id          UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    sailor_id   UUID             NOT NULL REFERENCES sailor(id) ON DELETE CASCADE,
    type        TEXT             NOT NULL,
    file_url    TEXT             NOT NULL,
    status      document_status  NOT NULL DEFAULT 'pending',
    checked_by  UUID             REFERENCES "user"(id) ON DELETE SET NULL,
    comment     TEXT,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_document_sailor ON document (sailor_id);

-- 7. Race group (a single start session with multiple boats)
CREATE TABLE IF NOT EXISTS race_group (
    id          UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    regatta_id  UUID               NOT NULL REFERENCES regatta(id) ON DELETE CASCADE,
    category_id UUID               NOT NULL REFERENCES category(id),
    number      INT                NOT NULL,
    start_at    TIMESTAMPTZ        NOT NULL,
    status      race_group_status  NOT NULL DEFAULT 'scheduled',
    created_at  TIMESTAMPTZ        NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ        NOT NULL DEFAULT now(),
    UNIQUE (regatta_id, number)
);

CREATE INDEX IF NOT EXISTS idx_race_group_regatta ON race_group (regatta_id);

-- 8. Race (individual boat measurement)
CREATE TABLE IF NOT EXISTS race (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id    UUID        NOT NULL REFERENCES race_group(id) ON DELETE CASCADE,
    sailor_id   UUID        REFERENCES sailor(id) ON DELETE SET NULL,
    judge_id    UUID        NOT NULL REFERENCES "user"(id),
    boat_class  TEXT        NOT NULL,
    sail_number TEXT        NOT NULL,
    finish_time REAL,
    place       SMALLINT,
    penalty     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (group_id, boat_class, sail_number)
);

CREATE INDEX IF NOT EXISTS idx_race_group ON race (group_id);
CREATE INDEX IF NOT EXISTS idx_race_sailor ON race (sailor_id);
CREATE INDEX IF NOT EXISTS idx_race_class_number ON race (group_id, boat_class, sail_number);

-- 9. Sailor-Category link (many-to-many)
CREATE TABLE IF NOT EXISTS sailor_category (
    sailor_id   UUID NOT NULL REFERENCES sailor(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    PRIMARY KEY (sailor_id, category_id)
);

-- 10. Auto-update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_user_updated_at ON "user";
CREATE TRIGGER trg_user_updated_at
    BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_regatta_updated_at ON regatta;
CREATE TRIGGER trg_regatta_updated_at
    BEFORE UPDATE ON regatta FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_sailor_updated_at ON sailor;
CREATE TRIGGER trg_sailor_updated_at
    BEFORE UPDATE ON sailor FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_document_updated_at ON document;
CREATE TRIGGER trg_document_updated_at
    BEFORE UPDATE ON document FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_race_group_updated_at ON race_group;
CREATE TRIGGER trg_race_group_updated_at
    BEFORE UPDATE ON race_group FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_race_updated_at ON race;
CREATE TRIGGER trg_race_updated_at
    BEFORE UPDATE ON race FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();