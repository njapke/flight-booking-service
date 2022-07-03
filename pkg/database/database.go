package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var ErrDatabase = errors.New("database error")
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
	data *sync.Map
}

func New() *Database {
	return &Database{
		data: &sync.Map{},
	}
}

func (db *Database) getPrefixedKey(collection, key string) string {
	return fmt.Sprintf("%s/%s", collection, key)
}

func (db *Database) rawPut(collection, key string, value []byte) error {
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	db.data.Store(db.getPrefixedKey(collection, key), entry)
	return nil
}

func (db *Database) rawGet(collection, key string) (*Entry, error) {
	mEntry, ok := db.data.Load(db.getPrefixedKey(collection, key))
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
	keys := make([]string, 0)
	db.data.Range(func(key, value interface{}) bool {
		prefix, k, _ := strings.Cut(key.(string), "/")
		if prefix != collection {
			return true
		}
		keys = append(keys, k)
		return true
	})
	return keys, nil
}

func (db *Database) Entries(collection string) ([]*Entry, error) {
	entries := make([]*Entry, 0)
	db.data.Range(func(key, value interface{}) bool {
		if !strings.HasPrefix(key.(string), collection) {
			return true
		}
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
