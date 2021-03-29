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

package main

import (
	clustersMonitoring "atos/rotterdam/adaptation-engine/monitoring"
	urls "atos/rotterdam/caas/adapters"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/database/imec"
	rest_api "atos/rotterdam/rest-api"
	"encoding/json"
	"os"
	"strconv"
)

// path used in logs
const pathLOG string = "Rotterdam : "

// fileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/*
 * Main: launch server
 */
func main() {
	log.Printf("................................................................................")

	// set global properties from configuration file
	cfgPath := "./config/config.json"

	cfgFile := os.Getenv("ConfigPath")
	if cfgFile != "" {
		cfgPath = cfgFile
	}

	// arguments: [1] ... alternative configuration file
	if len(os.Args) == 2 {
		log.Println(pathLOG + "Arguments[1] = " + os.Args[1])
		if fileExists(os.Args[1]) {
			log.Println(pathLOG + "Using configuration from " + os.Args[1])
			cfgPath = os.Args[1]
		}
	} else {
		log.Println(pathLOG + "Using default configuration file")
	}
	log.Println(pathLOG + "Getting configuration values from file [" + cfgPath + "] ...")

	cfg.Config = cfg.Configuration{}

	file, err := os.Open(cfgPath)
	if err != nil {
		log.Error(pathLOG+"ERROR (1) Error opening configuration file. ", err)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.Config)
	if err != nil {
		log.Error(pathLOG+"ERROR (2) Error parsing configuration file conent (JSON). ", err)
		panic(err)
	}

	// check ENV values and modify global properties if specified
	cfg.InitConfig(&cfg.Config)

	// show global configuration
	log.Printf(pathLOG+"Rest Api Version ...... %s", cfg.Config.RestApiVersion)
	log.Printf(pathLOG+"Rules Engine Version .. %s", cfg.Config.RulesEngineVersion)
	log.Printf(pathLOG+"CaaS Version .......... %s", cfg.Config.CaaSVersion)
	log.Printf(pathLOG+"SLALite Version ....... %s", cfg.Config.SLALiteVersion)
	log.Printf(pathLOG+"IMEC Version .......... %s", cfg.Config.IMECVersion)
	log.Printf(pathLOG+"Server Port ........... %d", cfg.Config.ServerPort)
	log.Printf(pathLOG + "Clusters Configurations:")
	for index := range cfg.Config.Clusters {
		log.Printf(pathLOG+"  => %s", cfg.Config.Clusters[index].Name)
		log.Printf(pathLOG+"     ID: %s", cfg.Config.Clusters[index].ID)
		log.Printf(pathLOG+"     Description: %s", cfg.Config.Clusters[index].Description)
		log.Printf(pathLOG+"     Type: %s", cfg.Config.Clusters[index].Type)
		log.Printf(pathLOG+"     Default Dock / Namespace: %s", cfg.Config.Clusters[index].DefaultDock)
		log.Printf(pathLOG+"     S.O.: %s", cfg.Config.Clusters[index].SO)
		log.Printf(pathLOG+"     Host IP: %s", cfg.Config.Clusters[index].HostIP)
		log.Printf(pathLOG+"     Kubernetes Endpoint: %s", cfg.Config.Clusters[index].KubernetesEndPoint)
		log.Printf(pathLOG+"     Openshift Endpoint: %s", cfg.Config.Clusters[index].OpenshiftEndPoint)
		log.Printf(pathLOG+"     SLALite Endpoint: %s", cfg.Config.Clusters[index].SLALiteEndPoint)
		log.Printf(pathLOG+"     Prometheus Pushgateway Endpoint: %s", cfg.Config.Clusters[index].PrometheusPushgatewayEndPoint)
		log.Printf(pathLOG+"     Prometheus Endpoint: %s", cfg.Config.Clusters[index].PrometheusEndPoint)
	}
	log.Printf(pathLOG + " Swagger UI ............ http://localhost:8333/swaggerui/index.html")
	log.Println(pathLOG + " Starting server [port " + strconv.Itoa(cfg.Config.ServerPort) + "] ...")

	// initialize URLs
	urls.Initialize()

	// adding clusters to DB
	_, err = imec_db.AddConfigInfrsToDB()
	if err == nil {
		clustersMonitoring.StartCheckingClusters()
	} else {
		log.Error(pathLOG+"ERROR (3) Error adding Infrastructure to database. Cluster monitoring was not started. ", err)
	}

	// initialize REST API
	rest_api.InitializeRESTAPI()
}
