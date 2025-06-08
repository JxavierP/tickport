PRAGMA foreign_keys = OFF;

-- Step 1: Rename old table
ALTER TABLE users RENAME TO users_old;

-- Step 2: Recreate users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Copy data
INSERT INTO users (id, username, email, role, password_hash, created_at, updated_at)
SELECT id, username, email, role, password_hash, created_at, updated_at FROM users_old;

-- Step 4: Drop old table
DROP TABLE users_old;

PRAGMA foreign_keys = ON;

