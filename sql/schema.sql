
CREATE TABLE directories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    root TEXT,
    added DATETIME DEFAULT CURRENT_TIMESTAMP,
    props JSON DEFAULT '{}'
);

CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    ext TEXT,
    directory_id INTEGER NOT NULL,
    details JSON DEFAULT '{}',
    FOREIGN KEY (directory_id) REFERENCES directories(id) ON DELETE CASCADE,
    UNIQUE(name, directory_id)
);

CREATE VIEW directories_details AS
SELECT
    a.id,
    a.name,
    a.root,
    a.added,
    a.props,
    json_group_array(b.name) AS entries,
    json_group_array(
        json_object('id',b.id,'name',b.name,'ext',b.ext,'details',b.details)) AS entries_details
FROM directories a
JOIN entries b
ON b.directory_id=a.id
GROUP BY a.id, a.name;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    hash TEXT NOT NULL,
    admin BOOLEAN DEFAULT FALSE,
    preferences JSON DEFAULT '{}'
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT COLLATE NOCASE NOT NULL,
    author TEXT COLLATE NOCASE,
    narrator TEXT COLLATE NOCASE,
    isbn TEXT COLLATE NOCASE,
    asin TEXT COLLATE NOCASE,
    genre TEXT,
    language TEXT,
    year INTEGER,
    duration INTEGER,
    chapters JSONB,
    provider TEXT,
    added DATETIME DEFAULT CURRENT_TIMESTAMP
);

