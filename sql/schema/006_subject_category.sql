-- +goose up
CREATE TABLE SubjectCategory (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subject_id INT NOT NULL,
    category_id INT NOT NULL,

    CONSTRAINT fk_subject
        FOREIGN KEY (subject_id)
        REFERENCES Subjects(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_category
        FOREIGN KEY (category_id)
        REFERENCES Categories(id)
        ON DELETE CASCADE,

    UNIQUE KEY unique_subject_category (subject_id, category_id)
);

-- +goose down
DROP TABLE IF EXISTS SubjectCategory;