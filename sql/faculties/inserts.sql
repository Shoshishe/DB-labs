CREATE OR REPLACE PROCEDURE faculties_setup()
LANGUAGE plpgsql 
AS $$
BEGIN
    INSERT into faculties (name, university_id) 
    SELECT substr(md5(random()::text),1,10),
    (SELECT id FROM universities order by random() limit 1)  FROM generate_series(1,50000);
END
$$;
