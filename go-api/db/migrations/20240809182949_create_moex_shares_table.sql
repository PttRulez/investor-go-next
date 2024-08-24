-- +goose Up
-- +goose StatementBegin
create table if not exists moex_shares (
	id serial primary key,
	board varchar(50) not null,
	engine varchar(50) not null,
	lotsize integer not null,
	market varchar(50) not null,
	name varchar(120) not null,
	shortname varchar(50) not null,
	secid varchar(10) not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table moex_shares;
-- +goose StatementEnd
