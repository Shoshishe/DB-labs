CREATE OR REPLACE FUNCTION next_lessons(selected_group_id UUID, curr_time timestamptz) RETURNS TABLE (lesson_id uuid, start_time timestamptz, end_time timestamptz, lesson_name text, lesson_type text)
LANGUAGE plpgsql
AS $$
BEGIN 
    RETURN QUERY SELECT l.id, l.lesson_start, l.lesson_end, lp.name, CASE lp.lesson_type WHEN 1 THEN 'SM' WHEN 2 THEN 'LR' WHEN 3 THEN 'LC' END as lt FROM lessons AS l
    INNER JOIN lessons_prototype AS lp ON lp.id=l.prototype_id
    INNER JOIN prototype_to_groups AS pg ON pg.prototype_id=lp.id
    INNER JOIN groups as g ON g.id=pg.group_id WHERE g.id=selected_group_id AND l.lesson_start > curr_time ORDER BY l.lesson_start;
END;
$$;