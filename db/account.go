package db

import uuid "github.com/satori/go.uuid"

func (db *DB) FindAccountById(id string) (*Account, error) {
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	item, err := txn.First("account", "id", id)
	if err != nil {
		return nil, err
	}
	account := item.(*Account)

	return account, nil
}

func (db *DB) FindAccounts(username *string) ([]*Account, error) {

	accounts := []*Account{}
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	if username != nil {
		it, err := txn.Get("account", "username", *username)
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			account := obj.(*Account)
			accounts = append(accounts, account)
		}

	} else {
		it, err := txn.Get("account", "username")
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			account := obj.(*Account)
			accounts = append(accounts, account)
		}
	}
	return accounts, nil

}

func (db *DB) AddAccounts(accounts []*Account) ([]*Account, error) {
	for _, account := range accounts {
		db.AddAccount(account)
	}
	return accounts, nil
}

func (db *DB) AddAccount(account *Account) (*Account, error) {
	if account.Id == nil {
		uuid := uuid.NewV4().String()
		account.Id = &uuid
	}
	txn := db.MemDB.Txn(true)

	if err := txn.Insert("account", account); err != nil {
		return nil, err
	}
	txn.Commit()
	return account, nil
}

func (db *DB) DeleteAccount(id string) (*Account, error) {
	account, err := db.FindAccountById(id)
	if err != nil {
		return nil, err
	}
	txn := db.MemDB.Txn(true)
	txn.Delete("account", account)
	txn.Commit()
	return account, nil
}
