-- +goose Up
-- +goose StatementBegin
create table if not exists dividends (
		date date not null,
		exchange varchar(20) not null,
		id serial primary key,
		payment_per_share numeric(10, 2) not null,
		payment_period varchar(50) not null,
		portfolio_id integer references portfolios(id) not null,
		shares_count integer,
		ticker varchar(50) not null,
		user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table dividends;
-- +goose StatementEnd
