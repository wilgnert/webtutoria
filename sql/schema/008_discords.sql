-- +goose up
CREATE TABLE StudentDiscords (
  student_id INT NOT NULL,
  discord_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (student_id), -- Assuming one discord_id per student
  UNIQUE KEY unique_discord_student (discord_id),
  FOREIGN KEY (student_id) REFERENCES Students(id) ON DELETE CASCADE
);

CREATE TABLE TutorDiscords (
  tutor_id INT NOT NULL,
  discord_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (tutor_id), -- Assuming one discord_id per tutor
  UNIQUE KEY unique_discord_tutor (discord_id),
  FOREIGN KEY (tutor_id) REFERENCES Tutors(id) ON DELETE CASCADE
);


-- +goose down
DROP TABLE IF EXISTS StudentDiscords;
DROP TABLE IF EXISTS TutorDiscords;