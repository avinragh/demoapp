package db

import (
	memdb "github.com/hashicorp/go-memdb"
)

type DB struct {
	*memdb.MemDB
}

func (db *DB) Init() (*DB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"account": &memdb.TableSchema{
				Name: "account",
				Indexes: map[string]*memdb.IndexSchema{
					"username": &memdb.IndexSchema{
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
			"server": &memdb.TableSchema{
				Name: "server",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
			"alarm": &memdb.TableSchema{
				Name: "alarm",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	initdb, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	db = &DB{initdb}
	return db, nil
}
