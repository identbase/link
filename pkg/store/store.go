// Package store provides an in-memory key-value store in place of more
// permenant data store but will provide the intial rails to achieve our data
// storage requirements.
package store

import (
	"strings"
	"sync"
)

/*
Log Levels */
const (
	DEBUG   string = "DEBUG"
	INFO    string = "INFO"
	WARNING string = "WARN"
	ERROR   string = "ERROR"
	FATAL   string = "FATAL"
)

/*
Keyer is an object that can be kept in an InMemory store */
type Keyer interface {
	Key() string
}

/*
Logger is a standard logging interface that we can use to write information
about what the in-memory store is doing. */
type Logger interface {
	Print(i ...interface{})
	Printf(format string, args ...interface{})
	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Warn(i ...interface{})
	Warnf(format string, args ...interface{})
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(i ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(i ...interface{})
	Panicf(format string, args ...interface{})
}

/*
InMemory store is thread-safe and handles any object implemented as a Keyer. */
type InMemory struct {
	log Logger
	sync.RWMutex
	data map[string]Keyer
}

/*
New creates a new InMemory store. */
func New(l Logger) *InMemory {
	return &InMemory{
		log: l,
	}
}

/*
Put stores a value. */
func (db *InMemory) Put(v Keyer) error {
	db.Lock()
	defer db.Unlock()

	if db.data == nil {
		db.data = make(map[string]Keyer)
	}

	db.log.Debugf("Storing object with key %s", v.Key())

	db.data[v.Key()] = v

	return nil
}

/*
Get retrieves a value or panics if the value does not exist. */
func (db *InMemory) Get(k string) Keyer {
	db.RLock()
	defer db.RUnlock()

	db.log.Debugf("Getting stored object %s", k)

	return db.data[k]
}

/*
Lookup retrieves a value or returns nil if it doesnt exist. */
func (db *InMemory) Lookup(k string) *Keyer {
	db.RLock()
	defer db.RUnlock()

	db.log.Debugf("Getting stored object %s", k)
	for key, value := range db.data {
		if key == k {
			return &value
		}
	}
	return nil

}

/*
Remove a value */
func (db *InMemory) Remove(k string) error {
	db.Lock()
	defer db.Unlock()

	if db.data == nil {
		db.data = make(map[string]Keyer)
	}

	db.log.Debugf("Deleting stored object %s", k)

	delete(db.data, k)

	return nil
}

/*
Search on keys */
func (db *InMemory) Search(q string) ([]Keyer, error) {
	var result []Keyer

	db.log.Debugf("Searching for stored object with key %s ...", q)

	for key := range db.data {
		if strings.Contains(key, q) {
			value := db.Get(key)
			result = append(result, value)
		}
	}

	return result, nil
}
