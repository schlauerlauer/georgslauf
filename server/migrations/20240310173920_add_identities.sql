-- +goose Up
-- +goose StatementBegin
create table identities(
	id int primary key not null,
	kratos_id text not null,
	email text not null,
	created_at int not null,
	tribe_id int not null,
	"role" int not null default 0,
	foreign key(tribe_id) references tribes(id)
);
create index idx_identities_kratos_id on identities("kratos_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table identities;
-- +goose StatementEnd
