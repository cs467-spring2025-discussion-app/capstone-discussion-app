DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'fastapi_test') THEN
      CREATE USER fastapi_test WITH PASSWORD 'fastapi_test';
   END IF;
END
$do$;

CREATE DATABASE fastapi_test WITH OWNER fastapi_test ENCODING 'UTF8';

\connect fastapi_test

ALTER DATABASE fastapi_test SET timezone TO 'UTC';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

GRANT ALL PRIVILEGES ON DATABASE fastapi_test TO fastapi_test;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO fastapi_test;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO fastapi_test;
GRANT ALL PRIVILEGES ON SCHEMA public TO fastapi_test;

ALTER USER fastapi_test WITH SUPERUSER;
