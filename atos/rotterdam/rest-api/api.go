//
// Copyright 2018 Atos
//
// ROTTERDAM application
// CLASS Project: https://class-project.eu/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @author: ATOS
//

package rest_api

import (
	adaptation_engine "atos/rotterdam/adaptation-engine"
	clustersMonitoring "atos/rotterdam/adaptation-engine/monitoring"
	caas "atos/rotterdam/caas"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	db "atos/rotterdam/database/caas"
	faas "atos/rotterdam/faas"
	imec "atos/rotterdam/imec"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// path used in logs
const pathLOG string = "Rotterdam > REST API : "

/*
InitializeRESTAPI initialization function
*/
func InitializeRESTAPI() {
	log.Println(pathLOG + "[InitializeRESTAPI] Initializing REST API' ...")

	// initialize adapter after getting 'Mode' value
	caas.InitializeAdapter()

	// set paths / routes
	router := mux.NewRouter()

	// routes:
	// default
	router.HandleFunc("/", caas.HomePath).Methods("GET")

	// configuration
	router.HandleFunc("/api/", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/", caas.NotImplementedFunc).Methods("GET")
	// get configuration from K8s
	router.HandleFunc("/api/v1/version", caas.GetVersion).Methods("GET")
	router.HandleFunc("/api/v1/status", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/caas/config", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/rules-engine/config", adaptation_engine.NotImplementedFunc).Methods("GET")
	// configuration & infrastructures
	router.HandleFunc("/api/v1/config", caas.UpdateConfig).Methods("PUT")
	router.HandleFunc("/api/v1/config", caas.GetCurrentConfig).Methods("GET")

	// deployment and provisioning functions:
	// get all tasks
	router.HandleFunc("/api/v1/docks/tasks", caas.GetAllTasks).Methods("GET")
	// get all tasks qos - DEPRECATED
	router.HandleFunc("/api/v1/docks/tasksqos", caas.GetAllTasksQoS).Methods("GET")
	// deploy task - DEPRECATED
	router.HandleFunc("/api/v1/docks/tasks", DeployTask).Methods("POST")
	// deploy COMPSs task - DEPRECATED
	router.HandleFunc("/api/v1/docks/tasks-compss", DeployTaskCOMPSs).Methods("POST")
	// not-implemented - DEPRECATED
	router.HandleFunc("/api/v1/docks/{dock}/tasks", caas.GetDockTasks).Methods("GET")
	// get task - DEPRECATED
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{id}", caas.GetTask).Methods("GET")
	// not-implemented - DEPRECATED
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{id}", caas.NotImplementedFunc).Methods("PUT")
	// remove task - DEPRECATED
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{id}", caas.RemoveTask).Methods("DELETE")
	// not-implemented
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{name}/containers", caas.NotImplementedFunc).Methods("GET")

	/************************************
	 rules engine functions
	 Path: /api/v1/sla...
	************************************/
	// sla - violation
	router.HandleFunc("/api/v1/sla/tasks/{id}/guarantee/{guarantee}", adaptation_engine.ProcessViolation).Methods("POST")

	/************************************
	 Path: /api/v1/qos/definitions...
	************************************/
	// post qos definitions for SLAs creation
	router.HandleFunc("/api/v1/qos/definitions", caas.LoadQoSDefinitions).Methods("POST")
	// get qos definitions for SLAs creation
	router.HandleFunc("/api/v1/qos/definitions", caas.GetQoSDefinitions).Methods("GET")
	// get qos definition by name
	router.HandleFunc("/api/v1/qos/definitions/{name}", caas.GetQoSDefinition).Methods("GET")

	/************************************
	 Path: /api/v1/tasks...
	************************************/
	// get all tasks
	router.HandleFunc("/api/v1/tasks", caas.GetAllTasks).Methods("GET")
	// deploy task
	router.HandleFunc("/api/v1/tasks", Deploy).Methods("POST")
	// get task
	router.HandleFunc("/api/v1/tasks/{id}", caas.GetTask).Methods("GET")
	// get task
	router.HandleFunc("/api/v1/tasks/{id}/all", caas.GetTaskAllInfo).Methods("GET")
	// remove task
	router.HandleFunc("/api/v1/tasks/{id}", Remove).Methods("DELETE")

	/************************************
	 Path: /api/v1/functions...
	************************************/
	// get all functions
	router.HandleFunc("/api/v1/functions", faas.GetAllFunctions).Methods("GET")
	// deploy function
	router.HandleFunc("/api/v1/functions", Deploy).Methods("POST")
	// call function
	router.HandleFunc("/api/v1/functions/{id}", faas.CallFunction).Methods("POST")
	// get function
	router.HandleFunc("/api/v1/functions/{id}", faas.GetFunction).Methods("GET")
	// remove function
	router.HandleFunc("/api/v1/functions/{id}", Remove).Methods("DELETE")

	/************************************
	 imec functions
	 Path: /api/v1/imec...
	************************************/
	// INFRASTRUCTURES:
	// get all available and running infrastructures
	router.HandleFunc("/api/v1/imec", imec.GetAllInfrastructures).Methods("GET")
	// create new infrastructure
	router.HandleFunc("/api/v1/imec", imec.CreateInfrastructure).Methods("POST")
	// update infrastrucutre
	router.HandleFunc("/api/v1/imec/{id}", imec.UpdateInfrastructure).Methods("PUT")
	// get infrastructure
	router.HandleFunc("/api/v1/imec/{id}", imec.GetInfrastructure).Methods("GET")
	// delete infrastrucutre
	router.HandleFunc("/api/v1/imec/{id}", imec.DeleteInfrastructure).Methods("DELETE")
	// CLUSTERS:
	// deploy a cluster in an exisiting infrastructure
	router.HandleFunc("/api/v1/imec/{id}/cluster", imec.DeployCluster).Methods("POST")
	// get cluster info
	router.HandleFunc("/api/v1/imec/{id}/cluster", imec.GetCluster).Methods("GET")
	// delete cluster
	router.HandleFunc("/api/v1/imec/{id}/cluster", imec.DeleteCluster).Methods("DELETE")

	// get cluster info
	router.HandleFunc("/api/v1/imec/{id}/status", clustersMonitoring.GetClusterStatusInfo).Methods("GET")

	// swagger api
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./rest-api/swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	/////////////////////////////////////////////////////////////////
	// start server:
	server := &http.Server{Addr: ":8333", Handler: router}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		// init DB
		db.InitDB()

		// init DB
		log.Println(pathLOG + "[InitializeRESTAPI] Server listening on port " + strconv.Itoa(cfg.Config.ServerPort) + " ...")
		log.Printf("................................................................................")
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop

	/////////////////////////////////////////////////////////////////
	// stop server:
	log.Printf("................................................................................")
	// close DB
	db.CloseDB()

	// stop cluster monitoring process
	clustersMonitoring.StopCheckingClusters()

	// shutdown server
	log.Println(pathLOG + "[InitializeRESTAPI] Shutting down server ...")

	var shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second, "shutdown timeout (5s,5m,5h) before connections are cancelled")
	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println(pathLOG + "[InitializeRESTAPI] Terminated")
	log.Printf("................................................................................")
}
