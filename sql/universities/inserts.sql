CREATE OR REPLACE PROCEDURE universities_setup() 
LANGUAGE plpgsql
AS $$
BEGIN 
      INSERT INTO universities (full_name, shorthand) 
      SELECT substr((md5(random()::text)),1,5), 
             substr(md5(random()::text),1,3)
      FROM generate_series(1,20);       
END
$$;