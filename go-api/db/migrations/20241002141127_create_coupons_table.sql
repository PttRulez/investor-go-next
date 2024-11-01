-- +goose Up
-- +goose StatementBegin
create table if not exists coupons (
		bonds_count integer,
		date date not null,
		exchange varchar(20) not null,
		id serial primary key,
		payment_period varchar(50) not null,
		portfolio_id integer references portfolios(id) not null,
		shortname varchar(255)  not null,
		tax_paid numeric(15,2)  not null,
		ticker varchar(50) not null,
		total_payment numeric(15, 2) not null,
		user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table coupons;
-- +goose StatementEnd
