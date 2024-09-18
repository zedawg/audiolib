CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    status TEXT DEFAULT 'queued',
    params JSONB DEFAULT '[]',
    result JSONB DEFAULT '[]',
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER limit_tasks_after_insert
AFTER INSERT ON tasks
BEGIN
    DELETE FROM tasks
    WHERE id NOT IN (
        SELECT id FROM tasks
        ORDER BY created DESC
        LIMIT 1000
    );
END;

CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created DATETIME NOT NULL,
    modified DATETIME NOT NULL,
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

CREATE TABLE audiobooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    subtitle TEXT NOT NULL,
    authors TEXT,
    narrator TEXT,
    genre TEXT,
    isbn TEXT,
    asin TEXT,
    language TEXT,
    year INTEGER,
    duration INTEGER,
    chapters JSONB,
    cover BLOB,
    file_path TEXT NOT NULL,
    m4b_path TEXT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audiobooks_files (
	audiobook_id INTEGER,
	file_id INTEGER,
	FOREIGN KEY (audiobook_id) REFERENCES audiobooks(id),
	FOREIGN KEY (file_id) REFERENCES files(id)
);

CREATE TABLE libraries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    paths JSONB,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message TEXT NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER limit_logs_after_insert
AFTER INSERT ON logs
BEGIN
    DELETE FROM logs
    WHERE id NOT IN (
        SELECT id FROM logs
        ORDER BY created DESC
        LIMIT 1000
    );
END;

CREATE TABLE properties (
	name TEXT PRIMARY KEY,
	value TEXT NOT NULL DEFAULT ''
);

CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	is_admin BOOLEAN NOT NULL DEFAULT false
);