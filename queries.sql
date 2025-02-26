-- name: GetTribes :many
select
	t.*
	,ti.id as icon
from tribes t
left join tribe_icons ti on ti.id = t.id
;

-- -- name: CreateTribe :one
-- insert into tribes (updated_at, "name")
-- values (?, ?)
-- returning *;

-- -- name: GetTribe :one
-- select *
-- from tribes
-- where id = ?
-- limit 1;

-- -- name: DeleteTribe :exec
-- delete from tribes
-- where id = ?
-- limit 1;

-- name: GetSchedule :many
select *
from schedule
order by "start" asc;

-- -- name: GetGroups :many
-- select *
-- from groups;

-- -- name: GetGroup :one
-- select *
-- from groups
-- where id = ?
-- limit 1;

-- -- name: DeleteGroup :exec
-- delete from groups
-- where id = ?
-- limit 1;

-- TODO image_id
-- name: GetGroupsByTribe :many
select
	id
	,created_at
	,updated_at
	,name
	,size
	,grouping
from groups
where tribe_id = ?
order by created_at desc;

-- TODO preferred tribe limit 1
-- name: GetTribeRole :one
select
	tribe_role
	,tribe_id
from tribe_roles
where user_id = ?
limit 1;

-- name: CreateTribeRole :exec
insert into tribe_roles (user_id, tribe_id, tribe_role)
values (?,?,?);

-- name: GetUserIdByExt :one
select
	id
	,role
from users
where
	ext_id = ?
limit 1;

-- TODO imageId
-- name: GetStationsByTribe :many
select
	id
	,created_at
	,updated_at
	,name
	,size
	,lati
	,long
	,description
	,requirements
from stations
where tribe_id = ?
order by created_at desc;

-- TODO picture update
-- name: CreateUser :one
insert into users (ext_id, username, firstname, lastname, email, last_login)
values (?,?,?,?,?,?)
returning id;

-- name: UpdateUser :exec
update users
set
	last_login = ?
	,username = ?
	,firstname = ?
	,lastname = ?
	,email = ?
where id = ?;

-- name: GetTribeByEmail :one
select
	id
	,name
	,dpsg
from tribes
where email_domain = ?
limit 1;

-- name: CreateUserIcon :exec
insert into user_icons (id, image)
values (?,?);

-- name: GetUserIcon :one
select
	image
from user_icons
where id = ?
limit 1;

-- name: GetTribeIcon :one
select
	image
from tribe_icons
where id = ?
limit 1;

-- name: GetSettings :one
select
	*
from settings
limit 1;

-- name: InsertSettings :exec
insert into settings (data) values (?);

-- name: UpdateSettings :exec
update settings
set
	updated_at = unixepoch()
	,data = ?;

--where id = ?