-- +goose Up
-- +goose StatementBegin
create table if not exists dividends_paid (
		amount integer not null,
		date date not null,
		exchange varchar(20) not null,
		id serial primary key,
		payment_year integer not null,
		portfolio_id integer references portfolios(id) not null,
		security_id integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table dividends_paid;
-- +goose StatementEnd
