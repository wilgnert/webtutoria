-- +goose up
CREATE TABLE StudentTutor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT NOT NULL,
    tutor_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraint for student_id
    CONSTRAINT fk_student
        FOREIGN KEY (student_id)
        REFERENCES Students(id)
        ON DELETE CASCADE, -- If a student is deleted, remove their entries from StudentTutor

    -- Foreign key constraint for tutor_id
    CONSTRAINT fk_tutor
        FOREIGN KEY (tutor_id)
        REFERENCES Tutors(id)
        ON DELETE CASCADE -- If a tutor is deleted, remove their entries from StudentTutor
);

-- +goose down
DROP TABLE IF EXISTS StudentTutor;