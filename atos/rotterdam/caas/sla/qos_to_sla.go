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

package sla

import (
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"log"
)

// get QoS template from global list
func GetQoSElem(name string) (cfg.CLASS_QOS_TEMPLATE__ELEM, bool) {
	for _, qoselem := range cfg.QosTemplates {
		if qoselem.GuaranteeName == name {
			return qoselem, true
		}
	}

	return cfg.CLASS_QOS_TEMPLATE__ELEM{}, false
}

// create a DB Task QoS element
func genDbTaskQos(taskName string, elem cfg.CLASS_QOS_TEMPLATE__ELEM) *structs.DB_TASK_QOS {
	dbtaskqos := &structs.DB_TASK_QOS{
		DbId:            structs.DB_TABLE_QOS,
		Id:              taskName,
		IdTask:          taskName,
		Type:            elem.Type,
		Action:          elem.Action,
		ScaleFactor:     elem.ScaleFactor,
		MaxReplicas:     30,
		MinReplicas:     1,
		Guarantee:       elem.GuaranteeName,
		TotalViolations: 0,
		MaxAllowed:      elem.MaxAllowed}

	return dbtaskqos
}

// create Guarantees element
func genGuarantees(taskName string, elem cfg.CLASS_QOS_TEMPLATE__ELEM) []structs.SLA_AGREEMENT_DETAILS_GUARANTIES {
	guarantees := make([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, 0)

	for _, v := range elem.Guarantees {
		g := structs.SLA_AGREEMENT_DETAILS_GUARANTIES{}
		g.Name = elem.GuaranteeName // v.Name,   // "kubelet_running_pod_count"
		g.Constraint = v.Constraint // "kubelet_running_pod_count < 50"

		guarantees = append(guarantees, g)
	}

	return guarantees
}

// CreateSLAAgreemnt:
func CreateGuarantees(task structs.CLASS_TASK) ([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, *structs.DB_TASK_QOS, error) {
	log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Generating Guarantees from task qos [" + task.Qos.Name + "] ...")

	var guarantees []structs.SLA_AGREEMENT_DETAILS_GUARANTIES

	if task.Qos.Custom != nil && len(task.Qos.Custom) > 0 {
		log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Translating custom task QoS to SLA guarantees ...")
		// TODO support for more than one guarantee
		//for index, _ := range task.Qos.Custom {
		if task.Qos.Custom[0].Guarantees != nil {
			guarantees = make([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, 0)

			guarantee := structs.SLA_AGREEMENT_DETAILS_GUARANTIES{}
			guarantee.Name = task.Qos.Custom[0].Guarantees[0].Metric
			guarantee.Constraint = task.Qos.Custom[0].Guarantees[0].Metric + " " +
				task.Qos.Custom[0].Guarantees[0].Condition + " " + task.Qos.Custom[0].Guarantees[0].Value

			guarantees = append(guarantees, guarantee)

			return guarantees, nil, nil
		}
		//}
	} else {
		log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Translating catalog task QoS to SLA guarantees ...")

		qosTemplate, found := GetQoSElem(task.Qos.Name)
		if found {
			guarantees = genGuarantees(task.Name, qosTemplate)
			dbtaskqos := genDbTaskQos(task.Name, qosTemplate)
			return guarantees, dbtaskqos, nil
		} else {
			log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] ERROR QoS element not found: [" + task.Qos.Name + "]")
		}
	}

	return guarantees, nil, nil
}
