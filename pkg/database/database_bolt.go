package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

func openBolt(constr string) (VMDatabase, error) {
	file, err := filepath.Abs(constr)
	_, err = os.Stat(file)
	if err != nil {
		return VMDatabase{nil}, fmt.Errorf("Unable to open '%s': %s", file, err)
	}
	d, err := bolt.Open(file, 0600, nil)

	return VMDatabase{d}, err
}
