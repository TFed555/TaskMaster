-- SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.todos
-- 		WHERE userid = 1 AND createdAt >= '2025-06-04T15:30:45Z'::TIMESTAMPTZ LIMIT 1 OFFSET 0

-- INSERT INTO notes.todos (title, priority, category, description, createdAt, completedAt, userId) 
-- VALUES ('dd', 'low', 'x', 's', '2025-06-07', '2025-06-07', 1)

-- Select * from notes.todos

-- select * from auth.refresh_tokens
-- SELECT * FROM auth.refresh_tokens
--         WHERE token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIwNDQ4NDIsInN1YiI6MX0.YEaXpkpn7usHHCezBA4hC2MEyvrfxqfo_L9shGAnVOY'
--         LIMIT 1

INSERT INTO notes.todos (userid, title, priority, description, category, createdat) VALUES (1, 'a', 'v', 'v', 'c', NOW())