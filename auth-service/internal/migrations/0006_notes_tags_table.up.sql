CREATE TABLE IF NOT EXISTS tags(
    id SERIAL PRIMARY key,
    userId int not null references auth.users(id) on DELETE CASCADE,
    name text not null
)

