package main

import (
	"demoapp/context"
	"demoapp/cron"
	"demoapp/web"
	"net/http"
)

func main() {
	ctx := &context.Context{}
	ctx = ctx.Init()
	go cron.ScheduleUpdateResources(ctx)
	go cron.ScheduleUpdateAlarms(ctx)
	go cron.ScheduleExecuteAlarms(ctx)
	siw := &web.ServerInterfaceWrapper{}
	r := web.Handler(ctx, siw.Handler)
	http.ListenAndServe(":8080", r)

}
