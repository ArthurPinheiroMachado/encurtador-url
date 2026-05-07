package db

import (
	"fmt"
	"golang/db/drivers/postgres"
	"golang/internal/util"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type _DatabaseDriver interface {
	InitDatabase() []string
	QueryString(user, pass, host, name, port string) string
	InsertStage() string
	LastPosition() string
	GetUrls() string
	DeleteUrl() string
	SaveUrl() string
	GetUrlByUrl() string
	UpdateAccesses() string
}

type Database struct {
	driver _DatabaseDriver
	sqlx   *sqlx.DB
}

func ConstructDatabase() (*Database, error) {
	trace := util.CreateErrorContext("db.ConstructDatabase")

	tp, name, host, port, pass, user, _env_errors := getDatabaseVars()
	if _env_errors != nil {
		return nil, trace.Apply(_env_errors)
	}

	driver, _driver := chooseDriver(tp)
	if _driver != nil {
		return nil, trace.Apply(_driver)
	}

	sqlx, _sqlx := sqlx.Connect(tp, driver.QueryString(user, pass, host, name, port))
	if _sqlx != nil {
		return nil, trace.Apply(_sqlx)
	}

	return &Database{driver, sqlx}, nil
}

func getDatabaseVars() (string, string, string, string, string, string, error) {
	trace := util.CreateErrorContext("db.getDatabaseVars")

	tp, _tp := util.EnvAsResult("DB_TYPE")
	name, _name := util.EnvAsResult("DB_NAME")
	host, _host := util.EnvAsResult("DB_HOST")
	port, _port := util.EnvAsResult("DB_PORT")
	pass, _pass := util.EnvAsResult("DB_PASS")
	user, _user := util.EnvAsResult("DB_USER")

	_env_errors := trace.Join(_tp, _name, _host, _port, _pass, _user)
	return tp, name, host, port, pass, user, _env_errors
}

func chooseDriver(db string) (_DatabaseDriver, error) {
	switch db {
	case "postgres":
		return _DatabaseDriver(postgres.Driver{}), nil

	default:
		return nil, fmt.Errorf("%s não é um driver válido", db)
	}
}

func (db *Database) Migrate() error {
	last_position := -1
	_ = db.sqlx.Get(&last_position, db.driver.LastPosition())

	for idx, query := range db.driver.InitDatabase() {
		if idx <= last_position {
			continue
		}

		_, _mig_error := db.sqlx.Exec(query)

		if _mig_error != nil {
			return fmt.Errorf("unable to migrate idx=%d error=%w", idx, _mig_error)
		}

		fmt.Println("idx=%w query=%w", idx, query)
		_, _stage_error := db.sqlx.Exec(db.driver.InsertStage(), idx, query)

		if _stage_error != nil {
			return fmt.Errorf("unable to insert migration stage idx=%d error=%w", idx, _stage_error)
		}
	}
	return nil
}
