CREATE OR REPLACE PROCEDURE next_lessons(group_id UUID, curr_time timestamptz)
LANGUAGE plpgsql
AS $$
BEGIN 
    SELECT l.id, l.lesson_start, l.lesson_end, lp.name, lp.lesson_type FROM lessons AS l
    INNER JOIN lessons_prototype AS lp ON lp.id=l.prototype_id
    INNER JOIN prototype_to_groups AS pg ON pg.prototype_id=lp.INDEX
    INNER JOIN groups as g ON g.id=pg.id WHERE g.id=group_id AND l.lesson > curr_time
END;
$$;