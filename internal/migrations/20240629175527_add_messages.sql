-- 20240629175527 - add_messages migration
ALTER TABLE conversation RENAME TO conversations;
ALTER TABLE conversations ADD COLUMN model TEXT DEFAULT 'llama3';
ALTER TABLE conversations ADD COLUMN context TEXT DEFAULT '';
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,

    message TEXT NOT NULL,
    kind TEXT NOT NULL DEFAULT 'user',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
