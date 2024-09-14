### DDL
```sql
CREATE DATABASE db_gomrp;

USE db_gomrp;

CREATE TABLE users (
    user_id VARCHAR(50) PRIMARY KEY, -- Custom employee ID
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    permission_int VARCHAR(50) NOT NULL,
    archived BOOLEAN DEFAULT FALSE  -- Soft delete flag
);

CREATE TABLE inventory (
    inventory_id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(100) NOT NULL,
    item_code VARCHAR(50) UNIQUE NOT NULL,
    quantity INT DEFAULT 0,
    status ENUM('OK', 'NG') DEFAULT 'OK',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE
);

CREATE TABLE manufacturing_orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT,
    quantity INT NOT NULL,
    status ENUM('Pending', 'In Progress', 'Completed') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (product_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE manufacturing_components (
    component_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    component_item_id INT,
    quantity_required INT,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (order_id) REFERENCES manufacturing_orders(order_id),
    FOREIGN KEY (component_item_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE sales (
    sale_id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT,
    quantity INT NOT NULL,
    sale_status ENUM('Pending', 'Paid', 'Delivered') DEFAULT 'Pending',
    sale_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (item_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE purchases (
    purchase_id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT,
    quantity INT NOT NULL,
    purchase_status ENUM('Ordered', 'Received') DEFAULT 'Ordered',
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (item_id) REFERENCES inventory(inventory_id)
);

CREATE TABLE income (
    income_id INT AUTO_INCREMENT PRIMARY KEY,
    sale_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    received_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sale_id) REFERENCES sales(sale_id)
);

CREATE TABLE outcome (
    outcome_id INT AUTO_INCREMENT PRIMARY KEY,
    purchase_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    spent_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (purchase_id) REFERENCES purchases(purchase_id)
);

-- in case i forgot
CREATE USER 'gomrp'@'%' IDENTIFIED BY 'gomrp123';
GRANT ALL PRIVILEGES ON *.* to 'gomrp'@'%';
SELECT user, host FROM mysql.user;
```

### Restore DB
```sh
# mysqldump -u gomrp -pgomrp123 db_gomrp > db.sql
mysql -u gomrp -pgomrp123 db_gomrp < db.sql
```

### Run
```sh
# go get github.com/gofiber/fiber/v2 github.com/go-sql-driver/mysql github.com/joho/godotenv
go mod tidy
go run main.go
```

### etc in case i forgot
```sh
go mod init rizaldyaristyo-fiber-boiler
```