package database

import (
	"encoding/json"
	"fmt"
	"io"
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

func (db *Database) toEntry(m Model) (*badger.Entry, error) {
	entryValue, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return badger.NewEntry(db.getPrefixedKey(m.Collection(), m.Key()), entryValue), nil
}

func (db *Database) Put(models ...Model) error {
	return db.db.Update(func(txn *badger.Txn) error {
		for _, m := range models {
			e, err := db.toEntry(m)
			if err != nil {
				return err
			}
			if err := txn.SetEntry(e); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *Database) Get(key string, val Model) error {
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.getPrefixedKey(val.Collection(), key))
		if err != nil {
			return err
		}
		return item.Value(func(value []byte) error {
			return json.Unmarshal(value, val)
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func Get[T Model](db *Database, key string) (T, error) {
	var val T
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.getPrefixedKey(val.Collection(), key))
		if err != nil {
			return err
		}
		return item.Value(func(value []byte) error {
			return json.Unmarshal(value, &val)
		})
	})
	if err != nil {
		return val, err
	}
	return val, nil
}

func (db *Database) RawGet(collection, key string) ([]byte, error) {
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

func Values[T Model](db *Database, prefixes ...string) ([]T, error) {
	values := make([]T, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		var collectionType T
		prefix := db.getPrefixedKey(collectionType.Collection(), strings.Join(prefixes, "/"))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var modelVal T
				if err := json.Unmarshal(val, &modelVal); err != nil {
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

func (db *Database) RawValues(w io.Writer, prefixes ...string) error {
	return db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		_, err := w.Write([]byte("["))
		if err != nil {
			return err
		}
		prefix := []byte(strings.Join(prefixes, "/"))
		for it.Seek(prefix); it.ValidForPrefix(prefix); {
			item := it.Item()
			err = item.Value(func(val []byte) error {
				_, err = w.Write(val)
				return err
			})
			if err != nil {
				return err
			}
			it.Next()
			if it.ValidForPrefix(prefix) {
				_, err = w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte("]"))
		return err
	})
}

func (db *Database) Close() error {
	return db.db.Close()
}
