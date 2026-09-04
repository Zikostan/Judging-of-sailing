-- ============================================================
-- Migration 001: Rollback
-- ============================================================

DROP TRIGGER IF EXISTS trg_user_updated_at    ON "user";
DROP TRIGGER IF EXISTS trg_regatta_updated_at ON regatta;
DROP TRIGGER IF EXISTS trg_sailor_updated_at  ON sailor;
DROP TRIGGER IF EXISTS trg_document_updated_at ON document;
DROP TRIGGER IF EXISTS trg_race_group_updated_at ON race_group;
DROP TRIGGER IF EXISTS trg_race_updated_at    ON race;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS sailor_category;
DROP TABLE IF EXISTS race;
DROP TABLE IF EXISTS race_group;
DROP TABLE IF EXISTS document;
DROP TABLE IF EXISTS sailor;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS regatta;
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS race_group_status;
DROP TYPE IF EXISTS document_status;
DROP TYPE IF EXISTS user_role;