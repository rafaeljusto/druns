REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM druns;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM druns;
REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public FROM druns;

DROP USER IF EXISTS druns;
CREATE USER druns WITH PASSWORD 'abc123';

/****************************************/

DROP TYPE IF EXISTS LOG_OPERATION CASCADE;
CREATE TYPE LOG_OPERATION AS ENUM ('CREATE', 'UPDATE', 'DELETE');

DROP TABLE IF EXISTS log CASCADE;
CREATE TABLE log (
	id SERIAL PRIMARY KEY,
	agent INT REFERENCES adm_user(id),
	ip_address INET,
	changed_at TIMESTAMPTZ,
	operation LOG_OPERATION 
);

/****************************************/

DROP TABLE IF EXISTS session CASCADE;
CREATE TABLE session (
	id SERIAL PRIMARY KEY,
	adm_user_id INT REFERENCES adm_user(id),
	ip_address INET,
	created_at TIMESTAMPTZ,
	last_access_at TIMESTAMPTZ
);

/****************************************/

DROP TABLE IF EXISTS adm_user CASCADE;
CREATE TABLE adm_user (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	email VARCHAR NOT NULL DEFAULT '',
	password VARCHAR NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX ON adm_user(email);

/****************************************/

DROP TABLE IF EXISTS adm_user_log CASCADE;
CREATE TABLE adm_user_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	email VARCHAR
);

/****************************************/

DROP TABLE IF EXISTS client CASCADE;
CREATE TABLE client (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	birthday TIMESTAMPTZ
);

/****************************************/

DROP TABLE IF EXISTS client_log CASCADE;
CREATE TABLE client_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	birthday TIMESTAMPTZ
);

/****************************************/

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO druns;
GRANT SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO druns;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO druns;