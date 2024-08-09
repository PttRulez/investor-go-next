-- +goose Up
-- +goose StatementBegin
create table if not exists portfolios (
	id serial primary key,
	compound boolean not null,
	name varchar(50) not null,
	user_id integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table portfolios;
-- +goose StatementEnd
