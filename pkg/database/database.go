package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

type Entry struct {
	Key   string
	Value []byte
}

type Model interface {
	Collection() string
	Key() string
}

type Database struct {
	db *badger.DB
}

func New() (*Database, error) {
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (db *Database) getPrefixedKey(collection, key string) []byte {
	return []byte(fmt.Sprintf("%s/%s", collection, key))
}

func (db *Database) rawPut(collection, key string, value []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(db.getPrefixedKey(collection, key), value)
	})
}

func (db *Database) rawGet(collection, key string) ([]byte, error) {
	var val []byte
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.getPrefixedKey(collection, key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (db *Database) Put(m Model) error {
	entryValue, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return db.rawPut(m.Collection(), m.Key(), entryValue)
}

func (db *Database) Get(key string, val Model) error {
	value, err := db.rawGet(val.Collection(), key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(value, val); err != nil {
		return err
	}
	return nil
}

func (db *Database) Values(forModel Model, prefixes ...string) ([]Model, error) {
	values := make([]Model, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := db.getPrefixedKey(forModel.Collection(), strings.Join(prefixes, "/"))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				modelVal := reflect.New(reflect.TypeOf(forModel).Elem()).Interface().(Model)
				if err := json.Unmarshal(val, modelVal); err != nil {
					return err
				}
				values = append(values, modelVal)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}
