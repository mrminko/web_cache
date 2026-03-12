package main

type UnixTime int64

type Database map[string]*Cache

type Cache struct {
	obj          []byte
	lastModified UnixTime
}

func (db Database) Add(key string, value *Cache) {
	db[key] = value
}

func (db Database) Update(key string, value *Cache) {
	delete(db, key)
	db.Add(key, value)
}

func (db Database) Has(key string) (*Cache, bool) {
	v, ok := db[key]
	if !ok {
		return nil, false
	}
	return v, true
}
