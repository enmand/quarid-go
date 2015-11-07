package database

import (
	"upper.io/db"
	"upper.io/db/postgresql"
)

func openPostgresql(constr string) (Database, error) {
	settings, err := postgresql.ParseURL(constr)
	if err != nil {
		return Database{nil}, err
	}

	d, err := db.Open(postgresql.Adapter, settings)
	if err != nil {
		return Database{nil}, err
	}

	return Database{d}, nil
}
