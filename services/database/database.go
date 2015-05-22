package database

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

const DB_DIR = "dbs/%s"
const DB_FILEPATH = "%s.db"

type Database *bolt.DB

func GetDatabase(dir, path string) (Database, error) {
	file := fmt.Sprintf("%s/%s/%s.db", DB_DIR, dir, DB_FILEPATH)
	_, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to open '%s': %s", file, err)
	}

	full_dir := fmt.Sprintf(DB_DIR, dir)
	if err := os.Mkdir(dir, os.ModeDir); err != nil {
		return nil, fmt.Errorf(
			"Could not create DB directory %s: %s",
			full_dir,
			err,
		)
	}
	return bolt.Open(file, 0600, nil)
}
