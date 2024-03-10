-- name: ListTribes :many
select *
from tribes;

-- name: GetTribe :one
select *
from tribes
where id = ?
limit 1;

-- name: DeleteTribe :exec
delete from tribes
where id = ?
limit 1;

-- name: GetSchedule :many
select *
from schedule
order by "start" asc;

-- name: ListGroups :many
select *
from groups;

-- name: GetGroup :one
select *
from groups
where id = ?
limit 1;

-- name: DeleteGroup :exec
delete from groups
where id = ?
limit 1;

-- name: ListGroupsByTribe :many
select *
from groups
where tribe_id = ?;
