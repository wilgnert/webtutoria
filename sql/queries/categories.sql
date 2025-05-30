-- name: GetCategoryByName :one
select id, name from Categories
where name = ?;

-- name: CreateCategory :execresult
insert into Categories (name) values (?);

-- name: GetCategoryByID :one
select id, name from Categories
where id = ?;

-- name: GetAllCategories :many
select id, name from Categories
order by name;

-- name: UpdateCategory :execresult
update Categories
set name = ?
where id = ?;

-- name: DeleteCategory :execresult
delete from Categories
where id = ?;

