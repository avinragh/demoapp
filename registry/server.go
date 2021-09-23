package registry

import (
	"demoapp/context"
	"demoapp/db"
	"demoapp/evaluator"
	serverEval "demoapp/evaluator/server"
)

var serverPoweredOffAlarm = serverEval.NewPoweredOff()

func RegisterDeploymentDeviceAlarms(ctx *context.Context) {
	Registry().AddNode(ctx, serverEval.NewPoweredOff(), nil)
	Registry().AddNode(ctx, serverEval.NewErrored(), nil)

}

func RunServerEvaluators(ctx *context.Context, account *db.Account, server *db.Server) ExecutionResult {
	logger := ctx.GetLogger()
	eResult := ExecutionResult{
		ModifiedAccounts: make(map[*string]*db.Account),
		ModifiedServers:  make(map[*string]*db.Server),
	}

	alarmInfo := account.AlarmInfo
	serverAlarmInfo := server.AlarmInfo

	if alarmInfo == nil {
		alarmInfo = &db.AccountAlarmInfo{}
		account.AlarmInfo = alarmInfo
	}
	if serverAlarmInfo == nil {
		serverAlarmInfo = &db.ServerAlarmInfo{}
		server.AlarmInfo = serverAlarmInfo
	}

	var err error

	queue := make([]*node, 0)
	queue = append(queue, Registry().root.children...)
	for len(queue) > 0 {
		currentEvaluator := queue[0]
		queue = queue[1:]

		if eval, ok := interface{}(currentEvaluator.entity).(evaluator.ServerEvaluatorI); ok {
			var result *db.ServerAlarmInfo
			result, err = eval.Evaluate(ctx, account, server)
			if err != nil {
				logger.Println("error evaluating alarm")
			}
			if result != nil && server.Id != nil {
				server.AlarmInfo = result
				serverAlarmInfo = result
				eResult.ModifiedServers[server.Id] = server
			}
			if len(currentEvaluator.children) > 0 {
				if eval.IsSet(serverAlarmInfo) {
					for _, childEvaluator := range Registry().GetDescendents(currentEvaluator.entity) {
						if serverVal, ok := interface{}(childEvaluator).(evaluator.ServerEvaluatorI); ok {
							if serverVal.Unset(serverAlarmInfo) {
								server.AlarmInfo = serverAlarmInfo
								eResult.ModifiedServers[server.Id] = server
							} else if accountVal, ok := interface{}(childEvaluator).(evaluator.AccountEvaluatorI); ok {
								if accountVal.Unset(alarmInfo) {
									account.AlarmInfo = alarmInfo
									eResult.ModifiedAccounts[account.Id] = account
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
	return eResult
}
