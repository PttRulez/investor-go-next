-- +goose Up
-- +goose StatementBegin
create table if not exists moex_bonds (
	id serial primary key,
	board varchar(50) not null,
	coupon_percent decimal(10, 2) not null,
	coupon_value decimal(10, 2) not null,
	coupon_frequency integer not null,
	engine varchar(50) not null,
	face_value integer not null,
	issue_date date not null,
	lotsize integer not null,
	market varchar(50) not null,
	mat_date date not null,
	name varchar(100) not null,
	shortname varchar(50) not null,
	secid varchar(20) not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table moex_bonds;
-- +goose StatementEnd
