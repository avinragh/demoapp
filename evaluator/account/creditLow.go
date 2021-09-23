package account

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
	"demoapp/util"
)

type CreditLow struct{}

func NewCreditLow() evaluator.AccountEvaluatorI {
	return &CreditLow{}
}

func (evaluator *CreditLow) GetName() string {
	return "CreditLow"
}

func (evaluator *CreditLow) getCreditLowConfig() float64 {
	return util.GetEnvAsFloat64OrDefault("CREDIT_LOW_THRESHOLD", 10.0)
}

func (evaluator *CreditLow) Evaluate(ctx *context.Context, account *db.Account, servers []*db.Server) (*db.AccountAlarmInfo, error) {
	if account != nil {
		if account.Id != nil {
			hasChanges := false
			alarmInfo := account.AlarmInfo
			if alarmInfo == nil {
				alarmInfo = &db.AccountAlarmInfo{}
				hasChanges = true
			}
			if servers == nil {
				if *alarmInfo.IsCreditLow {
					*alarmInfo.IsCreditLow = false
				}
			} else {
				isCreditLow := *account.Credits <= evaluator.getCreditLowConfig()
				if isCreditLow == *alarmInfo.IsCreditLow {
					*alarmInfo.IsCreditLow = isCreditLow
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

func (evaluator *CreditLow) Unset(alarmInfo *db.AccountAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	if *alarmInfo.IsCreditLow {
		*alarmInfo.IsCreditLow = false
		return true
	}
	return false
}

func (evaluator *CreditLow) IsSet(alarmInfo *db.AccountAlarmInfo) bool {
	if alarmInfo == nil {
		return false
	}
	return *alarmInfo.IsCreditLow
}
