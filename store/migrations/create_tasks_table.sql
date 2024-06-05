CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP,
    completed_at TIMESTAMP
);
