-- +goose up
CREATE TABLE Students (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS Students;