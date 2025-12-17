CREATE OR REPLACE FUNCTION student_id_select(id UUID) RETURNS TABLE (university_id uuid, name text, surname text, patronymic text, email text, group_name text)
LANGUAGE PLPGSQL
as $$
BEGIN
   RETURN QUERY SELECT u.id, u.university_id, u.name, u.surname, u.patronymic, u.email, g.name
                FROM users AS u INNER JOIN students AS s ON s.user_id = u.id 
                INNER JOIN groups AS g on s.group_id=g.id WHERE u.id=id;
END;
$$;

-- CREATE OR REPLACE PROCEDURE student_skipped_hours(id UUID) 
-- LANGUAGE SQL
-- as $$
-- BEGIN
--     SELECT SUM(skipped_hours) FROM skipped_hours WHERE student_id=id;
-- END
-- $$;

CREATE OR REPLACE FUNCTION student_skipped_lessons(user_id UUID) RETURNS TABLE(lesson_start timestamptz, lesson_end timestamptz, name text, lesson_type smallint, proto_name text, skipped_hours bigint)
LANGUAGE PLPGSQL
as $$
BEGIN
    RETURN QUERY SELECT l.lesson_start, l.lesson_end, p.name, p.lesson_type, s.skipped_hours FROM skipped_hours as s 
    INNER JOIN lessons as l ON s.skipped_lesson_id=l.id 
    INNER JOIN lessons_prototypes as p ON p.id=l.prototype_id
    WHERE s.student_id=user_id;
END
$$;

CREATE OR REPLACE FUNCTION students_by_group(group_id UUID) RETURNS TABLE(id uuid, university_id uuid, name text, surname text, patronymic text, email text, group_name text)
LANGUAGE plpgsql
as $$
BEGIN
    RETURN QUERY SELECT u.id, u.university_id, u.name, u.surname, u.patronymic, u.email, g.name FROM users AS u 
    INNER JOIN students AS s ON s.user_id = u.id 
    INNER JOIN groups AS g on s.group_id=g.id WHERE u.id=id;
END;
$$;
