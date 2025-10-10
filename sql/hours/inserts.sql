CREATE OR REPLACE PROCEDURE hours_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
BEGIN
  FOR i in 1..5000 LOOP
      INSERT into skipped_hours (skipped_hours, student_id, group_id, skipped_lesson_id)
      SELECT ((1 + random() * 3)::int2), u.user_id, u.group_id, (SELECT id FROM lessons ORDER BY RANDOM() LIMIT 1)  
      FROM (SELECT user_id,group_id FROM students ORDER BY RANDOM() LIMIT 1) as u;
  END LOOP;    
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
