-- +goose Up
-- +goose StatementBegin
create index idx_schedule_start on schedule("start");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_schedule_start;
-- +goose StatementEnd
