-- name: GetTribes :many
select
	t.*
	,ti.id as icon
from tribes t
left join tribe_icons ti on ti.id = t.id
;

-- name: GetStationCategories :many
select
	id
	,name
	,max
from station_categories;

-- name: GetSchedule :many
select *
from schedule
order by "start" asc;

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

-- name: InsertGroup :one
insert into groups (name, size, tribe_id, grouping, comment, created_by, updated_by, created_at, updated_at)
values (?,?,?,?,?,?,?,?,?)
returning id;

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

-- name: UpdateTribeRole :exec
update tribe_roles
set tribe_role = ?
where id = ?;

-- name: GetUsersRoleLargerNone :many
select
	email
	,last_login
	,created_at
	,role
	,firstname
	,lastname
from users
where role > 0;

-- name: GetUsersRoleNone :many
select
	id
	,email
	,created_at
	,firstname
	,lastname
from users
where role = 0
order by created_at desc;

-- TODO join icons
-- name: GetTribeRolesOpen :many
select
	tr.id
	,tr.created_at
	,t.name as tribe_name
	,u.email as user_email
from tribe_roles tr
inner join users u on tr.user_id = u.id
inner join tribes t on tr.tribe_id = t.id
where
	tribe_role = 0
order by tr.created_at desc;

-- name: GetUserIdByExt :one
select
	id
	,role
	,last_login
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
values (?,?)
on conflict(id) do update set
	image = excluded.image;

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

-- NTH where id = ?
-- name: UpdateSettings :exec
update settings
set
	updated_at = unixepoch()
	,updated_by = ?
	,data = ?;
