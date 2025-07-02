CREATE TABLE IF NOT EXISTS todos(
    id SERIAL PRIMARY key,
    title varchar(255) not null,
    userId int not null references auth.users(id) on DELETE CASCADE,
    priority varchar(64) not null,
    category varchar(255) not null,
    description text not null,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completedAt TIMESTAMP
)