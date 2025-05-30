-- name: GetStudentDiscordByDiscordID :one
select * from StudentDiscords
where discord_id = ?;

-- name: GetTutorDiscordByDiscordID :one
select * from TutorDiscords
where discord_id = ?;

-- name: ListStudentDiscords :many
select * from StudentDiscords
order by created_at desc;

-- name: ListTutorDiscords :many
select * from TutorDiscords
order by created_at desc;

-- name: CreateStudentDiscord :exec
insert into StudentDiscords (student_id, discord_id)
values (?, ?)
on duplicate key update
discord_id = values(discord_id);

-- name: CreateTutorDiscord :exec
insert into TutorDiscords (tutor_id, discord_id)
values (?, ?)
on duplicate key update
discord_id = values(discord_id);

-- name: GetStudentDiscordByStudentID :one
select * from StudentDiscords
where student_id = ?;

-- name: GetTutorDiscordByTutorID :one
select * from TutorDiscords
where tutor_id = ?;

-- name: DeleteStudentDiscordByStudentID :exec
delete from StudentDiscords
where student_id = ?;
-- name: DeleteTutorDiscordByTutorID :exec
delete from TutorDiscords
where tutor_id = ?;