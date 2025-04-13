ALTER TABLE  IF EXISTS users
ADD COLUMN role_id INT REFERENCES roles(id) ;

update users set role_id = (
    select id from roles where name = 'user'
) where role_id is null;


ALTER TABLE IF EXISTS users
ALTER COLUMN role_id SET NOT NULL;
