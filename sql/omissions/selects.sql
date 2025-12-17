CREATE OR REPLACE FUNCTION get_student_omissions(student_id uuid) RETURNS TABLE (ugroup_id UUID, info TEXT, start_time timestamptz, end_time timestamptz, type_name text)
LANGUAGE PLPGSQL
AS $$
BEGIN
 RETURN QUERY SELECT o.student_id, o.group_id, o.info, o.start_time, o.end_time, t.name FROM omissions as o INNER JOIN omission_types as t ON o.type_id=t.id WHERE o.student_it=student_id;
END
$$;


