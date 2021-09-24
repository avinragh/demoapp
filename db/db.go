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
					"queryIndex": &memdb.IndexSchema{
						Name:    "queryIndex",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "AccountId"},
					},
					"uuid": &memdb.IndexSchema{
						Name:    "uuid",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Uuid"},
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
					"alarmType": &memdb.IndexSchema{
						Name:    "alarmType",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "AlarmType"},
					},
					"resourceId": &memdb.IndexSchema{
						Name:    "resourceId",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ResourceId"},
					},
					"name": &memdb.IndexSchema{
						Name:    "name",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
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
