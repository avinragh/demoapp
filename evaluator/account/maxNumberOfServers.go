package account

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
	"demoapp/util"
)

type MaxNumberOfServers struct{}

func NewMaxNumberOfServers() evaluator.AccountEvaluatorI {
	return &MaxNumberOfServers{}
}

func (evaluator *MaxNumberOfServers) GetName() string {
	return "MaxNumberOfServers"
}

func (evaluator *MaxNumberOfServers) getMaxNumberOfServersConfig() int {
	return util.GetEnvAsIntOrDefault("MAX_NUMBER_OF_SERVERS", 1)
}

func (evaluator *MaxNumberOfServers) Evaluate(ctx *context.Context, account *db.Account, servers []*db.Server) (*db.AccountAlarmInfo, error) {
	if account != nil {
		if account.Id != nil {
			hasChanges := false
			alarmInfo := account.AlarmInfo
			if alarmInfo == nil {
				alarmInfo = &db.AccountAlarmInfo{}
				hasChanges = true
			}
			if servers == nil {
				if *alarmInfo.IsMaxNumberOfServers {
					*alarmInfo.IsMaxNumberOfServers = false
				}
			} else {
				isMaxNumberOfServers := len(servers) <= evaluator.getMaxNumberOfServersConfig()
				if isMaxNumberOfServers == *alarmInfo.IsMaxNumberOfServers {
					*alarmInfo.IsMaxNumberOfServers = isMaxNumberOfServers
					hasChanges = true
				}
			}
			if hasChanges {
				return alarmInfo, nil
			}
		}
	}
	return nil, nil
}

func (evaluator *MaxNumberOfServers) Unset(alarmInfo *db.AccountAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	if *alarmInfo.IsMaxNumberOfServers {
		*alarmInfo.IsMaxNumberOfServers = false
		return true
	}
	return false
}

func (evaluator *MaxNumberOfServers) IsSet(alarmInfo *db.AccountAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	return *alarmInfo.IsMaxNumberOfServers
}
