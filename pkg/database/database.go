package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var ErrDatabase = errors.New("database error")
var ErrCollectionNotFound = fmt.Errorf("%w: collection not found", ErrDatabase)
var ErrEntryNotFound = fmt.Errorf("%w: entry not found", ErrDatabase)

type Entry struct {
	Key   string
	Value []byte
}

type Model interface {
	Collection() string
	Key() string
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

func (db *Database) Put(m Model) error {
	entryValue, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return db.rawPut(m.Collection(), m.Key(), entryValue)
}

func (db *Database) Get(key string, val Model) error {
	entry, err := db.rawGet(val.Collection(), key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(entry.Value, val); err != nil {
		return err
	}
	return nil
}

func (db *Database) Keys(collection string) ([]string, error) {
	if db.collections[collection] == nil {
		return nil, ErrCollectionNotFound
	}
	keys := make([]string, 0)
	db.collections[collection].Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}

func (db *Database) Entries(collection string) ([]*Entry, error) {
	if db.collections[collection] == nil {
		return nil, ErrCollectionNotFound
	}
	entries := make([]*Entry, 0)
	db.collections[collection].Range(func(key, value interface{}) bool {
		entries = append(entries, value.(*Entry))
		return true
	})
	return entries, nil
}

func (db *Database) Values(forModel Model) ([]Model, error) {
	entries, err := db.Entries(forModel.Collection())
	if err != nil {
		return nil, err
	}
	values := make([]Model, len(entries))
	for i, entry := range entries {
		val := reflect.New(reflect.TypeOf(forModel).Elem()).Interface().(Model)
		if err := json.Unmarshal(entry.Value, val); err != nil {
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}
