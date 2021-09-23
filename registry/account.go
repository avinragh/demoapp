package registry

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
	"demoapp/evaluator/account"
)

func RegisterAccountAlarms(ctx *context.Context) {
	Registry().AddNode(ctx, account.NewCreditLow(), nil)
	Registry().AddNode(ctx, account.NewMaxNumberOfServers(), nil)
}

func RunAcountEvaluators(ctx *context.Context, account *db.Account, servers []*db.Server) ExecutionResult {
	var err error

	hasChanges := false

	logger := ctx.GetLogger()

	var result ExecutionResult = ExecutionResult{
		ModifiedAccounts: make(map[*string]*db.Account),
		ModifiedServers:  make(map[*string]*db.Server),
	}

	alarmInfo := account.AlarmInfo

	if alarmInfo == nil {
		alarmInfo = &db.AccountAlarmInfo{}
		account.AlarmInfo = alarmInfo
		hasChanges = true
	}

	queue := make([]*node, 0)
	queue = append(queue, Registry().root.children...)

	if len(queue) > 0 {
		currentEvaluator := queue[0]
		queue = queue[1:]

		if eval, ok := interface{}(currentEvaluator.entity).(evaluator.AccountEvaluatorI); ok {
			alarmInfo, err = eval.Evaluate(ctx, account, servers)
			if err != nil {
				logger.Printf("error evaluating account alarm %s", err)
			}
			if alarmInfo != nil {
				account.AlarmInfo = alarmInfo
				hasChanges = true
			}

			alarmInfo = account.AlarmInfo

			if len(currentEvaluator.children) > 0 {
				if eval.IsSet(alarmInfo) {
					for _, childEvaluator := range Registry().GetDescendents(currentEvaluator.entity) {
						if childEval, ok := interface{}(childEvaluator).(evaluator.AccountAlarmStateI); ok {
							if childEval.Unset(alarmInfo) {
								account.AlarmInfo = alarmInfo
								hasChanges = true
							}
						}

						if childEval, ok := interface{}(childEvaluator).(evaluator.ServerAlarmStateI); ok {
							for _, server := range servers {
								serverAlarmInfo := server.AlarmInfo

								if serverAlarmInfo == nil {
									serverAlarmInfo = &db.ServerAlarmInfo{}
									server.AlarmInfo = serverAlarmInfo
								}

								if childEval.Unset(serverAlarmInfo) && server.Id != nil {
									server.AlarmInfo = serverAlarmInfo
									result.ModifiedServers[server.Id] = server
								}
							}
						}
					}
				} else {
					queue = append(queue, currentEvaluator.children...)
				}
			}
		} else {
			queue = append(queue, currentEvaluator.children...)
		}
	}
	if hasChanges {
		result.ModifiedAccounts[account.Id] = account
	}
	return result
}
