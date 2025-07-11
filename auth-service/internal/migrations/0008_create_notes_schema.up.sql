CREATE SCHEMA IF NOT EXISTS notes;

ALTER TABLE public.todos SET SCHEMA notes;
ALTER TABLE public.archived_todos SET SCHEMA notes;
ALTER TABLE public.tags SET SCHEMA notes;
ALTER TABLE public.todo_tags SET SCHEMA notes;