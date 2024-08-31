-- +goose Up
-- +goose StatementBegin
create table if not exists opinions_on_positions (
		id serial primary key,
		opinion_id integer references opinions(id) not null,
		position_id integer references positions(id) not null,
		CONSTRAINT unique_position_opinion UNIQUE (position_id, opinion_id)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table opinions_on_positions;
-- +goose StatementEnd
