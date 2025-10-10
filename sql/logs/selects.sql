CREATE OR REPLACE FUNCTION logs_from_to(start_point timestamptz, end_point timestamptz) RETURNS TABLE(action jsonb, serialized_at timestamptz, user_id uuid, email text, full_name text)
LANGUAGE plpgsql
as $$
BEGIN
    RETURN QUERY SELECT 
        l.action, l.serialized_at, l.user_id, u.email, (u.surname ||  ' ' || u.name || ' ' || u.patronymic) as full_name 
    FROM logs as l 
    LEFT JOIN users as u ON u.id=l.user_id WHERE l.serialized_at >= start_point and l.serialized_at <= end_point
    ORDER BY l.serialized_at;
END;
$$;