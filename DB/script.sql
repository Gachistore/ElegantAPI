INSERT INTO account (number, first_name, last_name, e_mail, encrypted_password)
VALUES
    (1,'Admin', 'Cool', 'adm@cool', '1337'),


-- Populating the product table
INSERT INTO product (name, price, measurements, description, packaging)
VALUES
    ('Bed', 599.99, '160x200 cm', 'Comfortable bed with orthopedic mattress', 'Box'),
    ('Wardrobe', 399.99, '200x150x60 cm', 'Spacious wardrobe for storing clothes and belongings', 'Bag'),
    ('Dining Table', 199.99, '120x80x75 cm', 'Wooden dining table for dining area', 'Box');

-- Populating the review table
INSERT INTO review (accID, prodID, rating_given, text)
VALUES
    (1, 1, 4.0, 'Great bed, comfortable and beautiful'),
    (2, 2, 4.5, 'The wardrobe exceeded my expectations, the assembly quality is excellent'),
    (3, 3, 5.0, 'The table fits perfectly into the interior, the material quality is pleasing');

-- Populating the category table
INSERT INTO category (name)
VALUES
    ('Bedroom Furniture'),
    ('Living Room Furniture'),
    ('Dining Furniture');

-- Populating the product_category table
INSERT INTO product_category (prodID, category_name)
VALUES
    (1, 'Bedroom Furniture'),
    (2, 'Living Room Furniture'),
    (3, 'Dining Furniture');