CREATE OR REPLACE PROCEDURE omissions_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
BEGIN
  FOR i in 1..5000 LOOP
      INSERT INTO omissions (type_id, student_id, group_id, info, start_time, end_time)
      SELECT (SELECT id FROM omission_types ORDER BY RANDOM() LIMIT 1), u.user_id, u.group_id, substr(md5(random()::text),1,20), (NOW() + (random() * INTERVAL '1 week')), 
      (NOW() + (random() * INTERVAL '1 week' * 2))::timestamptz
      FROM (SELECT user_id,group_id FROM students ORDER BY RANDOM() LIMIT 1) as u;
  END LOOP;    
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;