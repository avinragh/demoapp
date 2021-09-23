package registry

import "demoapp/db"

type ExecutionResult struct {
	ModifiedAccounts map[*string]*db.Account
	ModifiedServers  map[*string]*db.Server
}
