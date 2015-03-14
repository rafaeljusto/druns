REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM druns;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM druns;
REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public FROM druns;

DROP USER IF EXISTS druns;
CREATE USER druns WITH PASSWORD 'abc123';

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

DROP TYPE IF EXISTS LOG_OPERATION CASCADE;
CREATE TYPE LOG_OPERATION AS ENUM ('Create', 'Update', 'Delete');

DROP TABLE IF EXISTS log CASCADE;
CREATE TABLE log (
	id SERIAL PRIMARY KEY,
	agent INT REFERENCES adm_user(id),
	ip_address INET,
	changed_at TIMESTAMPTZ,
	operation LOG_OPERATION 
);

/****************************************/

DROP TABLE IF EXISTS adm_user_log CASCADE;
CREATE TABLE adm_user_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	email VARCHAR
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

DROP TABLE IF EXISTS client CASCADE;
CREATE TABLE client (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	birthday DATE
);

/****************************************/

DROP TABLE IF EXISTS client_log CASCADE;
CREATE TABLE client_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	birthday DATE
);

/****************************************/

DROP TABLE IF EXISTS place CASCADE;
CREATE TABLE place (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	address VARCHAR
);

/****************************************/

DROP TABLE IF EXISTS place_log CASCADE;
CREATE TABLE place_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	address VARCHAR
);

/****************************************/

DROP TYPE IF EXISTS WEEKDAY CASCADE;
CREATE TYPE WEEKDAY AS ENUM ('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');

DROP TYPE IF EXISTS CLIENT_GROUP_TYPE CASCADE;
CREATE TYPE CLIENT_GROUP_TYPE AS ENUM ('Weekley', 'Once');

DROP TABLE IF EXISTS client_group CASCADE;
CREATE TABLE client_group (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	place_id INT REFERENCES place(id),
	weekday WEEKDAY,
	time TIME,
	duration INT,
	type CLIENT_GROUP_TYPE,
	capacity INT
);

/****************************************/

DROP TABLE IF EXISTS client_group_log CASCADE;
CREATE TABLE client_group_log (
	log_id INT REFERENCES log(id),
	id INT,
	name VARCHAR,
	place_id INT,
	weekday WEEKDAY,
	time TIME,
	duration INT,
	type CLIENT_GROUP_TYPE,
	capacity INT
);

/****************************************/

DROP TYPE IF EXISTS ENROLLMENT_TYPE CASCADE;
CREATE TYPE ENROLLMENT_TYPE AS ENUM ('Reservation', 'Regular', 'Replacement');

DROP TABLE IF EXISTS enrollment CASCADE;
CREATE TABLE enrollment (
	id SERIAL PRIMARY KEY,
	client_id INT REFERENCES client(id),
	client_group_id INT REFERENCES client_group(id),
	type ENROLLMENT_TYPE
);

CREATE UNIQUE INDEX ON enrollment(client_id, client_group_id);
CREATE INDEX ON enrollment(client_id);
CREATE INDEX ON enrollment(client_group_id);

/****************************************/

DROP TABLE IF EXISTS enrollment_log CASCADE;
CREATE TABLE enrollment_log (
	log_id INT REFERENCES log(id),
	id INT,
	client_id INT,
	client_group_id INT,
	type ENROLLMENT_TYPE
);

/****************************************/

DROP TABLE IF EXISTS class CASCADE;
CREATE TABLE class (
	id SERIAL PRIMARY KEY,
	client_group_id INT REFERENCES client_group(id),
	class_date TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON class(client_group_id, class_date);

/****************************************/

DROP TABLE IF EXISTS class_log CASCADE;
CREATE TABLE class_log (
	log_id INT REFERENCES log(id),
	id INT,
	client_group_id INT,
	class_date TIMESTAMPTZ
);

/****************************************/

DROP TABLE IF EXISTS student CASCADE;
CREATE TABLE student (
	id SERIAL PRIMARY KEY,
	class_id INT REFERENCES class(id),
	enrollment_id INT REFERENCES enrollment(id),
	attended BOOLEAN
);

/****************************************/

DROP TABLE IF EXISTS student_log CASCADE;
CREATE TABLE student_log (
	log_id INT REFERENCES log(id),
	id INT,
	class_id INT,
	enrollment_id INT,
	attended BOOLEAN
);

/****************************************/

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO druns;
GRANT SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO druns;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO druns;