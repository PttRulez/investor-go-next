-- +goose Up
-- +goose StatementBegin
create table if not exists experts (
	id serial primary key,
	avatar_url varchar(100),
	name varchar(50) not null,
	user_id  integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table experts;
-- +goose StatementEnd
