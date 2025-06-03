-- name: CreateTutor :execresult
insert into Tutors (name, role_id, channel_id) value (?, ?, ?);

-- name: GetTutorByID :one
select * from Tutors
where id = ?;

-- name: GetAllTutors :many
select * from Tutors;

-- name: GetAllTutorsWithNameLike :many
select * from Tutors
where name like ?;

-- name: UpdateTutor :execresult
update Tutors
set name = ?,
    role_id = ?,
    channel_id = ?
where id = ?;

