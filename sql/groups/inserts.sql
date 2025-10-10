CREATE OR REPLACE PROCEDURE groups_setup()
LANGUAGE plpgsql 
AS $$
BEGIN
    FOR i IN 1..500 LOOP
       INSERT into groups (name, faculty_id) 
       SELECT substr(md5(random()::text),1,4),
       (SELECT id from faculties ORDER by random() LIMIT 1);
    END LOOP;
END
$$;
