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

// FindServerByid operation middleware
func (siw *ServerInterfaceWrapper) FindServerById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Find Server By id")

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	database := ctx.GetDB()

	server, err := database.FindServerById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(server); err != nil {
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

// AddServer operation middleware
func (siw *ServerInterfaceWrapper) AddServers(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("InAddServes")

	servers := []*db.Server{}

	database := ctx.GetDB()

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &servers)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	servers, err = database.AddServers(servers)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(servers); err != nil {
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

// DeleteServer operation middleware
func (siw *ServerInterfaceWrapper) DeleteServer(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Delete Server")

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	database := ctx.GetDB()

	server, err := database.DeleteServer(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(server); err != nil {
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

// FindServerById operation middleware
func (siw *ServerInterfaceWrapper) FindServers(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Find Servwr By Id")

	database := ctx.GetDB()

	servers := []*db.Server{}

	var err error

	var params db.FindServersParams

	// ------------- Optional query parameter "username" -------------
	if paramValue := r.URL.Query().Get("accountId"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "accountId", r.URL.Query(), &params.AccountId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter username: %s", err), http.StatusBadRequest)
		return
	}

	servers, err = database.FindServers(params.AccountId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter username: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(servers); err != nil {
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

// AddServerById operation middleware
func (siw *ServerInterfaceWrapper) AddServerById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Add Server By ID")

	server := &db.Server{}

	database := ctx.GetDB()

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &server)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	server, err = database.AddServer(server)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(server); err != nil {
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
