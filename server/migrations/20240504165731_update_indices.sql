-- +goose Up
-- +goose StatementBegin
drop index idx_identities_kratos_id;
create index idx_identities_idp_id on identities("idp_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_identities_idp_id;
create index idx_identities_kratos_id on identities("kratos_id");
-- +goose StatementEnd
