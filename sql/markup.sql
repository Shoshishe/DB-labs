CREATE SCHEMA IF NOT EXISTS schema_iis;

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;

CREATE TABLE IF NOT EXISTS lessons_types (
	id smallint NOT NULL GENERATED ALWAYS AS IDENTITY,
	name text NOT NULL,
	
	CONSTRAINT pk_lessons_types PRIMARY KEY (id),
	CONSTRAINT lessons_types_uq UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS universities (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	full_name text NOT NULL,
	shorthand text NOT NULL,
	
	CONSTRAINT universities_full_name_ck CHECK(LENGTH(full_name) < 100),
	CONSTRAINT universities_shorthand_ck CHECK(LENGTH(shorthand) < 10),
	CONSTRAINT pk_universities PRIMARY KEY (id),
	CONSTRAINT universities_name_uq UNIQUE (full_name)
);

CREATE INDEX IF NOT EXISTS universities_names_idx ON universities USING gin (daitch_mokotoff(full_name)) WITH (fastupdate = off);

CREATE TABLE IF NOT EXISTS users (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	password text NOT NULL,
	surname text NOT NULL,
	name text NOT NULL,
	email text NOT NULL,
	patronymic text NOT NULL,

	CONSTRAINT users_name_ck CHECK(LENGTH(name) < 50),
	CONSTRAINT users_surname_ck CHECK(LENGTH(surname) < 50),
	CONSTRAINT users_password_ck CHECK(LENGTH(password) < 256),
	CONSTRAINT users_email_ck CHECK(LENGTH(email) < 50),
	CONSTRAINT users_patronymic_ck CHECK(LENGTH(patronymic) < 50),
	CONSTRAINT pk_users PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS roles (
	id smallint NOT NULL GENERATED ALWAYS AS IDENTITY,
	name text NOT NULL,

	CONSTRAINT roles_name_check CHECK(LENGTH(NAME) < 20),
	CONSTRAINT pk_roles_0 PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS faculties (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	name text NOT NULL,
	CONSTRAINT faculties_name_ck CHECK (LENGTH(name) < 30),
	university_id uuid NOT NULL,
	
	CONSTRAINT pk_faculties PRIMARY KEY (id),
	CONSTRAINT faculties_name_unq UNIQUE (name, university_id),
	CONSTRAINT fk_faculties_universities FOREIGN KEY (university_id) REFERENCES universities(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS groups (
	name text NOT NULL,
	faculty_id uuid NOT NULL,

	CONSTRAINT groups_name_check CHECK (LENGTH(name) < 8),
	CONSTRAINT fk_groups_faculties FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT pk_groups PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS lessons_prototype (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	name text NOT NULL,
	lesson_type smallint NOT NULL,

	CONSTRAINT uq_lessons_prototype UNIQUE (name, lesson_type),
	CONSTRAINT pk_lessons_prototype PRIMARY KEY (id),
	CONSTRAINT fk_lessons_prototype_types FOREIGN KEY (lesson_type) REFERENCES lessons_types(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS audit_log_inserts (
	table_name TEXT NOT NULL,
	new_values jsonb NOT NULL,
	serialized_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_log_deletes (
	deleted_id uuid,
	table_name TEXT NOT NULL,
	old_values jsonb NOT NULL,
	serialized_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_log_updates (
	updated_id uuid NOT NULL,
	table_name TEXT NOT NULL, 
	old_values jsonb NOT NULL,
	new_values jsonb NOT NULL,
	serialized_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS prototype_to_groups (
	group_id uuid NOT NULL,
	prototype_id uuid NOT NULL,

	CONSTRAINT fk_prototype_to_groups_groups FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT fk_prototype_to_groups_lessons_prototype FOREIGN KEY (prototype_id) REFERENCES lessons_prototype(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS students (
	user_id uuid NOT NULL,
	group_id uuid NOT NULL,

	CONSTRAINT unq_students UNIQUE (user_id, group_id),
	CONSTRAINT fk_students_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT fk_students_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS omission_types(
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	name text not null
);

CREATE TABLE IF NOT EXISTS omissions(
	student_id uuid not null,
	group_id uuid not null,
	info text not null,
	start_time timestamptz not null,
	end_time timestamptz not null,
	type_id uuid not null,

	CONSTRAINT fk_omissions_students FOREIGN KEY (student_id, group_id) REFERENCES students(user_id, group_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_omissions_types FOREIGN KEY (type_id) REFERENCES omission_types(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT ck_certificate_reason CHECK (LENGTH(info) < 500)
);

CREATE TABLE IF NOT EXISTS teachers (
	user_id uuid NOT NULL,
	photo_name text,

	CONSTRAINT teachers_photo_name_ck CHECK(LENGTH(photo_name) < 256),
	CONSTRAINT teachers_photo_check CHECK(LENGTH(photo_name) < 100),
	CONSTRAINT fk_teachers_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS teachers_user_id_uq ON teachers (user_id);

CREATE TABLE IF NOT EXISTS teachers_to_lessons (
	lesson_prototype_id uuid NOT NULL,
	teacher_id uuid NOT NULL,

	CONSTRAINT fk_teachers_to_lessons_teachers FOREIGN KEY (teacher_id) REFERENCES teachers(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_teachers_to_lessons_lessons_prototype FOREIGN KEY (lesson_prototype_id) REFERENCES lessons_prototype(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS user_roles (
	role_id smallint NOT NULL,
	user_id uuid NOT NULL,
    university_id uuid NOT NULL,
	
	CONSTRAINT fk_users_universities FOREIGN KEY (university_id) REFERENCES universities(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT fk_user_roles_roles_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT fk_user_roles_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS lessons (
	prototype_id uuid NOT NULL,
	lesson_start timestamptz NOT NULL,
	lesson_end timestamptz NOT NULL,
	id uuid NOT NULL DEFAULT gen_random_uuid(),

	CONSTRAINT pk_lessons PRIMARY KEY (id),
	CONSTRAINT fk_lessons_lessons_prototype FOREIGN KEY (prototype_id) REFERENCES lessons_prototype(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS skipped_hours (
	student_id uuid NOT NULL,
	group_id uuid NOT NULL,
	skipped_hours bigint NOT NULL CHECK (skipped_hours >=0 ),
	lesson_id uuid NOT NULL,
	is_legitimate BOOLEAN NOT NULL DEFAULT FALSE,

	CONSTRAINT fk_skipped_hours_students FOREIGN KEY (student_id, group_id) REFERENCES students(user_id, group_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_skipped_hours_lessons FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE OR REPLACE FUNCTION hours_omissions_trigger() RETURNS TRIGGER AS $$
BEGIN
	IF EXISTS(SELECT 1 FROM omissions AS o 
				INNER JOIN skipped_hours AS sh ON sh.student_id=o.student_id 
				INNER JOIN lessons as l ON l.id=sh.lesson_id 
				WHERE o.end_time >= l.start_time) THEN
			NEW.is_legitimate=TRUE;
	END IF;	
END;
$$ LANGUAGE PLPGSQL;
DROP TRIGGER IF EXISTS hours_omissions_tg ON skipped_hours;
CREATE TRIGGER hours_omissions_tg BEFORE INSERT OR UPDATE ON skipped_hours FOR EACH ROW EXECUTE PROCEDURE hours_omissions_trigger();

CREATE OR REPLACE FUNCTION omissions_trigger() RETURNS TRIGGER AS $$
BEGIN
	IF (TG_OP='INSERT') THEN	
		UPDATE skipped_hours AS sh SET is_legitimate=TRUE
		FROM lessons AS l WHERE sh.lesson_id=l.id 
		AND sh.student_id = NEW.student_id AND l.start_time <= new.end_time;
	ELSIF (TG_OP='DELETE') THEN
		UPDATE skipped_hours AS sh SET is_legitimate=FALSE 
		FROM lessons AS l WHERE sh.lesson_id=l.id 
		AND sh.student_id = OLD.student_id AND l.start_time <= OLD.end_time;	
	ELSIF (TG_OP='UPDATE') THEN
		UPDATE skipped_hours AS sh SET is_legitimate=FALSE 
		FROM lessons AS l WHERE sh.lesson_id=l.id 
		AND sh.student_id = OLD.student_id AND l.start_time <= OLD.end_time;	

		UPDATE skipped_hours AS sh SET is_legitimate=TRUE
		FROM lessons AS l WHERE sh.lesson_id=l.id 
		AND sh.student_id = NEW.student_id AND l.start_time <= new.end_time;
	END IF;
END;
$$ LANGUAGE PLPGSQL;
DROP TRIGGER IF EXISTS omission_hours_tg ON omissions;
CREATE TRIGGER omission_hours_tg AFTER INSERT OR UPDATE OR DELETE ON omissions FOR EACH ROW EXECUTE PROCEDURE omissions_trigger();


CREATE INDEX IF NOT EXISTS skipped_hours_user_id_idx ON skipped_hours (student_id);

CREATE TABLE IF NOT EXISTS grades (
	student_id uuid NOT NULL,
	group_id uuid NOT NULL,
	grade smallint NOT NULL CHECK (grade>=0),
	lesson_id uuid NOT NULL,

	CONSTRAINT fk_grades_student_id FOREIGN KEY (student_id, group_id) REFERENCES students(user_id, group_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_grades_lesson_id FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION audit_trigger() RETURNS TRIGGER AS $$
DECLARE
    new_data jsonb;
    old_data jsonb;
    key text;
    new_values jsonb;
    old_values jsonb;
    user_id uuid;
BEGIN
	new_values = '{}';
	old_values = '{}';

	IF (TG_OP='INSERT') THEN 
		new_data := to_jsonb(NEW);
		new_values := new_data;
		INSERT INTO audit_log_inserts (table_name, old_values, new_values) 
		VALUES (TG_TABLE_NAME, old_values, new_values);
		RETURN NEW;

	ELSIF (TG_OP='UPDATE') THEN
		new_data := to_jsonb(NEW);
		old_data := to_jsonb(OLD);

		FOR key in SELECT jsonb_object_keys(new_data) INTERSECT SELECT jsonb_object_keys(old_data)
		LOOP
			IF new_data ->> key != old_data ->> key THEN
				new_values := new_values || jsonb_build_object(key, new_data ->> key);
				old_values := old_values || jsonb_build_object(key, old_data ->> key);
			END IF;
		END LOOP;	

		INSERT INTO audit_log_updates (table_name, old_values, new_values, updated_id) 
		VALUES (TG_TABLE_NAME, old_values, new_values, NEW.id);
		RETURN NEW;

	ELSIF TG_OP = 'DELETE' THEN
			old_data := to_jsonb(OLD);
			old_values := old_data;
			FOR key IN SELECT jsonb_object_keys(old_data)
			LOOP
				old_values := old_values || jsonb_build_object(key, old_data ->> key);
			END LOOP;

			INSERT INTO audit_log_deletes (table_name, old_values) 
			VALUES (TG_TABLE_NAME, old_values);
			RETURN NEW;
    END IF;
END;
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE PROCEDURE DDL_APPLY_AUDIT_TRIGGER() AS $$
DECLARE
_sql TEXT;
BEGIN 
	FOR _sql IN SELECT CONCAT('DROP TRIGGER IF EXISTS tg_audit_', quote_ident(it.table_name), ' ON ', quote_ident(it.table_name), '; CREATE TRIGGER tg_audit_',quote_ident(it.table_name), ' AFTER INSERT OR UPDATE OR DELETE ON ', quote_ident(it.table_name), 
	' FOR EACH STATEMENT EXECUTE PROCEDURE audit_trigger();') FROM information_schema.tables as it WHERE it.table_schema not in ('pg_catalog', 'information_schema') AND it.table_schema NOT LIKE 'pg_toast%' AND EXISTS(SELECT 1 FROM information_schema.columns as ic WHERE ic.column_name='id' AND ic.table_name=it.table_name) LOOP
	EXECUTE TRIM('\"' FROM _sql);
	END LOOP;
END;
$$ LANGUAGE PLPGSQL;

CALL DDL_APPLY_AUDIT_TRIGGER();


