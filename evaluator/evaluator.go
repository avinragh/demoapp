package evaluator

import (
	"demoapp/context"
	"demoapp/db"
)

type EvaluatorI interface {
	GetName() string
}

type AccountAlarmStateI interface {
	EvaluatorI
	Unset(aInAccountAlarmInfo *db.AccountAlarmInfo) bool
	IsSet(aInAccountAlarmInfo *db.AccountAlarmInfo) bool
}

type ServerAlarmStateI interface {
	EvaluatorI
	Unset(aInServerAlarmInfo *db.ServerAlarmInfo) bool
	IsSet(aInServerAlarmInfo *db.ServerAlarmInfo) bool
}

type AccountEvaluatorI interface {
	AccountAlarmStateI
	Evaluate(ctx *context.Context, aInAccount *db.Account, aInServers []*db.Server) (*db.AccountAlarmInfo, error)
}

type ServerEvaluatorI interface {
	ServerAlarmStateI
	Evaluate(ctx *context.Context, aInAccount *db.Account, aInServer *db.Server) (*db.ServerAlarmInfo, error)
}
