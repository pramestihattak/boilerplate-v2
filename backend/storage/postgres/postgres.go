package postgres

import (
	"fmt"
	"time"

	"boilerplate-v2/util"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// Postgres storage driver
	_ "github.com/lib/pq"
)

// Storage provides a wrapper around an sql database and provides
// required methods for interacting with the database
type Storage struct {
	logger logrus.FieldLogger
	db     *sqlx.DB
}

func NewStorage(logger logrus.FieldLogger, config *viper.Viper) (*Storage, error) {
	dbString, err := util.NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}
	db, err := NewDbConn(logger, dbString)
	if err != nil {
		return nil, err
	}

	return &Storage{logger: logger, db: db}, nil
}

func NewDbConn(logger logrus.FieldLogger, dbstring string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres '%s': %v", dbstring, err)
	}
	// TODO: is this a sane default?
	// The current max_connections in postgres is 100.
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
