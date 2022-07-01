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

func (db *Database) RawPut(collection, key string, value []byte) error {
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

func (db *Database) RawGet(collection, key string) (*Entry, error) {
	if db.collections[collection] == nil {
		return nil, ErrCollectionNotFound
	}
	mEntry, ok := db.collections[collection].Load(key)
	if !ok {
		return nil, ErrEntryNotFound
	}
	return mEntry.(*Entry), nil
}

func (db *Database) PutUser(key string, u *User) error {
	userData, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return db.RawPut("users", key, userData)
}

func (db *Database) GetUser(key string) (*User, error) {
	userData, err := db.RawGet("users", key)
	if err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal(userData.Value, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
