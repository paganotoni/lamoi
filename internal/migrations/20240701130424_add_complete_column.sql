-- 20240701130424 - add_complete_column migration
ALTER TABLE messages ADD COLUMN complete BOOLEAN DEFAULT TRUE;
