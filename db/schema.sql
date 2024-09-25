CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    priority INTEGER,
    status TEXT,
    params JSONB,
    result JSONB,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created DATETIME NOT NULL,
    modified DATETIME NOT NULL
);

CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT COLLATE NOCASE NOT NULL,
    authors TEXT COLLATE NOCASE,
    narrator TEXT COLLATE NOCASE,
    isbn TEXT COLLATE NOCASE,
    asin TEXT COLLATE NOCASE,
    genre TEXT,
    language TEXT,
    year INTEGER,
    duration INTEGER,
    chapters JSONB,
    added DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    type TEXT,
    size INTEGER,
    book_id INTEGER,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE books_files (
    book_id INTEGER,
    file_id INTEGER,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE books_images (
    book_id INTEGER NOT NULL,
    image_id INTEGER NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

CREATE TRIGGER limit_tasks_after_insert
AFTER INSERT ON tasks BEGIN
    DELETE FROM tasks
    WHERE id NOT IN (
        SELECT id FROM tasks
        ORDER BY created DESC
        LIMIT 100000
    );
END;