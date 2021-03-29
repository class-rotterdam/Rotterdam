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

package adaptation_engine

import (
	caas "atos/rotterdam/caas"
	"atos/rotterdam/caas/common"
	cfg "atos/rotterdam/config"
	db "atos/rotterdam/database/caas"
	db_imec "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// floatToString
func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

// resetDeadLinesMissedMetric
func resetDeadLinesMissedMetric(clusterIndex int, id string, job string) (string, error) {
	// CALL to Prometheus Pushgateway API to delete a deployment
	log.Println("Rotterdam > Adaptation-Engine > [resetDeadLinesMissedMetric] Deleting metric from Prometheus ...")
	log.Println("Rotterdam > Adaptation-Engine > [resetDeadLinesMissedMetric] DELETE [" + cfg.Config.Clusters[clusterIndex].PrometheusPushgatewayEndPoint + "/metrics/job/deadlines_missed_" + id + "]")

	// raw_data
	rawData := `# TYPE deadlines_missed_` + id + ` counter
	deadlines_missed_` + id + ` 0
	`

	// reset value
	status, _, err := common.HTTPPOSTRawData(
		cfg.Config.Clusters[clusterIndex].PrometheusPushgatewayEndPoint+"/metrics/job/"+job, true, rawData)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > [resetDeadLinesMissedMetric] ERROR", err)
		return strconv.Itoa(status), err
	}

	log.Println("Rotterdam > Adaptation-Engine > [resetDeadLinesMissedMetric] RESPONSE: OK")
	return strconv.Itoa(status), nil
}

/*
sendDataToPrometheusPushgateway data (workers) to prometheus via pushgateway
*/
func sendDataToPrometheusPushgateway(dbTask structs.DB_TASK, totalNewReplicas int) {
	log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] Sending data to Pushgateway (Prometheus) ... ")

	cl, err := db_imec.GetInfrByID(dbTask.ClusterId) // ([]DB_INFRASTRUCTURE_CLUSTER, error)
	if err == nil {
		pushgatewayURL := cl[0].PrometheusPushgatewayEndPoint

		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] Sending metrics to prometheus pushgateway [" + pushgatewayURL + "] ...")

		s := `
		workers_` + dbTask.Id + ` ` + strconv.Itoa(totalNewReplicas) + `
		`

		ioreader := strings.NewReader(s)

		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] POST " + pushgatewayURL + "/metrics/job/sla/instance/violations")
		res, err := http.Post(pushgatewayURL+"/metrics/job/sla/instance/violations", "binary/octet-stream", ioreader)
		if err != nil {
			log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] Error (1): " + err.Error())
		}
		defer res.Body.Close()
		message, _ := ioutil.ReadAll(res.Body)
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] Response: " + string(message))
	} else {
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [sendDataToPrometheusPushgateway] Error (2): " + err.Error())
	}
}

/*
Process handles violatiokns from SLALite
*/
func Process(w http.ResponseWriter, v structs.ViolationInfo) bool {
	log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Processing violation from SLA [" + v.Agreement_id + "] and Client ID [" + v.Client_id + "] ...")

	// GET QoS
	//dbtaskqos, err := common.ReadTaskQoSValue(v.Client_id)
	dbtaskqos, err := db.ReadTaskQoSValue(v.Agreement_id)
	if err == nil {
		dbtaskqos.TotalViolations = dbtaskqos.TotalViolations + 1
		err = db.SetTaskQoSValue(v.Client_id, *dbtaskqos)
		if err != nil {
			log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ERROR (1) updating TaskQoS")
		}
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Total violations: " + strconv.Itoa(dbtaskqos.TotalViolations))
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Max Allowed: " + strconv.Itoa(dbtaskqos.MaxAllowed))
		if dbtaskqos.TotalViolations > dbtaskqos.MaxAllowed {
			log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Taking action ...")
			// if COMPSs metric ==> delete from Prometheus / Pushgateway
			if dbtaskqos.Type == "app-compss" {
				resetDeadLinesMissedMetric(0, dbtaskqos.IdTask, "compss")
			}

			// ADAPT
			dbtask, err := db.ReadTaskValue(v.Client_id)
			var totalNewReplicas int
			if err == nil {
				log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtask.Replicas: " + strconv.Itoa(dbtask.Replicas))
				log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.ScaleFactor: " + floatToString(dbtaskqos.ScaleFactor))

				if dbtaskqos.Action == "scale_in" {
					resScaling := math.Floor(float64(dbtask.Replicas) * dbtaskqos.ScaleFactor)
					totalNewReplicas = int(resScaling) // dbtask.Replicas - (dbtask.Replicas / dbtaskqos.ScaleFactor)
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] totalNewReplicas: " + strconv.Itoa(totalNewReplicas))
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.MinReplicas: " + strconv.Itoa(dbtaskqos.MinReplicas))
					if totalNewReplicas < dbtaskqos.MinReplicas {
						totalNewReplicas = dbtaskqos.MinReplicas
					}

					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ... decreasing number of replicas from " + strconv.Itoa(dbtask.Replicas) + " to " + strconv.Itoa(totalNewReplicas))
				} else {
					resScaling := math.Ceil(float64(dbtask.Replicas) * dbtaskqos.ScaleFactor)
					totalNewReplicas = int(resScaling) // dbtask.Replicas + (dbtask.Replicas * dbtaskqos.ScaleFactor)
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] totalNewReplicas: " + strconv.Itoa(totalNewReplicas))
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.MaxReplicas: " + strconv.Itoa(dbtaskqos.MaxReplicas))
					if totalNewReplicas > dbtaskqos.MaxReplicas {
						totalNewReplicas = dbtaskqos.MaxReplicas
					}

					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ... increasing number of replicas from " + strconv.Itoa(dbtask.Replicas) + " to " + strconv.Itoa(totalNewReplicas))
				}

				dbtaskqos.TotalViolations = 0
				err = db.SetTaskQoSValue(v.Client_id, *dbtaskqos)
				if err != nil {
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ERROR (2) updating TaskQoS")
				}

				if totalNewReplicas == dbtask.Replicas {
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] No action was taken: total new replicas == current replicas")
					return true
				}

				// take action: scale up or down
				caas.ScaleUpDown(*dbtask, totalNewReplicas)

				// save to Prometheus / Pushgateway
				sendDataToPrometheusPushgateway(*dbtask, totalNewReplicas)

				return true
			}
		}
	}

	return false
}
