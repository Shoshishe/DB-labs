CREATE OR REPLACE PROCEDURE users_setup()
LANGUAGE plpgsql 
AS $$
DECLARE
   surnames TEXT[] :=ARRAY['Vladymtsev','Krupenich','Petrochenko', 'Mikhailov', 'Glamazdin', 'Tondel'];
   names TEXT[] := ARRAY['Egor', 'Andrei', 'Dmitriy','Vadim', 'Constantine', 'Mikhail', 'Aliaksei', 'Nikita'];
   patronymics TEXT[] := ARRAY['Pavlovich', 'Andreevich', 'Aliakseevich', 'Vadimovich', 'Nikitov', 'Nazarovich'];
BEGIN
    INSERT into users (surname, name, patronymic, email, password, university_id) 
    SELECT surnames[1 + trunc(random() * 6)::int],
    names[1 + trunc(random() * 8)::int],
    patronymics[1+trunc(random() * 6)::int],
    substr(md5(random()::text),1,5) || '@' || substr(md5(random()::text),1,5) || '.com',
    substr(md5(random()::text),1,200),
    (SELECT id from universities ORDER by random() LIMIT 1)  FROM generate_series(1,5000000);
    INSERT INTO user_roles (user_id, role_id) (SELECT id, (4 + trunc(random() * 2)::int) from users);
  EXCEPTION   
    WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END
$$;
