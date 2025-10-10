CREATE OR REPLACE FUNCTION group_grades_ordered(selected_group UUID) RETURNS TABLE(group_grade double precision, user_id uuid)
LANGUAGE PLPGSQL
AS $$
BEGIN 
    RETURN QUERY SELECT AVG(CAST(grd.grade AS FLOAT)), s.user_id FROM groups as grp 
    INNER JOIN students as s ON s.group_id = selected_group
    INNER JOIN grades as grd ON grd.student_id=s.user_id GROUP BY s.user_id 
    ORDER BY AVG(grd.grade);
END
$$;

CREATE OR REPLACE FUNCTION grades_from(avg_grade float) RETURNS TABLE (average_grade float, student_id uuid)
LANGUAGE plpgsql
as $$
BEGIN
    RETURN QUERY SELECT AVG(CAST((g.grade) AS FLOAT)), g.student_id
    FROM grades as g INNER JOIN users as u on u.id=g.student_id
    GROUP BY g.student_id HAVING AVG(CAST((g.grade) AS FLOAT)) > avg_grade ORDER BY AVG(CAST((g.grade) AS FLOAT));
END;
$$;