-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex_shares ADD COLUMN currency varchar(10) not null default 'RUB';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex_shares DROP COLUMN currency;
-- +goose StatementEnd
