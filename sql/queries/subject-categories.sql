-- name: CreateSubjectCategory :execresult
insert into SubjectCategory (subject_id, category_id) values (?, ?);

-- name: ListCategoriesBySubjectID :many
select c.name as name from SubjectCategory sc
join Categories c on sc.category_id = c.id
where sc.subject_id = ?
order by c.name;

-- name: DeleteSubjectCategoriesBySubjectID :exec
delete from SubjectCategory
where subject_id = ?;