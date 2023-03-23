CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY, 
    title TEXT,
    content TEXT,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
)
