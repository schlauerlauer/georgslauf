-- name: GetTribes :many
select
	t.*
	,ti.id as icon
from tribes t
left join tribe_icons ti on ti.id = t.id
;

-- -- name: CreateTribe :one
-- insert into tribes (updated_at, created_by, "name")
-- values (?,?,?)
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
	g.id
	,g.created_at
	,g.updated_at
	,g.name
	,g.size
	,g.grouping
	,g.comment
	,u.firstname
from groups g
left join users u on g.updated_by = u.id
where g.tribe_id = ?
order by g.created_at desc;

-- name: UpdateGroup :exec
update groups
set
	name = ?
	,size = ?
	,grouping = ?
	,comment = ?
	,updated_at = ?
	,updated_by = ?
where
	id = ?
	and tribe_id = ?;

-- name: GetTribeRoleWithIcon :one
select
	tr.tribe_role
	,tr.tribe_id
	,ti.id as icon_id
from tribe_roles tr
left join tribe_icons ti on ti.id = tr.id
where
	user_id = ?
limit 1;

-- name: GetTribeRoleByTribe :one
select
	tribe_role
from tribe_roles
where
	user_id = ?
	and tribe_id = ?
limit 1;

-- name: CreateTribeRole :exec
insert into tribe_roles (user_id, tribe_id, tribe_role, created_by)
values (?,?,?,?);

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

-- name: CreateTribeIcon :exec
insert into tribe_icons (id, created_by, image)
values (?,?,?);

-- name: UpdateTribeIcon :exec
update tribe_icons
set
	created_at = unixepoch()
	,created_by = ?
	,image = ?
where id = ?;

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
	,updated_by = ?
	,data = ?;

--where id = ?