CREATE table categories
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id INTEGER,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT unique_user_category UNIQUE (user_id, name)
);

CREATE INDEX idx_categories_user_id ON categories (user_id);