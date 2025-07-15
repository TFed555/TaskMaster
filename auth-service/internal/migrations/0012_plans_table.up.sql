CREATE TABLE IF NOT EXISTS plans_todo(
    ID SERIAL PRIMARY KEY,
    todoID int not null references notes.todos(id) on delete CASCADE,
    steps TEXT[] NOT NULL DEFAULT '{}'
)