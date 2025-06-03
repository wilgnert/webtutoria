-- +goose up
ALTER TABLE Tutors
ADD COLUMN channel_id TEXT NOT NULL,
ADD COLUMN role_id TEXT NOT NULL;

-- +goose down
ALTER TABLE Tutors
DROP COLUMN channel_id,
DROP COLUMN role_id;