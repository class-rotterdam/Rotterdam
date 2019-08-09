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

package adaptation_engine

import (
	caas "atos/rotterdam/caas"
	"atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	"log"
	"net/http"
	"strconv"

	"github.com/nikunjy/rules/parser"
)

// Process
func Process(w http.ResponseWriter, v structs.ViolationInfo) bool {
	log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Processing violation from SLA " + v.Agreement_id + " ...")

	// GET QoS
	dbtaskqos, err := common.ReadTaskQoSValue(v.Client_id)
	if err == nil {
		dbtaskqos.TotalViolations = dbtaskqos.TotalViolations + 1
		err = common.SetTaskQoSValue(v.Client_id, *dbtaskqos)
		if err != nil {
			log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ERROR (1) updating TaskQoS")
		}
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Total violations: " + strconv.Itoa(dbtaskqos.TotalViolations))
		log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Max Allowed: " + strconv.Itoa(dbtaskqos.MaxAllowed))
		if dbtaskqos.TotalViolations > dbtaskqos.MaxAllowed {
			log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] Taking action ...")
			// ADAPT
			dbtask, err := common.ReadTaskValue(v.Client_id)
			var totalNewReplicas int
			if err == nil {
				log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtask.Replicas: " + strconv.Itoa(dbtask.Replicas))
				log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.ScaleFactor: " + strconv.Itoa(dbtaskqos.ScaleFactor))

				if dbtaskqos.Action == "scale_down" {
					totalNewReplicas = dbtask.Replicas - (dbtask.Replicas / dbtaskqos.ScaleFactor)
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] totalNewReplicas: " + strconv.Itoa(totalNewReplicas))
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.MinReplicas: " + strconv.Itoa(dbtaskqos.MinReplicas))
					if totalNewReplicas < dbtaskqos.MinReplicas {
						totalNewReplicas = dbtaskqos.MinReplicas
					}

					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ... decreasing number of replicas from " + strconv.Itoa(dbtask.Replicas) + " to " + strconv.Itoa(totalNewReplicas))
				} else {
					totalNewReplicas = dbtask.Replicas + (dbtask.Replicas * dbtaskqos.ScaleFactor)
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] totalNewReplicas: " + strconv.Itoa(totalNewReplicas))
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] dbtaskqos.MaxReplicas: " + strconv.Itoa(dbtaskqos.MaxReplicas))
					if totalNewReplicas > dbtaskqos.MaxReplicas {
						totalNewReplicas = dbtaskqos.MaxReplicas
					}

					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ... increasing number of replicas from " + strconv.Itoa(dbtask.Replicas) + " to " + strconv.Itoa(totalNewReplicas))
				}

				dbtaskqos.TotalViolations = 0
				err = common.SetTaskQoSValue(v.Client_id, *dbtaskqos)
				if err != nil {
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] ERROR (2) updating TaskQoS")
				}

				if totalNewReplicas == dbtask.Replicas {
					log.Println("Rotterdam > Adaptation-Engine > adaptation engine [Process] No action was taken: total new replicas == current replicas")
					return true
				}

				// take action: scale up or down
				caas.ScaleUpDown(*dbtask, totalNewReplicas)
				return true
			}
		}
	}

	// EVALUATE
	type obj map[string]interface{}
	parser.Evaluate("x eq 1", obj{"x": 1})
	/*parser.Evaluate("x == 1", obj{"x": 1})
	parser.Evaluate("x lt 1", obj{"x": 1})
	parser.Evaluate("x < 1", obj{"x": 1})
	parser.Evaluate("x gt 1", obj{"x": 1})
	parser.Evaluate("x.a == 1 and x.b.c <= 2", obj{
		"x": obj{
			"a": 1,
			"b": obj{
				"c": 2,
			},
		},
	})
	parser.Evaluate("y == 4 and (x > 1)", obj{"x": 1})
	parser.Evaluate("y == 4 and (x IN [1,2,3])", obj{"x": 1})
	parser.Evaluate("y == 4 and (x eq 1.2.3)", obj{"x": "1.2.3"})*/

	return false
}
