SELECT table_name FROM information_schema.tables
WHERE table_schema = $1 AND table_name IN (%s)