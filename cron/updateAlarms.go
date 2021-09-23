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
	ServerAlarmType := ServerAlarmType

	addAlarms := []*db.Alarm{}
	for _, account := range accounts {
		if account.AlarmInfo != nil {
			accountAlarms, err := database.FindAlarms(&accountAlarmType, account.Id)
			if err != nil {
				return err
			}

			if len(accountAlarms) == 0 {
				if *account.AlarmInfo.IsCreditLow {
					alarm := &db.Alarm{}
					alarm.Name = CreditLowAlarmName
					alarm.AlarmType = AccountAlarmType
					alarm.ResourceId = *account.Id
					addAlarms = append(addAlarms, alarm)
				}
				if *account.AlarmInfo.IsMaxNumberOfServers {
					alarm := &db.Alarm{}
					alarm.Name = MaxNumberOfServersAlarmName
					alarm.AlarmType = AccountAlarmType
					alarm.ResourceId = *account.Id
					addAlarms = append(addAlarms, alarm)
				}
			} else {
				if !*account.AlarmInfo.IsCreditLow {
					delId := *accountAlarms[0].Id
					if delId != "" {
						_, err := database.DeleteAlarm(delId)
						if err != nil {
							return err
						}
					}

				}
				if !*account.AlarmInfo.IsMaxNumberOfServers {
					delId := *accountAlarms[0].Id
					if delId != "" {
						_, err := database.DeleteAlarm(delId)
						if err != nil {
							return err
						}

					}
				}
			}
		}
		servers, err := database.FindServers(account.Id)
		if err != nil {
			return err
		}
		for _, server := range servers {
			serverAlarms, err := database.FindAlarms(&ServerAlarmType, account.Id)
			if err != nil {
				return err
			}
			if len(servers) == 0 {
				if *server.AlarmInfo.IsPoweredOff {
					alarm := &db.Alarm{}
					alarm.Name = IsPoweredOffAlarmName
					alarm.AlarmType = ServerAlarmType
					alarm.ResourceId = *server.Id
					addAlarms = append(addAlarms, alarm)
				}
				if *server.AlarmInfo.IsErrored {
					alarm := &db.Alarm{}
					alarm.Name = IsErroredAlarmName
					alarm.AlarmType = ServerAlarmType
					alarm.ResourceId = *server.Id
					addAlarms = append(addAlarms, alarm)
				}
			} else {
				if !*server.AlarmInfo.IsPoweredOff {
					delId := *serverAlarms[0].Id
					if delId != "" {
						_, err := database.DeleteAlarm(delId)
						if err != nil {
							return err
						}
					}
				}
				if !*server.AlarmInfo.IsErrored {
					delId := *serverAlarms[0].Id
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
	if len(addAlarms) > 0 {
		database.AddAlarms(addAlarms)
	}

	return nil
}
