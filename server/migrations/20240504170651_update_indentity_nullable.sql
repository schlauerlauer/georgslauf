-- +goose Up
-- +goose StatementBegin
drop table identities;
create table identities(
	id int primary key not null,
	idp_id text not null,
	email text,
	created_at int not null,
	tribe_id int,
	"role" int not null default 0,
	foreign key(tribe_id) references tribes(id)
);
create index idx_identities_idp_id on identities("idp_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table identities;
create table identities(
	id int primary key not null,
	idp_id text not null,
	email text not null,
	created_at int not null,
	tribe_id int not null,
	"role" int not null default 0,
	foreign key(tribe_id) references tribes(id)
);
create index idx_identities_idp_id on identities("idp_id");
-- +goose StatementEnd
