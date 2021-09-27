package cron

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/util"
	"sync"
	"time"
)

const (
	ServerAlarmType             = "server"
	AccountAlarmType            = "account"
	CreditLowAlarmName          = "creditLow"
	MaxNumberOfServersAlarmName = "maxNumberOfServers"
	IsPoweredOffAlarmName       = "isPoweredOff"
	IsErroredAlarmName          = "isErrored"
	UpdateAlarmInterval         = "UPDATE_ALARM_INTERVAL"
)

func ScheduleUpdateAlarms(ctx *context.Context) {
	var wg sync.WaitGroup
	logger := ctx.GetLogger()
	updateAlarmsInterval := util.GetEnvAsIntOrDefault(UpdateAlarmInterval, 3)
	wg.Add(1)
	go func(ctx *context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		c := time.Tick(time.Duration(updateAlarmsInterval) * time.Minute)
		for range c {
			err := UpdateAlarms(ctx)
			if err != nil {
				logger.Printf("Error scheduling UpdateAlarms: %s", err)
			}
		}
	}(ctx, &wg)
	wg.Wait()
}

func UpdateAlarms(ctx *context.Context) error {
	database := ctx.GetDB()
	accounts, err := database.FindAccounts(nil)
	if err != nil {
		return err
	}

	accountAlarmType := AccountAlarmType
	// ServerAlarmType := ServerAlarmType

	addAlarms := []*db.Alarm{}
	for _, account := range accounts {
		if account.Id != nil {
			if account.AlarmInfo != nil {
				if account.AlarmInfo.IsCreditLow != nil {
					creditLowAlarmName := CreditLowAlarmName
					existingAlarms, err := database.FindAlarms(&accountAlarmType, account.Id, &creditLowAlarmName)
					if err != nil {
						return err
					}
					if *account.AlarmInfo.IsCreditLow {
						if len(existingAlarms) == 0 || existingAlarms == nil {
							alarm := &db.Alarm{}
							alarm.Name = CreditLowAlarmName
							alarm.AlarmType = AccountAlarmType
							alarm.ResourceId = *account.Id
							addAlarms = append(addAlarms, alarm)

						}
					} else {
						if len(existingAlarms) > 0 && existingAlarms != nil {
							delId := *existingAlarms[0].Id
							if delId != "" {
								_, err := database.DeleteAlarm(delId)
								if err != nil {
									return err
								}
							}

						}

					}
				}
			}

		}
	}
	_, err = database.AddAlarms(addAlarms)
	if err != nil {
		return err
	}
	return nil
}
