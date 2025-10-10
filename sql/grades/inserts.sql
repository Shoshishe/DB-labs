CREATE OR REPLACE PROCEDURE grades_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
BEGIN
  FOR i in 1..5000 LOOP
      INSERT into grades (grade, student_id, group_id, lesson_id)
      SELECT ((1 + random() * 9)::int2), u.user_id, u.group_id, (SELECT id FROM lessons ORDER BY RANDOM() LIMIT 1)  
      FROM (SELECT user_id,group_id FROM students ORDER BY RANDOM() LIMIT 1) as u;
  END LOOP;    
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
