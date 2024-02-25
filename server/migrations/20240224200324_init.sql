-- +goose Up
-- +goose StatementBegin
create table tribes(
	id int primary key not null,
	updated_at int not null,
	"name" text not null,
	short text,
	dpsg text,
	image_id text,
	email_domain text,
	stavo_email text
);
create index idx_tribes_name on tribes("name");

create table stations(
	id int primary key not null,
	created_at int not null,
	updated_at int not null,
	"name" text not null,
	short text,
	"size" int not null,
	tribe_id int not null,
	lati real,
	long real,
	image_id text,
	foreign key(tribe_id) references tribes(id)
);

create table groups(
	id int primary key not null,
	created_at int not null,
	updated_at int not null,
	"name" text not null,
	short text,
	"size" int,
	grouping int not null,
	tribe_id int not null,
	image_id text,
	foreign key(tribe_id) references tribes(id)
);

create table points_to_stations(
	id int primary key not null,
	created_at int not null,
	updated_at int not null,
	group_id int not null,
	station_id int not null,
	"value" int not null,
	foreign key(group_id) references groups(id),
	foreign key(station_id) references stations(id)
);
create unique index idx_pts on points_to_stations(group_id, station_id);

create table points_to_groups(
	id int primary key not null,
	created_at int not null,
	updated_at int not null,
	station_id int not null,
	group_id int not null,
	"value" int not null,
	foreign key(station_id) references stations(id),
	foreign key(group_id) references groups(id)
);
create unique index idx_ptg on points_to_groups(station_id, group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tribes;
drop table stations;
drop table groups;

drop index idx_pts;
drop table points_to_stations;
drop index idx_ptg;
drop table points_to_groups;

-- +goose StatementEnd
