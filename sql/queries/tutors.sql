-- name: CreateTutor :execresult
insert into Tutors (name) value (?);

-- name: GetTutorByID :one
select id, name from Tutors
where id = ?;

-- name: GetAllTutors :many
select id, name from Tutors;

-- name: GetAllTutorsWithNameLike :many
select id, name from Tutors
where name like ?;

-- name: UpdateTutor :execresult
update Tutors
set name = ?
where id = ?;

