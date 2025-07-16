ALTER TABLE todos_history
DROP CONSTRAINT todos_history_todoid_fkey,
ADD CONSTRAINT todos_history_todoid_fkey 
FOREIGN KEY (todoid) REFERENCES notes.todos(id) ON DELETE SET NULL;

ALTER TABLE todos_history
DROP CONSTRAINT todos_history_archived_todoid_fkey,
ADD CONSTRAINT todos_history_archived_todoid_fkey 
FOREIGN KEY (archived_todoid) REFERENCES notes.archived_todos(id) ON DELETE SET NULL;