CREATE OR REPLACE FUNCTION teacher_select(user_id uuid) RETURNS TABLE (photo text, user_surname text, user_name text, user_patronymic text, user_email text, user_uni uuid)
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY SELECT t.photo_name, u.surname, u.name, u.patronymic, u.email, u.university_id FROM teachers as t 
  INNER join users as u on u.id=t.user_id WHERE t.user_id=user_id;
END
$$;


