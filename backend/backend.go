package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/datastore"
	"github.com/valsgaard/interview-case/backend/endpoints"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

const (
	serviceName   = "8th wonder"
	listenPort    = 8000
	postgresqlURI = "postgresql://tester:flaf@localhost/sybo"
)

func main() {
	// Logging
	log := common.NewLog(serviceName, os.Stdout)
	log.Info("Initiating ...")

	// Database connection
	store, err := datastore.NewSQLConnection(postgresqlURI, 5)
	if err != nil {
		log.Fatal(err)
	}

	/**************************************************************************
	***************************************************************************
	**                                                                       **
	**   Prepare Context                                                     **
	**                                                                       **
	***************************************************************************
	**************************************************************************/

	ctx := context.Background()
	ctx = port.SetDatastore(ctx, store)
	ctx = common.SetLog(ctx, log)

	/**************************************************************************
	***************************************************************************
	**                                                                       **
	**    Prepare routing                                                    **
	**                                                                       **
	***************************************************************************
	**************************************************************************/

	r := mux.NewRouter()

	// User
	r.Handle("/user", common.NewHandlerFunc(ctx, endpoints.NewUserCreate)).
		Methods("POST")

	r.Handle("/user", common.NewHandlerFunc(ctx, endpoints.NewUserGet)).
		Methods("GET")

	// Game State
	r.Handle("/user/{id}/state", common.NewHandlerFunc(ctx, endpoints.NewGameStateGet)).
		Methods("GET")

	r.Handle("/user/{id}/state", common.NewHandlerFunc(ctx, endpoints.NewGameStateUpdate)).
		Methods("PUT")

	// Friends
	r.Handle("/user/{id}/friends", common.NewHandlerFunc(ctx, endpoints.NewFriendsGet)).
		Methods("GET")

	r.Handle("/user/{id}/friends", common.NewHandlerFunc(ctx, endpoints.NewFriendsUpdate)).
		Methods("PUT")

	/**************************************************************************
	***************************************************************************
	**                                                                       **
	**    Open listeners                                                     **
	**                                                                       **
	***************************************************************************
	**************************************************************************/

	// HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", listenPort),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverShutdown := common.ListenAndServe(ctx, server, 5*time.Second)
	log.Infof("Listening on port %d", 8000)

	/**************************************************************************
	***************************************************************************
	**                                                                       **
	**    Cleanup                                                            **
	**                                                                       **
	***************************************************************************
	**************************************************************************/

	<-common.SystemTermination()

	///////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////

	log.Info("Closing HTTP Listener connections ...")

	serverShutdown()

	log.Info("Closing service")
}
