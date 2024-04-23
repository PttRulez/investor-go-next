package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pttrulez/investor-go/internal/config"
	"github.com/pttrulez/investor-go/internal/types"
)

func NewPostgresRepo(cfg config.PostgresConfig) (*types.Repository, error) {
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("sql.Open err", err)
		return nil, err
	}

	createAllTables(db)

	return &types.Repository{
		Cashout: NewCashoutPostgres(db),
		Deal: types.DealRepository{
			MoexBond:  NewMoexBondDealPostgres(db),
			MoexShare: NewMoexShareDealPostgres(db),
		},
		Deposit: NewDepositPostgres(db),
		Expert:  NewExpertPostgres(db),
		Moex: types.MoexRepository{
			Bond:  NewMoexBondsPostgres(db),
			Share: NewMoexSharesPostgres(db),
		},
		Portfolio: NewPortfolioPostgres(db),
		Position: types.PositionRepository{
			MoexBond:  NewMoexBondPositionPostgres(db),
			MoexShare: NewMoexSharePositionPostgres(db),
		},
		User: NewUserPostgres(db),
	}, nil
}
func createAllTables(db *sql.DB) {
	createUsersTable(db)
	createPortfoliosTable(db)
	createExpertsTable(db)
	createMoexBondDealsTable(db)
	createMoexShareDealsTable(db)
	createOpinionsTable(db)
	createOpinionsOnPositionsTable(db)
	createMoexBondPositionsTable(db)
	createMoexSharePositionsTable(db)
	createDepositsTable(db)
	createCashoutsTable(db)

	//  ---------------------- MOEX ----------------------
	createMoexBondsTable(db)
	createMoexSharesTable(db)
	createMoexCurrenciesTable(db)
}

// ---------------------- FUNCS FOR TABLES CREATION  ----------------------
func createCashoutsTable(db *sql.DB) {
	queryString := `create table if not exists cashouts (
		id serial primary key,
		amount integer not null,
		date date not null,
		portfolio_id integer references portfolios(id) not null
		)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createCashoutsTable] err", err)
	}
}
func createMoexBondDealsTable(db *sql.DB) {
	queryString := `create table if not exists moex_bond_deals (
		id serial primary key,
		amount integer not null,
		date date not null,
		portfolio_id integer references portfolios(id) not null,
		price numeric(10, 2) not null,
		security_id integer references moex_bonds(id) not null,
		type varchar(50) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexBondDealsTable] err", err)
	}
}
func createMoexShareDealsTable(db *sql.DB) {
	queryString := `create table if not exists moex_share_deals (
		id serial primary key,
		amount integer not null,
		date date not null,
		portfolio_id integer references portfolios(id) not null,
		price numeric(10, 2) not null,
		security_id integer references moex_shares(id) not null,
		type varchar(50) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexShareDealsTable] err", err)
	}
}
func createDepositsTable(db *sql.DB) {
	queryString := `create table if not exists deposits (
		id serial primary key,
		amount integer not null,
		date date not null,
		portfolio_id integer references portfolios(id) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createDepositsTable] err", err)
	}
}
func createExpertsTable(db *sql.DB) {
	queryString := `create table if not exists experts (
		id serial primary key,
		avatar_url varchar(100),
		name varchar(50) not null,
		user_id integer references users(id) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createExpertsTable] err", err)
	}
}
func createOpinionsTable(db *sql.DB) {
	queryString := `create table if not exists opinions (
		id serial primary key,
		date date not null,
		exchange varchar(50) not null,
		expert_id integer references experts(id) not null,
		text text not null,
		security_id integer not null,
		security_type varchar(50) not null,
		source_link varchar(120),
		target_price numeric(10, 2),
		type varchar(50) not null,
		user_id integer references users(id) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createOpinionsTable] err", err)
	}
}
func createOpinionsOnPositionsTable(db *sql.DB) {
	queryString := `create table if not exists opinions_on_positions (
		id serial primary key,
		opinion_id integer references opinions(id) not null,
		portfolio_id integer references portfolios(id) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createOpinionsOnPositionsTable] err", err)
	}
}
func createPortfoliosTable(db *sql.DB) {
	queryString := `create table if not exists portfolios (
		id serial primary key,
		compound boolean not null,
		name varchar(50) not null,
		user_id integer references users(id) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createPortfoliosTable] err: ", err)
	}
}

func createMoexBondPositionsTable(db *sql.DB) {
	queryString := `create table if not exists moex_bond_positions (
		id serial primary key,
		amount integer not null,
		average_price numeric(10, 2) not null,
		comment text,
		portfolio_id integer references portfolios(id) not null,
		security_id integer not null,
		target_price numeric(10, 2))`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexBondPositionsTable] err", err)
	}
}
func createMoexSharePositionsTable(db *sql.DB) {
	queryString := `create table if not exists moex_share_positions (
		id serial primary key,
		amount integer not null,
		average_price numeric(10, 2) not null,
		comment text,
		portfolio_id integer references portfolios(id) not null,
		security_id integer not null,
		target_price numeric(10, 2))`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexSharePositionsTable] err", err)
	}
}
func createUsersTable(db *sql.DB) {
	queryString := `create table if not exists users (
		id serial primary key,
		email varchar(50) unique not null,
		hashed_password varchar(100) not null,
		name varchar(50) not null,
		role varchar(10) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createUsersTable] err", err)
	}
}
func createMoexBondsTable(db *sql.DB) {
	queryString := `create table if not exists moex_bonds (
		id serial primary key,
		board varchar(50) not null,
		engine varchar(50) not null,
		market varchar(50) not null,
		name varchar(100) not null,
		shortname varchar(50) not null,
		isin varchar(10) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexBondsTable] err", err)
	}
}
func createMoexCurrenciesTable(db *sql.DB) {
	queryString := `create table if not exists moex_currencies (
		id serial primary key,
		board varchar(50) not null,
		engine varchar(50) not null,
		market varchar(50) not null,
		name varchar(120) not null,
		shortname varchar(50) not null,
		ticker varchar(10) not null
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexSharesTable] err", err)
	}
}
func createMoexSharesTable(db *sql.DB) {
	queryString := `create table if not exists moex_shares (
		id serial primary key,
		board varchar(50) not null,
		engine varchar(50) not null,
		market varchar(50) not null,
		name varchar(120) not null,
		shortname varchar(50) not null,
		ticker varchar(10) not null unique
	)`

	_, err := db.Exec(queryString)
	if err != nil {
		log.Fatal("[createMoexSharesTable] err", err)
	}
}
