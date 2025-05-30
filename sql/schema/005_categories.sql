-- +goose up
CREATE TABLE Categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS Categories;