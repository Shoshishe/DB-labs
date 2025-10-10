CREATE OR REPLACE PROCEDURE teacher_select(user_id uuid)
LANGUAGE plpgsql
AS $$
BEGIN
  SELECT t.photo_name, u.surname, u.name, u.patronymic, u.email, u.university_id FROM teachers as t 
  INNER join users as u on u.id=t.user_id WHERE t.user_id=user_id 
END
$$;