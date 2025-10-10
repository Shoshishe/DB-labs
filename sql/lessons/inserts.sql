CREATE OR REPLACE PROCEDURE lessons_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
   names TEXT[] := ARRAY['ABC', 'SCI', 'DBS','VOP','BOAP', 'MTWB', 'BCYSEC', 'OOP'];
BEGIN
   
    FOR i IN 1..50 LOOP
      INSERT INTO lessons_prototype (name, lesson_type) 
      SELECT names[1 + trunc(random() * 8)::int], (SELECT id FROM lessons_types ORDER BY random() LIMIT 1) ON CONFLICT DO NOTHING;
    END LOOP;

    FOR i in 1..5000 LOOP
      INSERT INTO prototype_to_groups (group_id, prototype_id) SELECT (SELECT id FROM groups ORDER BY RANDOM() LIMIT 1), (SELECT id FROM lessons_prototype ORDER BY RANDOM() LIMIT 1);
    END LOOP;

    FOR i IN 1..100000 LOOP
      INSERT INTO lessons (prototype_id, lesson_start, lesson_end) 
      SELECT (SELECT id FROM lessons_prototype ORDER BY random() LIMIT 1), (NOW() + (random() * INTERVAL '1 week')), 
      (NOW() + (random() * INTERVAL '1 week' * 2))::timestamptz;
    END LOOP;

  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
