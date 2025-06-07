-- Step 1: Create new table with required columns and constraints
CREATE TABLE tickets_new (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    priority TEXT NOT NULL,
    status TEXT NOT NULL,
    creator_id TEXT NOT NULL DEFAULT 'default-admin-id',
    assignee_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Step 2: Copy data from old table to new table
INSERT INTO tickets_new (
    id, title, description, priority, status,
    creator_id, assignee_id, created_at, updated_at
)
SELECT
    id, title, description, priority, status,
    'default-admin-id', NULL, created_at, updated_at
FROM tickets;

-- Step 3: Drop old table
DROP TABLE tickets;

-- Step 4: Rename new table
ALTER TABLE tickets_new RENAME TO tickets;

