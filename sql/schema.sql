CREATE TABLE sources (
    id INTEGER PRIMARY KEY ,
    name TEXT NOT NULL UNIQUE,
    props TEXT DEFAULT '{}', -- json
    added DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE directories (
    id INTEGER PRIMARY KEY ,
    name TEXT NOT NULL,
    source_id INTEGER NOT NULL,
    added DATETIME DEFAULT CURRENT_TIMESTAMP,
    props TEXT DEFAULT '{}', -- json
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    UNIQUE(name, source_id)
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    dir_id INTEGER NOT NULL,
    added DATETIME DEFAULT CURRENT_TIMESTAMP,
    props TEXT DEFAULT '{}', -- json
    FOREIGN KEY (dir_id) REFERENCES directories(id) ON DELETE CASCADE,
    UNIQUE(name, dir_id)
);

CREATE TABLE images (
    id INTEGER PRIMARY KEY,
    data BLOB NOT NULL,
    dir_id INTEGER NOT NULL,
    props TEXT DEFAULT '{}',
    FOREIGN KEY (dir_id) REFERENCES directories(id) ON DELETE CASCADE,
    UNIQUE(dir_id)
);

CREATE TABLE provider_searches (
    query TEXT,
    results TEXT DEFAULT '[]',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    dir_id INTEGER NOT NULL,
    UNIQUE(query, dir_id)
);

CREATE VIEW entries AS
SELECT
    a.id AS id,
    a.name AS dir,
    c.name || '/' || a.name AS full_path,
    a.added AS added,
    a.props AS props,
    json_group_array(json_object('id',b.id,'name',b.name,'ext',b.props->>'$.ext')) AS files,
    d.id AS image_id
FROM (
    SELECT id,name,dir_id,props
    FROM files
    WHERE props->>'$.ext' IN ('mp3','m4a','m4b')
) b
JOIN directories a ON b.dir_id=a.id
JOIN sources c ON c.id=a.source_id
LEFT JOIN images d ON d.dir_id=a.id
GROUP BY a.id, a.name;


-- CREATE TABLE users (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT NOT NULL,
--     hash TEXT NOT NULL,
--     admin BOOLEAN DEFAULT FALSE,
--     preferences JSON DEFAULT '{}'
-- );

-- CREATE TABLE sessions (
--     id TEXT PRIMARY KEY,
--     user_id INTEGER NOT NULL,
--     created DATETIME DEFAULT CURRENT_TIMESTAMP,
--     expires DATETIME,
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
-- );


-- CREATE TABLE books (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     title TEXT COLLATE NOCASE NOT NULL,
--     author TEXT COLLATE NOCASE,
--     narrator TEXT COLLATE NOCASE,
--     isbn TEXT COLLATE NOCASE,
--     asin TEXT COLLATE NOCASE,
--     genre TEXT,
--     language TEXT,
--     year INTEGER,
--     duration INTEGER,
--     chapters JSONB,
--     provider TEXT,
--     added DATETIME DEFAULT CURRENT_TIMESTAMP
-- );

