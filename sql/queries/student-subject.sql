-- name: ListStudentSubjectCompletionsByStudent :many
SELECT * FROM StudentSubjectCompletion
WHERE student_id = ?;

-- name: ListStudentSubjectCompletionsBySubject :many
SELECT * FROM StudentSubjectCompletion
WHERE subject_id = ?;

-- name: ListStudentSubjectCompletions :many
SELECT * FROM StudentSubjectCompletion;

-- name: CreateStudentSubjectCompletion :execresult
INSERT INTO StudentSubjectCompletion (student_id, subject_id)
VALUES (?, ?);

-- name: GetStudentSubjectCompletionByID :one
SELECT * FROM StudentSubjectCompletion
WHERE id = ?;

-- name: DeleteStudentSubjectCompletion :exec
DELETE FROM StudentSubjectCompletion
WHERE id = ?;