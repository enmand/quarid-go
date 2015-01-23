package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/vmihailenco/msgpack.v2"
)

type DB map[string]interface{}
type Database struct {
	Path string
	File *os.File
	Data DB
}

const DB_FILEPATH = "%s.db"

var (
	CannotCreateDbFile = errors.New(
		"The database path provided cannot be created: %s",
	)
	CannotOpenDbError = errors.New(
		"The database path provided cannot be opened: %s",
	)
	CannotReadDbFileError = errors.New(
		"The database file can be found, but not read: %s",
	)
	CannotSaveDbFile = errors.New(
		"The database file cannot be saved: %s",
	)
)

func GetDatabase(path string) (*Database, error) {
	file, err := openDb(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
	}()

	dbFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(CannotReadDbFileError.Error(), err)
	}

	database := &Database{}
	database.Path = path
	database.File = file

	if err := msgpack.Unmarshal(dbFile, &database.Data); err != nil {
		return &Database{
			Path: path,
			File: file,
			Data: DB{},
		}, nil
	}

	return database, nil
}

func (d Database) Save() error {
	filePath := fmt.Sprintf(DB_FILEPATH, d.Path)
	dbFile, err := msgpack.Marshal(d.Data)
	if err != nil {
		return fmt.Errorf(CannotSaveDbFile.Error(), err)
	}
	ioutil.WriteFile(filePath, dbFile, os.ModePerm)

	return nil
}

func openDb(path string) (*os.File, error) {
	var file *os.File
	var err error
	filePath := fmt.Sprintf(DB_FILEPATH, path)

	stat, err := os.Stat(filePath)

	switch stat {
	case nil:
		file, err = os.Create(filePath)
		if err != nil {
			return nil, CannotCreateDbFile
		}
		break
	default:
		file, err = os.Open(filePath)
		if err != nil {
			return nil, CannotOpenDbError
		}
	}

	return file, nil
}
