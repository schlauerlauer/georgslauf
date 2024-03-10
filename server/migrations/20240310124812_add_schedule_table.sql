-- +goose Up
-- +goose StatementBegin
create table schedule(
	id int primary key not null,
	"start" int not null,
	"end" int,
	"name" text not null,
	"about" boolean not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table schedule;
-- +goose StatementEnd
