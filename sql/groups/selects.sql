CREATE OR REPLACE PROCEDURE group_skipped_ordered(group_id UUID)
LANGUAGE SQL
AS $$
BEGIN 
    SELECT COUNT(h.skipped_hours), g.name, g.id FROM groups as g INNER JOIN students as s ON s.group_id = group_id
    INNER JOIN skipped_hours as h ON h.student_id=s.user_id GROUP BY g.name ORDER BY COUNT(h.skipped_hours);
END
$$;

CREATE OR REPLACE PROCEDURE group_lesson_names(group_id UUID) 
LANGUAGE SQL 
AS $$
BEGIN
    SELECT DISTINCT p.name from lessons_prototypes INNER JOIN prototype_to_groups as p_g ON p_g.group_id=group_id INNER JOIN groups as g on g.id=group_id;
END
$$;

