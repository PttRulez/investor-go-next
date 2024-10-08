-- +goose Up
-- +goose StatementBegin
create table if not exists coupons (
		bonds_count integer,
		coupon_amount numeric(10, 2) not null,
		date date not null,
		exchange varchar(20) not null,
		id serial primary key,
		payment_period varchar(50) not null,
		portfolio_id integer references portfolios(id) not null,
		ticker varchar(50) not null,
		user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table coupons;
-- +goose StatementEnd
