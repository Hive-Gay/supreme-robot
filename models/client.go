package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	"github.com/rubenv/sql-migrate"
)

var logger = loggo.GetLogger("models")

type Client struct {
	client *sqlx.DB
}

func NewClient(dbEngine string, doMigration bool) (*Client, error) {
	client, err := sqlx.Connect("postgres", dbEngine)
	if err != nil {
		return nil, err
	}

	if doMigration {
		logger.Debugf("Loading Migrations")
		migrations := &migrate.HttpFileSystemMigrationSource{
			FileSystem: pkger.Dir("/models/migrations"),
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
	}
	return &Client{
		client: client,
	}, nil
}
