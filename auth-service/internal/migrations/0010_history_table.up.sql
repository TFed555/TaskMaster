CREATE TABLE IF NOT EXISTS todos_history (
    id SERIAL PRIMARY KEY,
    todoid int references notes.todos(id) NULL,
    archived_todoid int references notes.archived_todos(id) NULL,
    userId int not null references auth.users(id) on delete cascade,
    action text not null,
    old_value  text,
    new_value text,
    changedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)