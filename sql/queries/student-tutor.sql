-- name: ListStudentTutors :many
SELECT
    *
FROM
    StudentTutor
ORDER BY
    student_id, tutor_id;

-- name: ListStudentTutorsByStudent :many
select
    *
from
    StudentTutor
where
    student_id = ?
ORDER BY
    tutor_id;


-- name: ListStudentTutorsByTutor :many
select
    *
from
    StudentTutor
where
    tutor_id = ?
ORDER BY
    student_id;

-- name: CreateStudentTutor :execresult
insert into StudentTutor (student_id, tutor_id)
values (?, ?);

-- name: GetStudentTutorByID :one
select
    *
from
    StudentTutor
where
    id = ?;

-- name: DeleteStudentTutorByID :exec
delete from
    StudentTutor
where
    id = ?;
