-- 20240713172524 - pragmas migration
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 2000;
PRAGMA temp_store = memory;
