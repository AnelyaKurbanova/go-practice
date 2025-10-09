CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    amount NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    spent_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    note TEXT,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);

CREATE INDEX idx_expenses_user_id ON expenses (user_id);
CREATE INDEX idx_expenses_user_spent_at ON expenses (user_id, spent_at);