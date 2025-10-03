CREATE OR REPLACE PROCEDURE lessons_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
   names TEXT[] := ARRAY['ABC', 'SCI', 'DBS','VOP','BOAP', 'MTWB', 'BCYSEC', 'OOP'];
BEGIN
    INSERT INTO lessons_prototype (name, lesson_type) 
    SELECT names[1 + trunc(random() * 8)::int],
    (SELECT id FROM lessons_types ORDER BY random() LIMIT 1)
    FROM generate_series(1,50);
    
    INSERT INTO lessons (prototype_id, lesson_start, lesson_end) 
    SELECT (SELECT id FROM lessons_prototype ORDER BY RANDOM() LIMIT 1), (NOW() + (random() * INTERVAL '1 week')), 
    (NOW() + (random() * INTERVAL '1 week' * 2))::timestamptz FROM generate_series(1,10000);
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
