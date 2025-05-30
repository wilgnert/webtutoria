-- +goose up
CREATE TABLE StudentSubjectCompletion (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT NOT NULL,
    subject_id INT NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- To track when the completion occurred

    -- Composite unique constraint to ensure a student completes a subject only once
    UNIQUE KEY unique_student_subject (student_id, subject_id),

    -- Foreign key to the Students table
    CONSTRAINT fk_ssc_student
        FOREIGN KEY (student_id)
        REFERENCES Students(id)
        ON DELETE CASCADE,

    -- Foreign key to the Subjects table
    CONSTRAINT fk_ssc_subject
        FOREIGN KEY (subject_id)
        REFERENCES Subjects(id)
        ON DELETE CASCADE
);

-- +goose down
DROP TABLE IF EXISTS StudentSubjectCompletion;