CREATE TABLE IF NOT EXISTS categories (
  id   SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
  id          SERIAL PRIMARY KEY,
  name        TEXT NOT NULL,
  category_id INT  NOT NULL REFERENCES categories(id),
  price       INT  NOT NULL
);

INSERT INTO categories (name) VALUES
  ('phones'),
  ('laptops'),
  ('accessories')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, category_id, price)
SELECT 'iPhone',       c1.id, 400000 FROM categories c1 WHERE c1.name='phones'
UNION ALL
SELECT 'Pixel',        c2.id, 300000 FROM categories c2 WHERE c2.name='phones'
UNION ALL
SELECT 'MacBook Air',  c3.id, 700000 FROM categories c3 WHERE c3.name='laptops'
UNION ALL
SELECT 'USB-C Cable',  c4.id, 5000   FROM categories c4 WHERE c4.name='accessories';
