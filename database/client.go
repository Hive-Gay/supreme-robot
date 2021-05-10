package database

import (
	"github.com/Hive-Gay/supreme-robot/config"
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	"github.com/rubenv/sql-migrate"
)

var logger = loggo.GetLogger("database")

type Client struct {
	db *sqlx.DB
}

func NewClient(cfg *config.Config) (*Client, error) {
	client, err := sqlx.Connect("postgres", cfg.PostgresDsn)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Loading Migrations")
	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/database/migrations"),
	}

	logger.Debugf("Applying Migrations")
	n, err := migrate.Exec(client.DB, "postgres", migrations, migrate.Up)
	if n > 0 {
		logger.Infof("Applied %d migrations!", n)
	}
	if err != nil {
		logger.Criticalf("Could not migrate database: %s", err)
		return nil, err
	}

	return &Client{
		db: client,
	}, nil
}
