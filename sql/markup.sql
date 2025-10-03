CREATE EXTENSION fuzzystrmatch;

CREATE TABLE IF NOT EXISTS lessons_types (
	id smallint NOT NULL GENERATED ALWAYS AS IDENTITY,
	name text NOT NULL,
	CONSTRAINT pk_lessons_types PRIMARY KEY (id),
	CONSTRAINT lessons_types_uq UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS roles (
	id smallint NOT NULL GENERATED ALWAYS AS IDENTITY,
	name text NOT NULL,
	CONSTRAINT roles_name_check CHECK(LENGTH(NAME) < 20),
	CONSTRAINT pk_roles_0 PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS universities (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	full_name text NOT NULL,
	CONSTRAINT universities_full_name_ck CHECK(LENGTH(full_name) < 50),
	shorthand text NOT NULL,
	CONSTRAINT universities_shorthand_ck CHECK(LENGTH(shorthand) < 10),
	CONSTRAINT pk_universities PRIMARY KEY (id),
	CONSTRAINT universities_name_uq UNIQUE (full_name)
);

CREATE INDEX universities_names_idx ON universities USING gin (daitch_mokotoff(full_name)) WITH (fastupdate = off);

CREATE TABLE IF NOT EXISTS users (
	password text NOT NULL,
	CONSTRAINT users_password_ck CHECK(LENGTH(password) < 256),
	university_id uuid NOT NULL,
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	surname text NOT NULL,
	CONSTRAINT users_surname_ck CHECK(LENGTH(surname) < 50),
	name text NOT NULL,
	CONSTRAINT users_name_ck CHECK(LENGTH(name) < 50),
	email text NOT NULL,
	CONSTRAINT users_email_ck CHECK(LENGTH(email) < 50),
	patronymic text NOT NULL,
	CONSTRAINT users_patronymic_ck CHECK(LENGTH(patronymic) < 50),
	CONSTRAINT unq_users UNIQUE (password, university_id),
	CONSTRAINT pk_users PRIMARY KEY (id),
	CONSTRAINT fk_users_universities FOREIGN KEY (university_id) REFERENCES universities(id) ON DELETE RESTRICT ON UPDATE RESTRICT
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
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	name text NOT NULL,
	CONSTRAINT groups_name_check CHECK (LENGTH(name) < 8),
	faculty_id uuid NOT NULL,
	CONSTRAINT fk_groups_faculties FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT pk_groups PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS lessons_prototype (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	name text NOT NULL,
	lesson_type smallint NOT NULL,
	CONSTRAINT pk_lessons_prototype PRIMARY KEY (id),
	CONSTRAINT fk_lessons_prototype_types FOREIGN KEY (lesson_type) REFERENCES lessons_types(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS logs (
	user_id uuid DEFAULT NULL,
	action jsonb NOT NULL,
	CONSTRAINT cns_logs CHECK(LENGTH(action::jsonb::text) < 1500),
	serialized_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT fk_logs_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
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
	CONSTRAINT unq_students_user_id UNIQUE (user_id),
	CONSTRAINT fk_students_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT fk_students_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
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
	skipped_hours bigint NOT NULL CHECK (skipped_hours >=0 ),
	skipped_lesson_id uuid NOT NULL,
	CONSTRAINT fk_skipped_hours_students FOREIGN KEY (student_id) REFERENCES students(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_skipped_hours_lessons FOREIGN KEY (skipped_lesson_id) REFERENCES lessons(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS skipped_hours_user_id_idx ON skipped_hours (student_id);

CREATE TABLE IF NOT EXISTS grades (
	student_id uuid NOT NULL,
	grade smallint NOT NULL CHECK (grade>=0),
	lesson_id uuid NOT NULL,
	CONSTRAINT fk_grades_lesson_id FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS grades_user_id_idx ON grades (student_id);