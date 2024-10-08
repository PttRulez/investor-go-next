-- +goose Up
-- +goose StatementBegin
create table if not exists expenses (
		amount numeric(10, 2) not null,
		date date not null,
		description varchar(100) not null,
		id serial primary key,
		portfolio_id integer references portfolios(id) not null,
		user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table expenses;
-- +goose StatementEnd
