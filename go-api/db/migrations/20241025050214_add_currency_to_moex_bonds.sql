-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex_bonds ADD COLUMN currency varchar(10) not null default 'RUB';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex_bonds DROP COLUMN currency;
-- +goose StatementEnd
