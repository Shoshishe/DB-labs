CREATE OR REPLACE FUNCTION group_skipped_ordered() RETURNS TABLE (skipped_hours bigint, group_name text)
LANGUAGE plpgsql
AS $$
BEGIN 
    RETURN QUERY 
    SELECT COUNT(h.skipped_hours), g.name FROM skipped_hours as h 
        INNER JOIN groups as g ON g.id=h.group_id
        GROUP BY h.group_id, g.name 
        ORDER BY COUNT(h.skipped_hours);
END
$$;

CREATE OR REPLACE FUNCTION group_lesson_names(selected_group UUID) RETURNS TABLE (name TEXT)
LANGUAGE plpgsql 
AS $$
BEGIN
    RETURN QUERY SELECT DISTINCT p.name from lessons_prototype as p
    INNER JOIN prototype_to_groups as p_g ON p_g.group_id=selected_group
    INNER JOIN groups as g on g.id=selected_group;
END
$$;
