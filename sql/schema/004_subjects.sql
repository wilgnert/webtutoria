-- +goose up
CREATE TABLE Subjects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  class VARCHAR(255) NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS Subjects;