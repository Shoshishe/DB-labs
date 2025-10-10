CREATE OR REPLACE PROCEDURE faculties_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
 names TEXT[] := ARRAY['FCSaN','FRE', 'FoGaAP', 'FoAM', 'FoCP'];
BEGIN
    FOR i IN 1..100 LOOP
        INSERT into faculties (name, university_id) 
        SELECT names[MOD(i, 5)+1],
        (SELECT id FROM universities order by random() limit 1) ON CONFLICT DO NOTHING;
    END LOOP;
END;
$$;
