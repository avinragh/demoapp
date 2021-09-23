package account

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
)

type PoweredOff struct{}

func NewPoweredOff() evaluator.ServerEvaluatorI {
	return &PoweredOff{}
}

func (evaluator *PoweredOff) GetName() string {
	return "PoweredOff"
}

func (evaluator *PoweredOff) Evaluate(ctx *context.Context, account *db.Account, server *db.Server) (*db.ServerAlarmInfo, error) {
	if account != nil {
		if account.Id != nil {
			hasChanges := false
			alarmInfo := server.AlarmInfo
			if alarmInfo == nil {
				alarmInfo = &db.ServerAlarmInfo{}
				hasChanges = true
			}
			if server == nil {
				if *alarmInfo.IsPoweredOff {
					*alarmInfo.IsPoweredOff = false
				}
			} else {
				isPoweredOff := *server.State == "stopped"
				if isPoweredOff == *alarmInfo.IsPoweredOff {
					*alarmInfo.IsPoweredOff = isPoweredOff
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

func (evaluator *PoweredOff) Unset(alarmInfo *db.ServerAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	if *alarmInfo.IsPoweredOff {
		*alarmInfo.IsPoweredOff = false
		return true
	}
	return false
}

func (evaluator *PoweredOff) IsSet(alarmInfo *db.ServerAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	return *alarmInfo.IsPoweredOff
}
