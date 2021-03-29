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

package sla

import (
	"atos/rotterdam/caas/common"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	structs "atos/rotterdam/globals/structs"
	"encoding/json"
	"strconv"
	"strings"
)

/*
PredictionParamenters ...
*/
type PredictionParamenters struct {
	Id            string `json:"id,omitempty"`
	NumWorkers    string `json:"numWorkers,omitempty"`
	ExecutionTime string `json:"executionTime,omitempty"`
	Result        string `json:"result,omitempty"`
}

// extract execution time from QoS
func getExecutionTime(q []structs.CLASS_COMPSS_TASK_QOS) (string, error) {
	for _, value := range q {
		if strings.HasPrefix(value.Metric, "execution_time") {
			return strconv.Itoa(value.Value), nil
		}
	}

	return "", nil
}

/*
stringToPredictionParamenters Parses a string to a struct of type PredictionParamenters
*/
func stringToPredictionParamenters(ct string) (*PredictionParamenters, error) {
	p := &PredictionParamenters{}
	err := json.Unmarshal([]byte(ct), p)
	if err != nil {
		log.Error(pathLOG+"[stringToPredictionParamenters] ERROR", err)
		return p, err
	}

	return p, nil
}

/*
Predict ...
*/
func Predict(task *structs.CLASS_TASK) (*structs.CLASS_TASK, error) {
	log.Println(pathLOG + "[Predict] Call to SLA-Predict ...")

	var params *PredictionParamenters
	params = new(PredictionParamenters)
	params.Id = "compssAPP"
	params.NumWorkers = strconv.Itoa(task.Replicas)

	v, err := getExecutionTime(task.QoSCOMPSs)
	if err != nil {
		log.Error(pathLOG+"[Predict] ERROR (1): ", err)
		return task, err
	}

	params.ExecutionTime = v

	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/sla-predict
	_, data, err := common.HTTPPUTString(cfg.Config.Clusters[0].SLALiteEndPoint+"/sla-predict", false, params)
	if err != nil {
		log.Error(pathLOG+"[Predict] ERROR (2): ", err)
		log.Println(pathLOG + "[Predict] RESPONSE: Returning initial number of replicas")
		return task, err
	} else {
		log.Println(pathLOG + "[Predict] RESPONSE data: " + data)
		data = strings.Replace(data, "[\"http://prometheus.192.168.7.28.xip.io\"]", "", 1)
		log.Println(pathLOG + "[Predict] RESPONSE data: " + data)
		// data (string) to PredictionParamenters
		predRes, err := stringToPredictionParamenters(data)
		if err != nil {
			log.Error(pathLOG+"[Predict] ERROR (3): ", err)
			log.Println(pathLOG + "[Predict] RESPONSE: Returning initial number of replicas")
			return task, err
		}

		log.Println(pathLOG + "[Predict] RESPONSE: " + predRes.NumWorkers + ", " + predRes.ExecutionTime + " => " + predRes.Result)
		if n, err := strconv.Atoi(predRes.Result); err == nil {
			if n != -1 {
				task.Replicas = n
				log.Println(pathLOG + "[Predict] RESPONSE: Number of replicas taken from SLA Predict response: " + strconv.Itoa(task.Replicas))
			}
		} else {
			log.Error(pathLOG+"[Predict] ERROR (4): ", err)
			log.Println(pathLOG + "[Predict] RESPONSE: Returning initial number of replicas")
			return task, err
		}
	}

	log.Println(pathLOG + "[Predict] RESPONSE: Number of replicas: " + strconv.Itoa(task.Replicas))

	return task, nil
}
