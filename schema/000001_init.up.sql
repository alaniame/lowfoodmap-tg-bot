CREATE TABLE IF NOT EXISTS product_categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS carb_types (
    carb_id SERIAL PRIMARY KEY,
    carb_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    category_id INT NOT NULL,
    stage INT NOT NULL,
    portion_high INT,
    portion_medium INT,
    portion_low INT,
    portion_size VARCHAR(255),
    FOREIGN KEY (category_id) REFERENCES product_categories (category_id)
);

CREATE TABLE IF NOT EXISTS product_carb_types (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    carb_id INT NOT NULL,
    FOREIGN KEY (carb_id) REFERENCES carb_types(carb_id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

INSERT INTO carb_types (carb_name) VALUES ('фруктаны') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('фруктоза') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('сорбитол') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('маннитол') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('галактаны') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('галактоолигосахариды(гос)') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('лактоза') ON CONFLICT (carb_name) DO NOTHING;
INSERT INTO carb_types (carb_name) VALUES ('талактаны') ON CONFLICT (carb_name) DO NOTHING;

INSERT INTO product_categories (category_name) VALUES ('Фрукты и ягоды') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Овощи и грибы') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Молоко, молочные продукты и их заменители') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Орехи, бобовые и семена') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Напитки') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Специи, травы, соусы и спреды') ON CONFLICT (category_name) DO NOTHING;
INSERT INTO product_categories (category_name) VALUES ('Хлеб, злаки и макаронные изделия') ON CONFLICT (category_name) DO NOTHING;
