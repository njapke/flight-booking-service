package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var ErrDatabase = errors.New("database error")
var ErrCollectionNotFound = fmt.Errorf("%w: collection not found", ErrDatabase)
var ErrEntryNotFound = fmt.Errorf("%w: entry not found", ErrDatabase)

type Entry struct {
	Key   string
	Value []byte
}

type Database struct {
	collections map[string]*sync.Map
}

func New() *Database {
	return &Database{
		collections: make(map[string]*sync.Map),
	}
}

func (db *Database) rawPut(collection, key string, value []byte) error {
	if db.collections[collection] == nil {
		db.collections[collection] = &sync.Map{}
	}
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	db.collections[collection].Store(key, entry)
	return nil
}

func (db *Database) rawGet(collection, key string) (*Entry, error) {
	if db.collections[collection] == nil {
		return nil, ErrCollectionNotFound
	}
	mEntry, ok := db.collections[collection].Load(key)
	if !ok {
		return nil, ErrEntryNotFound
	}
	return mEntry.(*Entry), nil
}

func (db *Database) Put(u Model) error {
	userData, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return db.rawPut("users", u.Key(), userData)
}

func (db *Database) Get(key string, val Model) error {
	userData, err := db.rawGet(val.Collection(), key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(userData.Value, val); err != nil {
		return err
	}
	return nil
}
