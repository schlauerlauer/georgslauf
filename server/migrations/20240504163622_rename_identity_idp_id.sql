-- +goose Up
-- +goose StatementBegin
alter table identities rename column kratos_id to idp_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table identities rename column idp_id to kratos_id;
-- +goose StatementEnd
