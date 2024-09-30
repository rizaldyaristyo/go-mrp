DROP DATABASE db_gomrp;

CREATE DATABASE db_gomrp;

USE db_gomrp;

CREATE TABLE users (
    user_id VARCHAR(50) PRIMARY KEY, -- Custom employee ID
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    permission_val VARCHAR(50) NOT NULL,
    archived BOOLEAN DEFAULT FALSE  -- Soft delete flag
);

CREATE TABLE vendors (
    vendor_id INT AUTO_INCREMENT PRIMARY KEY,
    vendor_name VARCHAR(100) NOT NULL,
    vendor_address VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE inventory (
    inventory_id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(100) NOT NULL,
    vendor_id INT,
    item_code VARCHAR(50) UNIQUE NOT NULL,
    item_code_2 VARCHAR(50) UNIQUE,
    item_type ENUM('Product', 'Raw Material', 'Processed Material', 'Consumable') DEFAULT 'Product',
    sellable BOOLEAN DEFAULT TRUE,
    purchaseable BOOLEAN DEFAULT TRUE,
    manufacturable BOOLEAN DEFAULT FALSE,
    price DECIMAL(10, 2) NOT NULL,
    price_2 DECIMAL(10, 2),
    currency VARCHAR(5) DEFAULT 'IDR',
    quantity INT DEFAULT 0,
    minimum_stock_warning INT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (vendor_id) REFERENCES vendors(vendor_id)
);

CREATE TABLE manufacturing_orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    manufacture_order_number VARCHAR(50) UNIQUE NOT NULL,
    product_id INT,
    quantity INT NOT NULL,
    status ENUM('Pending', 'In Progress', 'Completed', 'Cancelled') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (product_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE manufacturing_recipes (
    material_id INT AUTO_INCREMENT PRIMARY KEY,
    material_inventory_id INT NOT NULL, -- Raw Material
    needed_to_produce_product_id INT NOT NULL, -- Manufactured Product (can be product or processed material)
    material_quantity_to_produce_product INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete
    FOREIGN KEY (material_inventory_id) REFERENCES inventory(inventory_id),
    FOREIGN KEY (needed_to_produce_product_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE customers (
    customer_id INT AUTO_INCREMENT PRIMARY KEY,
    customer_name VARCHAR(100) NOT NULL,
    customer_bank_name VARCHAR(50),
    customer_bank_account_number VARCHAR(50),
    customer_cc_number VARCHAR(50),
    customer_address VARCHAR(255) NOT NULL,
    customer_email VARCHAR(100),
    customer_phone VARCHAR(30),
    customer_tax_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE -- Soft delete flag
);

CREATE TABLE sales (
    sale_id INT AUTO_INCREMENT PRIMARY KEY,
    sales_order_number VARCHAR(50) UNIQUE NOT NULL,
    item_id INT,
    quantity INT NOT NULL,
    sent_quantity INT DEFAULT 0, -- Tracks how much has been delivered
    sale_price_per_unit DECIMAL(10, 2) NOT NULL,
    tax_percent DECIMAL(10, 2) NOT NULL,
    customer_id INT,
    sales_channel ENUM('Online', 'In-Store', 'Distributor') DEFAULT 'Online',
    payment_method ENUM('Cash', 'Credit Card', 'Bank Transfer'),
    payment_status ENUM('Pending', 'Paid Partially', 'Paid') DEFAULT 'Pending',
    delivery_status ENUM('Pending', 'Delivered Partially', 'Delivered') DEFAULT 'Pending',
    canceled BOOLEAN DEFAULT FALSE, -- Marks if the order is canceled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_date TIMESTAMP,
    delivery_date TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (item_id) REFERENCES inventory(inventory_id),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);


CREATE TABLE purchases (
    purchase_id INT AUTO_INCREMENT PRIMARY KEY,
    purchase_order_number VARCHAR(50) UNIQUE NOT NULL,
    item_id INT,
    quantity INT NOT NULL,
    purchase_price_per_unit DECIMAL(10, 2) NOT NULL,
    purchase_status ENUM('Ordered', 'Received', 'Cancelled') DEFAULT 'Ordered',
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (item_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE income (
    income_id INT AUTO_INCREMENT PRIMARY KEY,
    sale_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(5) DEFAULT 'IDR',
    received_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (sale_id) REFERENCES sales(sale_id)
);

CREATE TABLE outcome (
    outcome_id INT AUTO_INCREMENT PRIMARY KEY,
    purchase_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(5) DEFAULT 'IDR',
    spent_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE, -- Soft delete flag
    FOREIGN KEY (purchase_id) REFERENCES purchases(purchase_id)
);

INSERT INTO users (user_id, username, password, email, permission_val, archived)
VALUES
    ('USR0001', 'john_doe', '$2y$04$MkJJMmwUOZXdtGJ1fTjVEuRBjwdNQNz5XkfXSNHCDc1Ymiv3dntgi', 'john.doe@example.com', '3021', FALSE),
    ('USR0002', 'jane_smith', '$2y$04$PYI1sj4oJJxTolfCWLknp.0gOxYdAPDKADHq3AMrtkjy2vo1tjcH6', 'jane.smith@example.com', '1203', FALSE),
    ('USR0003', 'alice_johnson', '$2y$04$sjhaJ2ul43EptKazMFATTut/MBoTw1FAuCDkTX6H2V/8Xs7BOS9za', 'alice.johnson@example.com', '3333', TRUE);

INSERT INTO vendors (vendor_name, vendor_address, tax_id)
VALUES
('ABC Corporation', 'Jl. ABC No. 123', '123456789'),
('XYZ Company', 'Jl. XYZ No. 456', '987654321');

INSERT INTO inventory (item_name, vendor_id, item_code, item_code_2, item_type, sellable, purchaseable, manufacturable, price, price_2, currency, quantity, minimum_stock_warning)
VALUES
('Bottle Cap', 1, 'CAP001', 'CAP-A', 'Raw Material', FALSE, TRUE, TRUE, 10.00, NULL, 'IDR', 1000, 200),
('Plastic Bottle', 1, 'BOTT001', 'BOTT-P', 'Raw Material', FALSE, TRUE, TRUE, 50.00, NULL, 'IDR', 500, 100),
('Mineral Water', 2, 'MW001', 'WATER-1L', 'Product', TRUE, FALSE, TRUE, 200.00, 220.00, 'IDR', 300, 50),
('Packaging Box', 2, 'BOX001', NULL, 'Consumable', FALSE, TRUE, FALSE, 100.00, NULL, 'IDR', 100, 20);

INSERT INTO manufacturing_orders (manufacture_order_number, product_id, quantity, status)
VALUES
('MO001', 3, 200, 'Pending'),
('MO002', 3, 100, 'Pending');

INSERT INTO manufacturing_recipes (material_inventory_id, needed_to_produce_product_id, material_quantity_to_produce_product)
VALUES
(1, 3, 100),
(2, 3, 50),
(3, 4, 50);

INSERT INTO customers (customer_name, customer_bank_name, customer_bank_account_number, customer_cc_number, customer_address, customer_email, customer_phone, customer_tax_id)
VALUES
('John Doe', 'Bank A', '1234567890', '1234-5678-9012-3456', 'Jl. Merah No. 1', 'john.doe@example.com', '08123456789', 'TAX001'),
('Jane Smith', 'Bank B', '0987654321', '9876-5432-1098-7654', 'Jl. Biru No. 2', 'jane.smith@example.com', '08234567890', 'TAX002'),
('Alice Johnson', 'Bank C', '5551234567', NULL, 'Jl. Kuning No. 3', 'alice.johnson@example.com', '08345678901', 'TAX003'),
('Bob Brown', 'Bank D', '4441234567', '1234-4321-5678-9012', 'Jl. Hijau No. 4', 'bob.brown@example.com', '08456789012', 'TAX004'),
('Charlie Green', 'Bank E', '7771234567', NULL, 'Jl. Putih No. 5', 'charlie.green@example.com', '08567890123', 'TAX005');

INSERT INTO sales (sales_order_number, item_id, quantity, sent_quantity, sale_price_per_unit, tax_percent, customer_id, sales_channel, payment_method, payment_status, delivery_status, canceled, payment_date, delivery_date)
VALUES
('SO001', 3, 50, 0, 300.00, 10.00, 1, 'Online', 'Credit Card', 'Pending', 'Delivered', FALSE, NULL, NULL),
('SO002', 3, 30, 30, 300.00, 10.00, 2, 'In-Store', 'Cash', 'Paid', 'Delivered', FALSE, '2024-09-02 15:30:00', '2024-09-02 15:30:00'),
('SO004', 2, 10, 0, 50.00, 10.00, 3, 'Online', 'Credit Card', 'Pending', 'Delivered', FALSE, NULL, NULL),
('SO005', 4, 5, 5, 100.00, 10.00, 2, 'In-Store', 'Cash', 'Paid', 'Delivered', FALSE, '2024-09-03 13:45:00', '2024-09-03 14:00:00'),
('SO006', 1, 100, 0, 10.00, 10.00, 1, 'Online', 'Credit Card', 'Pending', 'Delivered', TRUE, NULL, NULL);

INSERT INTO purchases (purchase_order_number, item_id, quantity, purchase_price_per_unit, purchase_status)
VALUES
('PO001', 1, 500, 300, 'Ordered'),
('PO002', 2, 200, 300, 'Received'),
('PO003', 4, 50, 300, 'Cancelled');

INSERT INTO income (sale_id, amount, currency, received_date)
VALUES
(2, 6600.00, 'IDR', '2024-09-01 10:00:00'),
(3, 4400.00, 'IDR', '2024-09-02 15:30:00');

INSERT INTO outcome (purchase_id, amount, currency, spent_date)
VALUES
(1, 5000.00, 'IDR', '2024-09-03 09:00:00'),
(2, 10000.00, 'IDR', '2024-09-04 13:45:00');