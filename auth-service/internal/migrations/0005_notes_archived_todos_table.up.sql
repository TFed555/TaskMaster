CREATE TABLE IF NOT EXISTS archived_todos(
    id SERIAL PRIMARY key,
    title varchar(255) not null,
    userId int not null unique references auth.users(id) on DELETE CASCADE,
    priority varchar(64) not null,
    category varchar(255) not null,
    description text not null,
    createdAt varchar(255) not null,
    completedAt varchar(255)
)