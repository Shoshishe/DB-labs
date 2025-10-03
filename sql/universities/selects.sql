CREATE OR REPLACE PROCEDURE universities_select()
LANGUAGE SQL
AS $$
BEGIN
    SELECT id, full_name, shorthand FROM universities;
END
$$;

CREATE OR REPLACE PROCEDURE universitires_id_select(uni_id uuid)
LANGUAGE SQL
AS $$
BEGIN
    SELECT full_name, shorthand FROM universities WHERE id=uni_id;
END
$$:

CREATE OR REPLACE PROCEDURE universities_name_fuzzy(name text)
LANGUAGE SQL
AS $$
BEGIN
    SELECT id, full_name, shorthand FROM universities WHERE daitch_mokotoff(full_name) && daitch_mokotoff(name)
END
$$;