-- name: CreateStudent :execresult
insert into Students (name) value (?);

-- name: GetStudentByID :one
select id, name from Students
where id = ?;

-- name: GetAllStudents :many
select id, name from Students;

-- name: GetAllStudentsWithNameLike :many
select id, name from Students
where name like ?;

-- name: UpdateStudent :execresult
update Students
set name = ?
where id = ?;

