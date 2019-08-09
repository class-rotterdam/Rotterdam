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
// Created on 28 May 2019
// @author: Roi Sucasas - ATOS
//

package main

import (
	adaptation_engine "atos/rotterdam/adaptation-engine"
	caas "atos/rotterdam/caas"
	common "atos/rotterdam/caas/common"
	cfg "atos/rotterdam/config"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

/*
 * Main: launch server
 */
func main() {
	log.Printf("................................................................................")

	// set global properties
	cfg_path := "./config/config.json"
	log.Println("Rotterdam > Rest API > Getting configuration values from file [" + cfg_path + "] ...")

	cfg.Config = cfg.Configuration{}

	file, err := os.Open(cfg_path)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.Config)
	if err != nil {
		panic(err)
	}

	// check ENV values and modify global properties if specified
	cfg.InitConfig(&cfg.Config)

	// show global configuration
	log.Printf("Rotterdam > Rest API > Rest Api Version ...... %s", cfg.Config.RestApiVersion)
	log.Printf("Rotterdam > Rest API > Rules Engine Version .. %s", cfg.Config.RulesEngineVersion)
	log.Printf("Rotterdam > Rest API > CaaS Version .......... %s", cfg.Config.CaaSVersion)
	log.Printf("Rotterdam > Rest API > SLALite Version ....... %s", cfg.Config.SLALiteVersion)
	log.Printf("Rotterdam > Rest API > Server Port ........... %d", cfg.Config.ServerPort)
	log.Printf("Rotterdam > Rest API > Clusters Configurations:")
	for index, _ := range cfg.Config.Clusters {
		log.Printf("Rotterdam > Rest API >  => %s", cfg.Config.Clusters[index].Name)
		log.Printf("Rotterdam > Rest API >     Description: %s", cfg.Config.Clusters[index].Description)
		log.Printf("Rotterdam > Rest API >     Mode: %s", cfg.Config.Clusters[index].Mode)
		log.Printf("Rotterdam > Rest API >     Server IP: %s", cfg.Config.Clusters[index].ServerIP)
		log.Printf("Rotterdam > Rest API >     Kubernetes Endpoint: %s", cfg.Config.Clusters[index].KubernetesEndPoint)
		log.Printf("Rotterdam > Rest API >     Openshift Endpoint: %s", cfg.Config.Clusters[index].OpenshiftEndPoint)
		log.Printf("Rotterdam > Rest API >     SLALite Endpoint: %s", cfg.Config.Clusters[index].SLALiteEndPoint)
	}
	log.Printf("Rotterdam > Rest API > Swagger UI ............ http://localhost:8333/swaggerui/index.html")
	log.Println("Rotterdam > Rest API > Starting server [port " + strconv.Itoa(cfg.Config.ServerPort) + "] ...")

	// initialize adapter after getting 'Mode' value
	caas.InitializeAdapter()

	// set paths / routes
	router := mux.NewRouter()

	// routes:
	// default
	router.HandleFunc("/", caas.HomePath).Methods("GET")

	// tests
	router.HandleFunc("/api/v1/test/req", caas.TestGetRequest).Methods("GET")
	router.HandleFunc("/api/v1/test/req", caas.TestPostRequest).Methods("POST")

	// configuration
	router.HandleFunc("/api/", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/", caas.NotImplementedFunc).Methods("GET")
	// get configuration from K8s
	router.HandleFunc("/api/v1/config", caas.GetConfig).Methods("GET")
	router.HandleFunc("/api/v1/version", caas.GetVersion).Methods("GET")
	router.HandleFunc("/api/v1/status", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/caas/config", caas.NotImplementedFunc).Methods("GET")
	router.HandleFunc("/api/v1/rules-engine/config", adaptation_engine.NotImplementedFunc).Methods("GET")

	// rules engine functions
	// sla - violation
	router.HandleFunc("/api/v1/sla/tasks/{name}/guarantee/{guarantee}", adaptation_engine.ProcessViolation).Methods("POST")

	// deployment and provisioning functions:
	// get all tasks
	router.HandleFunc("/api/v1/docks/tasks", caas.GetAllTasks).Methods("GET")
	// get all tasks qos
	router.HandleFunc("/api/v1/docks/tasksqos", caas.GetAllTasksQoS).Methods("GET")
	// post qos definitions for SLAs creation
	router.HandleFunc("/api/v1/qos/definitions", caas.LoadQoSDefinitions).Methods("POST")
	// get qos definitions for SLAs creation
	router.HandleFunc("/api/v1/qos/definitions", caas.GetQoSDefinitions).Methods("GET")
	// get qos definition by name
	router.HandleFunc("/api/v1/qos/definitions/{name}", caas.GetQoSDefinition).Methods("GET")
	// deploy task
	router.HandleFunc("/api/v1/docks/tasks", caas.DeployTask).Methods("POST")
	// deploy COMPSs task
	router.HandleFunc("/api/v1/docks/tasks-compss", caas.DeployTaskCOMPSs).Methods("POST")
	// not-implemented
	router.HandleFunc("/api/v1/docks/{dock}/tasks", caas.GetDockTasks).Methods("GET")
	// get task
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{name}", caas.GetTask).Methods("GET")
	// not-implemented
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{name}", caas.NotImplementedFunc).Methods("PUT")
	// remove task
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{name}", caas.RemoveTask).Methods("DELETE")
	// not-implemented
	router.HandleFunc("/api/v1/docks/{dock}/tasks/{name}/containers", caas.NotImplementedFunc).Methods("GET")

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
		common.InitDB()

		// init DB
		log.Println("Rotterdam > Rest API > Server listening on port " + strconv.Itoa(cfg.Config.ServerPort) + " ...")
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
	common.CloseDB()

	// shutdown server
	log.Println("Rotterdam > Rest API > Shutting down server ...")

	var shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second, "shutdown timeout (5s,5m,5h) before connections are cancelled")
	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Rotterdam > Rest API > Terminated")
	log.Printf("................................................................................")
}
