package db

import uuid "github.com/satori/go.uuid"

func (db *DB) FindServerById(id string) (*Server, error) {
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	item, err := txn.First("server", "id", id)
	if err != nil {
		return nil, err
	}
	server := item.(*Server)

	return server, nil
}

func (db *DB) FindServers(accountId *string) ([]*Server, error) {

	servers := []*Server{}
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	if accountId != nil {
		it, err := txn.Get("server", "id", "accountId", accountId)
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			server := obj.(*Server)
			servers = append(servers, server)
		}
	} else {
		it, err := txn.Get("server", "id")
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			server := obj.(*Server)
			servers = append(servers, server)
		}
	}
	return servers, nil

}

func (db *DB) AddServers(servers []*Server) ([]*Server, error) {
	for _, server := range servers {
		db.AddServer(server)
	}
	return servers, nil
}

func (db *DB) AddServer(server *Server) (*Server, error) {
	if server.Id == nil {
		uuid := uuid.NewV4().String()
		server.Id = &uuid
	}
	txn := db.MemDB.Txn(true)

	if err := txn.Insert("server", server); err != nil {
		return nil, err
	}
	txn.Commit()
	return server, nil
}

func (db *DB) DeleteServer(id string) (*Server, error) {
	server, err := db.FindServerById(id)
	if err != nil {
		return nil, err
	}
	txn := db.MemDB.Txn(true)
	txn.Delete("server", server)
	txn.Commit()
	return server, nil
}
