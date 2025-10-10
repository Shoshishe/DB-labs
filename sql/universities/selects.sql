CREATE OR REPLACE FUNCTION universities_select() RETURNS TABLE (uni_id uuid, uni_name text, uni_shorthand text)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT id, full_name, shorthand FROM universities;
END
$$;

CREATE OR REPLACE FUNCTION universitires_id_select(uni_id uuid) RETURNS TABLE (uni_name text, uni_shorthand text)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT full_name, shorthand FROM universities WHERE id=uni_id;
END
$$;

CREATE OR REPLACE FUNCTION universities_name_fuzzy(name text) RETURNS TABLE (uni_id uuid, uni_name text, uni_shorthand text)
LANGUAGE plpgSQL
AS $$
BEGIN
    RETURN QUERY SELECT id, full_name, shorthand FROM universities WHERE daitch_mokotoff(full_name) && daitch_mokotoff(name);
END
$$;
