CREATE OR REPLACE FUNCTION grades_from(avg_grade float) RETURNS TABLE (average_grade float, student_id uuid)
LANGUAGE plpgsql
as $$
BEGIN
    RETURN QUERY SELECT AVG(CAST((g.grade) AS FLOAT)), g.student_id
    FROM grades as g INNER JOIN users as u on u.id=g.student_id
    GROUP BY g.student_id HAVING AVG(CAST((g.grade) AS FLOAT)) > avg_grade;
END;
$$;

