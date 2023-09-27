CREATE TABLE IF NOT EXISTS privileges (
    id SERIAL PRIMARY KEY,
    privilege_title VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS privileged_users (
    user_id INTEGER,
    privilege_id INTEGER,
    assigned_at TIMESTAMP NOT NULL,

    FOREIGN KEY (privilege_id) REFERENCES privileges(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, privilege_id)
);