-- name: CreateSubject :execresult
insert into Subjects (code, name, description, class) 
values (?, ?, ?, ?);

-- name: GetSubjectByID :one
select id, code, name, description, class from Subjects
where id = ?;

-- name: ListSubjects :many
select id, code, name, description, class from Subjects
order by code;

-- name: ListSubjectsByClass :many
select id, code, name, description, class from Subjects
where class = ?
order by code;

-- name: UpdateSubject :execresult
update Subjects
set code = ?, name = ?, description = ?, class = ?
where id = ?;