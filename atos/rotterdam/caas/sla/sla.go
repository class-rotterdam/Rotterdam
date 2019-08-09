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
	"atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"log"
	"net/url"
	"strconv"

	"github.com/lithammer/shortuuid"
)

// CreateSLAAgreemnt:
func CreateSLAAgreemnt(task structs.CLASS_TASK) (string, error) {
	if task.Qos.Name == "" || task.Qos.Name == "None" {
		log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] No SLA agrement generation for task [" + task.Name + "]. QoS value is None.")
		return "NO_SLA", nil
	} else {
		log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] Generating SLA agrement for task [" + task.Name + "] ...")

		// Create uuid
		alphabet := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx"
		agreement_id := shortuuid.NewWithAlphabet(alphabet)

		// Create SLA Agreement structure
		var jsonAgreement *structs.SLA_AGREEMENT
		jsonAgreement = new(structs.SLA_AGREEMENT)

		jsonAgreement.Id = agreement_id
		jsonAgreement.Name = "Agreement " + agreement_id
		jsonAgreement.State = "started"
		jsonAgreement.Details.Id = agreement_id
		jsonAgreement.Details.Name = "Agreement " + agreement_id
		jsonAgreement.Details.Type = "agreement"
		jsonAgreement.Details.Provider.Id = "CLASSProvider"
		jsonAgreement.Details.Provider.Name = "CLASS Platform"
		jsonAgreement.Details.Client.Id = task.Name
		jsonAgreement.Details.Client.Name = "CLASS task " + task.Name
		jsonAgreement.Details.Creation = cfg.Config.SLAs.CreationDate     //"2019-01-01T00:00:00Z"
		jsonAgreement.Details.Expiration = cfg.Config.SLAs.ExpirationDate //"2025-01-01T00:00:00Z"

		res_guarantees, dbtaskqos, _ := CreateGuarantees(task)

		jsonAgreement.Details.Guarantees = res_guarantees

		// Call to SLALite API to create the agreement
		// ==> curl -k -X POST -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements
		status, _, err := common.HttpPOST_GenericStruct(
			cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements",
			jsonAgreement)

		if err != nil {
			log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] ERROR", err)
			return "Error creating the SLA Agreement", err
		}

		// Save QoS in DB
		if dbtaskqos != nil {
			err = common.SetTaskQoSValue(task.Name, *dbtaskqos)
			if err != nil {
				log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] ERROR creating QoS:", err)
			}
		} else {
			log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] ERROR creating QoS. dbtaskqos is nil")
		}

		log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] RESPONSE: OK (" + strconv.Itoa(status) + ")")

		return agreement_id, nil
	}
}

// StartAgreemnt: Call to SLALite API to start the agreement
func StartAgreemnt(agreement_id string) (string, error) {
	log.Println("Rotterdam > CAAS > SLA [StartAgreemnt] Starting SLA agrement [" + agreement_id + "] ...")
	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements/a03/start
	data := url.Values{}
	status, _, err := common.HttpPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreement_id+"/start", data)

	if err != nil {
		log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] ERROR", err)
		return "Error starting the SLA Agreement", err
	}
	log.Println("Rotterdam > CAAS > SLA [CreateSLAAgreemnt] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// StopAgreemnt: Call to SLALite API to stop the agreement
func StopAgreemnt(agreement_id string) (string, error) {
	log.Println("Rotterdam > CAAS > SLA [StopAgreemnt] Stopping SLA agrement [" + agreement_id + "] ...")
	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements/a03/stop
	data := url.Values{}
	status, _, err := common.HttpPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreement_id+"/stop", data)

	if err != nil {
		log.Println("Rotterdam > CAAS > SLA [StopAgreemnt] ERROR", err)
		return "Error stopping the SLA Agreement", err
	}
	log.Println("Rotterdam > CAAS > SLA [StopAgreemnt] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// TerminateAgreemnt: Call to SLALite API to terminate the agreement
func TerminateAgreemnt(agreement_id string) (string, error) {
	log.Println("Rotterdam > CAAS > SLA [TerminateAgreemnt] Terminating SLA agrement [" + agreement_id + "] ...")
	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements/a03/terminate
	data := url.Values{}
	status, _, err := common.HttpPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreement_id+"/terminate", data)

	if err != nil {
		log.Println("Rotterdam > CAAS > SLA [TerminateAgreemnt] ERROR", err)
		return "Error terminating the SLA Agreement", err
	}
	log.Println("Rotterdam > CAAS > SLA [TerminateAgreemnt] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// CreateStartSLA:
func CreateStartSLA(task structs.CLASS_TASK) error {
	log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] Creating SLA agreement ...")
	agreement_id, err := CreateSLAAgreemnt(task)
	if err != nil {
		log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] ERROR creating SLA Agreement")
	} else {
		log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] Starting SLA agreement " + agreement_id + " ...")
		_, err = StartAgreemnt(agreement_id)
		if err != nil {
			log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] ERROR starting SLA Agreement")
		} else {
			log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] SLA agreement " + agreement_id + " started")

			dbTask, err := common.ReadTaskValue(task.Name)
			if err == nil {
				dbTask.AgreementId = agreement_id
				err = common.SetTaskValue(task.Name, *dbTask)
				if err == nil {
					log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] Task updated")
				} else {
					log.Println("Rotterdam > CAAS > SLA [CreateStartSLA] ERROR setting 'id agreement' in task")
				}
			}
		}
	}

	return err
}

// StopTerminateSLA: stop and terminate SLA agreement after TASK is successfully deployed
func StopTerminateSLA(task_name string) {
	log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] Stopping SLA agreement from task [" + task_name + "] ...")

	dbTask, err := common.ReadTaskValue(task_name)
	if err == nil {
		agreement_id := dbTask.AgreementId
		// TODO agreement id to DB
		_, err := StopAgreemnt(agreement_id)

		if err != nil {
			log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] ERROR Stopping SLA Agreement")
		} else {
			log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] Terminating SLA agreement " + agreement_id + " ...")
			_, err := TerminateAgreemnt(agreement_id)
			if err != nil {
				log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] ERROR Terminating SLA Agreement")
			} else {
				log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] SLA agreement " + agreement_id + " terminated")
				_, err = common.DBDeleteTaskQos(task_name)
				if err != nil {
					log.Println("Rotterdam > CAAS > SLA [StopTerminateSLA] ERROR deleting QoS:", err)
				}
			}
		}
	}
}
