CREATE TABLE IF NOT EXISTS archived_todo_tags(
    todoId int not null references archived_todos(id) on DELETE CASCADE,
    tagId int not null references tags(id) on DELETE CASCADE
)
