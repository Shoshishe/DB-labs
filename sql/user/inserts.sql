CREATE OR REPLACE PROCEDURE users_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
   surnames TEXT[] :=ARRAY['Vladymtsev','Krupenich','Petrochenko', 'Mikhailov', 'Glamazdin', 'Tondel'];
   names TEXT[] := ARRAY['Egor', 'Andrei', 'Dmitriy','Vadim', 'Constantine', 'Mikhail', 'Aliaksei', 'Nikita'];
   patronymics TEXT[] := ARRAY['Pavlovich', 'Andreevich', 'Aliakseevich', 'Vadimovich', 'Nikitov', 'Nazarovich'];
BEGIN
  FOR i in 1..500 LOOP
      INSERT into users (surname, name, patronymic, email, password) 
      SELECT surnames[1 + trunc(random() * 6)::int],
      names[1 + trunc(random() * 8)::int],
      patronymics[1+trunc(random() * 6)::int],
      substr(md5(random()::text),1,5) || '@' || substr(md5(random()::text),1,5) || '.com',
      substr(md5(random()::text),1,200);
      INSERT INTO user_roles (user_id, role_id, university_id) (SELECT id, (2 + trunc(random() * 1)::int), (SELECT id from universities ORDER by random() LIMIT 1) from users);
    END LOOP;    
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
