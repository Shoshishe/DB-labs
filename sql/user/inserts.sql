CREATE OR REPLACE PROCEDURE users_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
   surnames TEXT[] :=ARRAY['Vladymtsev','Krupenich','Petrochenko', 'Mikhailov', 'Glamazdin', 'Tondel'];
   names TEXT[] := ARRAY['Egor', 'Andrei', 'Dmitriy','Vadim', 'Constantine', 'Mikhail', 'Aliaksei', 'Nikita'];
   patronymics TEXT[] := ARRAY['Pavlovich', 'Andreevich', 'Aliakseevich', 'Vadimovich', 'Nikitov', 'Nazarovich'];
   row RECORD;
BEGIN
   FOR i in 1..500 LOOP
      INSERT into users (surname, name, patronymic, email, password) 
      SELECT surnames[1 + trunc(random() * 6)::int],
      names[1 + trunc(random() * 8)::int],
      patronymics[1+trunc(random() * 6)::int],
      substr(md5(random()::text),1,5) || '@' || substr(md5(random()::text),1,5) || '.com',
      substr(md5(random()::text),1,200);
      INSERT INTO user_roles (user_id, role_id, university_id) (SELECT id, ROUND(1 + (random() * 1))::int, (SELECT id from universities ORDER by random() LIMIT 1) from users);
   END LOOP;    
   
   FOR row IN SELECT * from user_roles WHERE role_id=1 LOOP
      INSERT INTO students (user_id, group_id) SELECT row.user_id, (SELECT id from groups order by random() LIMIT 1) ON CONFLICT DO NOTHING;
   END LOOP;


   FOR row IN SELECT * from user_roles WHERE role_id=2 LOOP
      INSERT INTO teachers (user_id, photo_name) SELECT row.user_id, substr(md5(random()::text),1,200) ON CONFLICT DO NOTHING;
   END LOOP;
   
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
