package account

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
)

type Errored struct{}

func NewErrored() evaluator.ServerEvaluatorI {
	return &Errored{}
}

func (evaluator *Errored) GetName() string {
	return "Errored"
}

func (evaluator *Errored) Evaluate(ctx *context.Context, account *db.Account, server *db.Server) (*db.ServerAlarmInfo, error) {
	if account != nil {
		if account.Id != nil {
			hasChanges := false
			alarmInfo := server.AlarmInfo
			if alarmInfo == nil {
				alarmInfo = &db.ServerAlarmInfo{}
				hasChanges = true
			}
			if server == nil {
				if *alarmInfo.IsErrored {
					*alarmInfo.IsErrored = false
				}
			} else {
				isErrored := *server.State == "error"
				if isErrored == *alarmInfo.IsErrored {
					*alarmInfo.IsErrored = isErrored
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

func (evaluator *Errored) Unset(alarmInfo *db.ServerAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	if *alarmInfo.IsErrored {
		*alarmInfo.IsErrored = false
		return true
	}
	return false
}

func (evaluator *Errored) IsSet(alarmInfo *db.ServerAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	return *alarmInfo.IsErrored
}
