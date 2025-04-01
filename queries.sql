-- name: GetTribes :many
select
	t.*
	,ti.id as icon
from tribes t
left join tribe_icons ti on ti.id = t.id
;

-- name: GetTribesName :many
select
	t.id
	,t.short
	,t.name
from tribes t;

-- name: GetStationCategories :many
select
	sc.id
	,sc.name
	,sc.max
	,count(s.id) as count
from station_categories sc
left join stations s on s.category_id = sc.id
group by sc.id;

-- name: GetCategoryOfStation :one
select
	s.category_id
	,sc.name
from stations s
left join station_categories sc on s.category_id = sc.id
where s.id = ?
limit 1;

-- name: GetStationCategory :one
select
	id
	,name
	,max
from station_categories
where id = ?
limit 1;

-- name: GetStationCategoryCount :one
select
	count(id)
from stations
where category_id = ?
limit 1;

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
	,g.vegan
	,u.firstname
	,ui.image as user_image
from groups g
left join users u on g.updated_by = u.id
left join user_icons ui on g.updated_by = ui.id
where g.tribe_id = ?
order by g.created_at desc;

-- name: GetStationOverview :one
select
	count(id)
	,sum(size) as size
	,sum(vegan) as vegan
from stations
limit 1;

-- name: GetGroupOverview :one
select
	count(id)
	,sum(size) as size
	,sum(vegan) as vegan
from groups
limit 1;

-- name: GetStation :one
select
	s.id
	,s.created_at
	,s.updated_at
	,s.name
	,s.abbr
	,s.size
	,s.category_id
	,s.description
	,s.requirements
	,s.vegan
	,s.tribe_id
	,s.updated_by
	,sc.name as category_name
	,u.firstname
	,ui.image as user_image
from stations s
left join station_categories sc on category_id = sc.id
left join users u on s.updated_by = u.id
left join user_icons ui on s.updated_by = ui.id
where s.id = ?
limit 1;

-- name: GetGroup :one
select
	g.id
	,g.created_at
	,g.updated_at
	,g.name
	,g.abbr
	,g.size
	,g.comment
	,g.grouping
	,g.vegan
	,g.tribe_id
	,g.updated_by
	,u.firstname
	,ui.image as user_image
from groups g
left join users u on g.updated_by = u.id
left join user_icons ui on g.updated_by = ui.id
where g.id = ?
limit 1;

-- name: GetStationsDetails :many
select
	s.id
	,s.name
	,s.abbr
	,sc.name as category
	,t.name as tribe
	,s.size
	,ti.id as tribe_icon
from stations s
left join tribes t on s.tribe_id = t.id
left join tribe_icons ti on ti.id = t.id
left join station_categories sc on s.category_id = sc.id
order by s.created_at desc;

-- name: GetTribeNameIcon :one
select
	t.name
	,ti.id as tribe_icon
from tribes t
left join tribe_icons ti on ti.id = t.id
where
	t.id = ?
limit 1;

-- name: GetGroupsDetails :many
select
	g.id
	,g.name
	,g.grouping
	,g.size
	,g.abbr
	,t.name as tribe
	,ti.id as tribe_icon
from groups g
left join tribes t on g.tribe_id = t.id
left join tribe_icons ti on ti.id = t.id
order by g.created_at desc;

-- name: GetGroupsHost :many
select
	g.id
	,g.name
	,g.grouping
	,g.tribe_id
from groups g
order by g.tribe_id asc;

-- name: UpdateGroup :exec
update groups
set
	name = ?
	,size = ?
	,grouping = ?
	,comment = ?
	,updated_at = ?
	,updated_by = ?
	,vegan = ?
where
	id = ?
	and tribe_id = ?;

-- name: UpdateGroupHost :exec
update groups
set
	name = ?
	,size = ?
	,grouping = ?
	,comment = ?
	,updated_at = ?
	,updated_by = ?
	,vegan = ?
	,abbr = ?
	,tribe_id = ?
where
	id = ?;

-- name: InsertStation :one
insert into stations (name, size, tribe_id, category_id, description, requirements, created_by, created_at, updated_at, updated_by, vegan, pref_loc)
values (?,?,?,?,?,?,?,?,?,?,?,?)
returning id;

-- name: InsertGroup :one
insert into groups (name, size, tribe_id, grouping, comment, created_by, updated_by, created_at, updated_at, vegan)
values (?,?,?,?,?,?,?,?,?,?)
returning id;

-- name: GetTribeRoleWithIcon :one
select
	tr.tribe_role
	,tr.tribe_id
	,ti.id as icon_id
	,tr.accepted_at
	,t.name
