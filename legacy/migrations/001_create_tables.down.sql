-- Drop tables (employees first due to foreign key constraint)
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";