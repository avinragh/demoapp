package cron

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/upcloud"
	"demoapp/util"
	"sync"
	"time"
)

const (
	UpdateResourcesInterval = "UpdateResourcesInterval"
)

func ScheduleUpdateResources(ctx *context.Context) {
	var wg sync.WaitGroup
	logger := ctx.GetLogger()
	updateResourcesInterval := util.GetEnvAsIntOrDefault(UpdateResourcesInterval, 1)

	wg.Add(1)
	go func(ctx *context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		c := time.Tick(time.Duration(updateResourcesInterval) * time.Minute)
		for range c {
			err := UpdateResources(ctx)
			if err != nil {
				logger.Printf("Error scheduling UpdateAlarms: %s", err)
			}
		}
	}(ctx, &wg)
	wg.Wait()
	// To stop the task
}

func UpdateResources(ctx *context.Context) error {

	database := ctx.GetDB()

	accounts, err := database.FindAccounts(nil)
	if err != nil {
		return err
	}

	for _, account := range accounts {

		account, err = upcloud.GetAccount(ctx, account)
		if err != nil {
			return err
		}
		account, err = database.AddAccount(account)
		if err != nil {
			return err
		}
		servers, err := upcloud.GetServers(account)
		if err != nil {
			return err
		}

		dbServers, err := database.FindServers(account.Id)
		if err != nil {
			return err
		}

		updatedServers := []*db.Server{}
		for _, server := range servers {
			for _, dbServer := range dbServers {
				updatedServer := server
				if server.Uuid == dbServer.Uuid {
					updatedServer.Id = dbServer.Id
					updatedServers = append(updatedServers, updatedServer)
				}
			}
		}

		_, err = database.AddServers(updatedServers)
		if err != nil {
			return err
		}

	}

	return nil
}
