-- +goose Up
-- +goose StatementBegin
create table if not exists users (
		id serial primary key,
		email varchar(50) unique not null,
		hashed_password varchar(100) not null,
		name varchar(50) not null,
		role varchar(10) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
