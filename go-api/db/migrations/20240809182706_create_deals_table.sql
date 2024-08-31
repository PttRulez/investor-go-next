-- +goose Up
-- +goose StatementBegin
create table if not exists deals (
	amount integer not null,
	commission numeric(10, 2) not null, 
	date date not null,
	exchange varchar(20) not null,
	id serial primary key,
	portfolio_id integer references portfolios(id) not null,
	price numeric(10, 6) not null,
	security_type varchar(20) not null, 
	ticker varchar(50) not null,
	type varchar(50) not null,
	user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table deals;
-- +goose StatementEnd
