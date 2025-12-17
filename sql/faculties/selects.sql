CREATE OR REPLACE FUNCTION faculties_select() RETURNS TABLE (id uuid, name text, university_id uuid) AS $$
BEGIN
    RETURN QUERY SELECT f.id,f.name,f.university_id FROM faculties as f;
END;
$$ LANGUAGE PLPGSQL;