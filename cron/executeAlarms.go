package cron

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/registry"
	"demoapp/util"
	"sync"
	"time"
)

const (
	ExecuteAlarmInterval = "UPDATE_ALARM_INTERVAL"
)

func ScheduleExecuteAlarms(ctx *context.Context) {
	var wg sync.WaitGroup
	logger := ctx.GetLogger()
	executeAlarmsInterval := util.GetEnvAsIntOrDefault(ExecuteAlarmInterval, 2)
	wg.Add(1)
	go func(ctx *context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		c := time.Tick(time.Duration(executeAlarmsInterval) * time.Minute)
		for range c {
			err := ExecuteAlarms(ctx)
			if err != nil {
				logger.Printf("Error scheduling UpdateAlarms: %s", err)
			}
		}
	}(ctx, &wg)
	wg.Wait()

}
func ExecuteAlarms(ctx *context.Context) error {
	database := ctx.GetDB()
	accounts, err := database.FindAccounts(nil)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		changeAccountMap := make(map[*string]*db.Account)
		changeServerMap := make(map[*string]*db.Server)

		for _, account := range accounts {
			if account != nil && account.Id != nil {
				servers, err := database.FindServers(account.Id, nil)
				if err != nil {
					return err
				}
				for _, server := range servers {
					if server.Id != nil {
						result := registry.RunServerEvaluators(ctx, account, server)
						for k, v := range result.ModifiedAccounts {
							changeAccountMap[k] = v
						}
						for k, v := range result.ModifiedServers {
							changeServerMap[k] = v
						}
					}
				}

				accountResult := registry.RunAcountEvaluators(ctx, account, servers)
				for k, v := range accountResult.ModifiedAccounts {
					changeAccountMap[k] = v
				}
				for k, v := range accountResult.ModifiedServers {
					changeServerMap[k] = v
				}

			}
		}
		if len(changeAccountMap) > 0 {
			for _, v := range changeAccountMap {
				_, err = database.AddAccount(v)
				if err != nil {
					return err
				}
			}
		}

		if len(changeServerMap) > 0 {
			for _, v := range changeServerMap {
				_, err = database.AddServer(v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
