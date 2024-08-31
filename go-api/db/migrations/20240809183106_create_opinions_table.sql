-- +goose Up
-- +goose StatementBegin
create table if not exists opinions (
	id serial primary key,
	date date not null,
	exchange varchar(50) not null,
	expert_id integer references experts(id) not null,
	text text not null,
	security_id integer not null,
	security_type varchar(50) not null,
	source_link varchar(120),
	target_price numeric(15, 6),
	ticker varchar(50) not null,
	type varchar(50) not null,
	user_id integer references users(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table opinions;
-- +goose StatementEnd
