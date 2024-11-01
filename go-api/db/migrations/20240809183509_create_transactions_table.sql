-- +goose Up
-- +goose StatementBegin
create table if not exists transactions (
		id serial primary key,
		amount numeric(15, 2) not null,
		date date not null,
		portfolio_id integer references portfolios(id) not null,
		type varchar(50) not null,
		user_id  integer references users(id) not null
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transactions;
-- +goose StatementEnd
