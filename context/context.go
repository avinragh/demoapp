package context

import (
	"demoapp/db"
	"log"
	"os"
)

type Context struct {
	DB     *db.DB
	Logger *log.Logger
}

func (ctx *Context) Init() *Context {
	var err error
	db := &db.DB{}
	db, err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
	ctx.DB = db
	file, err := os.OpenFile("demoapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Openfile error %s", err)
	}
	logger := &log.Logger{}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ctx.Logger = logger
	return ctx
}

func (ctx *Context) GetDB() (db *db.DB) {
	db = ctx.DB
	return
}

func (ctx *Context) GetLogger() (logger *log.Logger) {
	logger = ctx.Logger
	return
}
