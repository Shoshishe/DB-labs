CREATE OR REPLACE PROCEDURE groups_setup()
LANGUAGE plpgsql 
AS $$
BEGIN
    INSERT into groups (name, faculty_id) 
    SELECT (substr(md5(random()::text),1,4), 1),
    (SELECT id from faculties ORDER by random() LIMIT 1)  FROM generate_series(1,50000);
END
$$;
