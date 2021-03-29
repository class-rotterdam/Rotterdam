/*
Copyright 2017 Atos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

CLASS Project: https://class-project.eu/

@author: ATOS
*/
package ml

import (
	"SLALite/model"
	"SLALite/utils"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

/*
URLMLSubsystem ML Endpoint
*/
var URLMLSubsystem string // endpoint

// initialization function
func init() {
	log.SetLevel(log.DebugLevel)
	log.Println("SLALite > Assessment > ML > init > Initializing ML adapter ...")

	// ENVIRONMENT VARIABLES
	// URLMLSubsystem
	if os.Getenv("URLMLSubsystem") != "" {
		log.Println("SLALite > Assessment > ML > init > Setting 'URLMLSubsystem' value ... " + os.Getenv("URLMLSubsystem"))
		URLMLSubsystem = os.Getenv("URLMLSubsystem")
	} else {
		URLMLSubsystem = "http://192.168.7.42:5002"
	}
	log.Println("SLALite > Assessment > ML > init > 'URLMLSubsystem' = " + URLMLSubsystem)
}

/*
SLAPredict call to ML subsystem to get recommended number of workers
*/
func SLAPredict(pred *model.PredictionParamenters) (*model.PredictionParamenters, error) {
	log.Println("SLALite > Assessment > ML [SLAPredict] PredictionParamenters: " + pred.NumWorkers + ", " + pred.ExecutionTime)

	log.Println("SLALite > Assessment > ML [SLAPredict] Call to: " + URLMLSubsystem + "/predictSLA?workers=" + pred.NumWorkers + "&exectime=" + pred.ExecutionTime)
	_, data, err := utils.HTTPGETString(URLMLSubsystem+"/predictSLA?workers="+pred.NumWorkers+"&exectime="+pred.ExecutionTime, false, "")
	if err != nil {
		log.Error("SLALite > Assessment > ML [SLAPredict] ERROR (1) executing query [GET "+URLMLSubsystem+"/predictSLA?workers="+pred.NumWorkers+"&exectime="+pred.ExecutionTime+"]", err)
		pred.Result = "error"
	} else {
		log.Println("SLALite > Assessment > ML [SLAPredict] Results: " + data)

		// chequear numero
		if n, err := strconv.Atoi(data); err == nil {
			if n != -1 {
				pred.Result = data
			}
		} else {
			log.Error("SLALite > Assessment > ML [SLAPredict] ERROR (2): ", err)
			pred.Result = "error"
		}

	}

	return pred, err
}
