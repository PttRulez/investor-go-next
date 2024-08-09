-- +goose Up
-- +goose StatementBegin
create table if not exists opinions_on_positions (
		id serial primary key,
		opinion_id integer references opinions(id) not null,
		portfolio_id integer references portfolios(id) not null
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table opinions_on_positions;
-- +goose StatementEnd