from tribe_roles tr
left join tribe_icons ti on ti.id = tr.id
left join tribes t on tr.tribe_id = t.id
where
	tr.user_id = ?
limit 1;

-- name: GetTribeRoleById :one
select
	tr.id
	,tr.tribe_role
	,tr.accepted_at
	,u.email
	,ui.image
	,u.firstname
	,u.lastname
	,t.name as tribe_name
	,t.short
	,t.email_domain
	,ti.id as tribe_icon
from tribe_roles tr
inner join users u on tr.user_id = u.id
inner join tribes t on tr.tribe_id = t.id
left join user_icons ui on ui.id = tr.user_id
left join tribe_icons ti on ti.id = tr.tribe_id
where
	tr.id = ?
limit 1;

-- name: GetTribeRoleByTribe :one
select
	tr.tribe_role
	,ti.id as icon_id
	,tr.accepted_at
	,t.name
from tribe_roles tr
left join tribe_icons ti on ti.id = tr.id
left join tribes t on tr.tribe_id = t.id
where
	tr.user_id = ?
	and tr.tribe_id = ?
limit 1;

-- name: CreateTribeRole :exec
insert into tribe_roles (user_id, tribe_id, tribe_role, created_by, accepted_at)
values (?,?,?,?,?);

-- name: UpdateTribeRole :exec
update tribe_roles
set
	tribe_role = ?
	,accepted_at = ?
	,updated_by = ?
	,updated_at = ?
where id = ?;

-- name: UpdateUserRole :exec
update users
set role = ?
where id = ?;

-- name: GetUserRole :one
select
	role
from users
where id = ?
limit 1;

-- name: GetUsersRoleLargerNone :many
select
	u.id
	,u.email
	,u.last_login
	,u.created_at
	,u.role
	,u.firstname
	,u.lastname
	,ui.image
from users u
left join user_icons ui on ui.id = u.id
where u.role > 0;

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

-- NTH multiple smaller queries
-- name: GetTribeRolesOpen :many
select
	tr.id
	,tr.created_at
	,t.name as tribe_name
	,u.email as user_email
	,u.firstname
	,u.lastname
	,ui.image as user_icon
	,ti.id as tribe_icon_id
from tribe_roles tr
inner join users u on u.id = tr.user_id
inner join tribes t on t.id = tr.tribe_id
left join user_icons ui on ui.id = tr.user_id
left join tribe_icons ti on ti.id = tr.tribe_id
where
	tribe_role = 0
	and accepted_at is null
order by tr.created_at desc;

-- name: GetTribeRolesAssigned :many
select
	tr.id
	,tr.tribe_id
	,u.email
	,tr.tribe_role
from tribe_roles tr
inner join users u on u.id = tr.user_id
where
	accepted_at is not null
	or tr.tribe_role = -1
order by tr.tribe_id asc;

-- name: GetUserIdByExt :one
select
	id
	,role
	,last_login
from users
where
	ext_id = ?
limit 1;

-- name: GetStationsByTribe :many
select
	s.id
	,s.created_at
	,s.updated_at
	,s.name
	,s.abbr
	,s.size
	,s.category_id
	,s.description
	,s.requirements
	,s.vegan
	,s.pref_loc
	,sc.name as category_name
	,u.firstname
	,ui.image as user_image
from stations s
left join station_categories sc on category_id = sc.id
left join users u on s.updated_by = u.id
left join user_icons ui on s.updated_by = ui.id
where s.tribe_id = ?
order by s.created_at desc;

-- name: UpdateStation :exec
update stations
set
	updated_at = ?
	,updated_by = ?
	,name = ?
	,size = ?
	,category_id = ?
	,description = ?
	,requirements = ?
	,vegan = ?
	,pref_loc = ?
where
	id = ?
	and tribe_id = ?;

-- name: UpdateStationHost :exec
update stations
set
	updated_at = ?
	,updated_by = ?
	,name = ?
	,size = ?
	,category_id = ?
	,description = ?
	,requirements = ?
	,vegan = ?
	,tribe_id = ?
	,abbr = ?
where
	id = ?;

-- name: GetStationsHost :many
select
	id
	,name
	,tribe_id
from stations
order by tribe_id asc;

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

-- name: DeleteGroup :exec
delete from groups
where id = ?;

-- name: DeleteStation :exec
delete from stations
where id = ?;

-- name: UpdateStationCategory :exec
update station_categories
set
	name = ?
	,max = ?
where id = ?;

-- name: DeleteStationCategory :exec
delete from station_categories
where id = ?;

-- name: InsertStationCateogy :one
insert into station_categories (name, max)
values (?,?)
returning id;