package web

import (
	"demoapp/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// FindAlarmByid operation middleware
func (siw *ServerInterfaceWrapper) FindAlarmById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Find Alarm By id")

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	database := ctx.GetDB()

	alarm, err := database.FindAlarmById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(alarm); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}

	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}

// AddAlarm operation middleware
func (siw *ServerInterfaceWrapper) AddAlarms(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("InAddServes")

	alarms := []*db.Alarm{}

	database := ctx.GetDB()

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &alarms)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alarms, err = database.AddAlarms(alarms)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(alarms); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}

// DeleteAlarm operation middleware
func (siw *ServerInterfaceWrapper) DeleteAlarm(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Delete Alarm")

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	database := ctx.GetDB()

	alarm, err := database.DeleteAlarm(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(alarm); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}

// FindAlarmById operation middleware
func (siw *ServerInterfaceWrapper) FindAlarms(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Find Servwr By Id")

	database := ctx.GetDB()

	alarms := []*db.Alarm{}

	var err error

	var params db.FindAlarmsParams

	// ------------- Optional query parameter "username" -------------
	if paramValue := r.URL.Query().Get("alarmType"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "alarmType", r.URL.Query(), &params.AlarmType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter alarmType: %s", err), http.StatusBadRequest)
		return
	}

	if paramValue := r.URL.Query().Get("resourceId"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "resourceId", r.URL.Query(), &params.ResourceId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter resourceId: %s", err), http.StatusBadRequest)
		return
	}
	if paramValue := r.URL.Query().Get("resourceId"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	alarms, err = database.FindAlarms(params.AlarmType, params.ResourceId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get Alarms: %s", err), http.StatusBadRequest)
		return

	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(alarms); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}

// AddAlarmById operation middleware
func (siw *ServerInterfaceWrapper) AddAlarmById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Add Alarm By ID")

	alarm := &db.Alarm{}

	database := ctx.GetDB()

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &alarm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alarm, err = database.AddAlarm(alarm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(alarm); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}
