-- +goose Up
-- +goose StatementBegin
create table if not exists positions (
	id serial primary key,
	amount integer not null,
	average_price numeric(10, 2) not null,
	board varchar(20) not null, 
	comment text,
	exchange varchar(20) not null,
	portfolio_id integer references portfolios(id) not null,
	security_type varchar(20),
	short_name varchar(20),
	ticker varchar(50) not null,
	target_price numeric(10, 2),
	user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table positions;
-- +goose StatementEnd
