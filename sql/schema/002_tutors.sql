-- +goose up
CREATE TABLE Tutors (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS Tutors;