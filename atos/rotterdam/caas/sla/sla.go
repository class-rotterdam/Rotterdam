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
	db "atos/rotterdam/database/caas"
	constants "atos/rotterdam/globals/constants"
	structs "atos/rotterdam/globals/structs"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// path used in logs
const pathLOG string = "Rotterdam > CAAS > SLA : "

///////////////////////////////////////////////////////////////////////////////

/*
createCOMPSsSLAAgreemnt ...
*/
func createCOMPSsSLAAgreemnt(task structs.CLASS_TASK, cnt int) (string, error) {
	log.Println(pathLOG + "[createCOMPSsSLAAgreemnt] Generating COMPSs SLA agrement for task [" + task.ID + "] ...")

	// Create uuid
	agreementID := strings.Replace(task.ID, "-", "_", -1)

	// Create SLA Agreement structure
	var jsonAgreement *structs.SLA_AGREEMENT
	jsonAgreement = new(structs.SLA_AGREEMENT)

	jsonAgreement.ID = agreementID
	jsonAgreement.Name = "agreement_" + agreementID + "_" + strconv.Itoa(cnt)
	jsonAgreement.State = "started"
	jsonAgreement.Details.ID = agreementID
	jsonAgreement.Details.Name = "agreement_" + agreementID + "_" + strconv.Itoa(cnt)
	jsonAgreement.Details.Type = "agreement"
	jsonAgreement.Details.Provider.ID = "CLASSProvider"
	jsonAgreement.Details.Provider.Name = "CLASS Platform"
	jsonAgreement.Details.Client.ID = task.ID
	jsonAgreement.Details.Client.Name = "compss_task_" + task.ID
	jsonAgreement.Details.Creation = cfg.Config.SLAs.CreationDate     //"2019-01-01T00:00:00Z"
	jsonAgreement.Details.Expiration = cfg.Config.SLAs.ExpirationDate //"2025-01-01T00:00:00Z"

	resGuarantees, dbtaskqos, _ := CreateCOMPSsGuarantees(task)

	jsonAgreement.Details.Guarantees = resGuarantees

	// Call to SLALite API to create the agreement
	// ==> curl -k -X POST -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements
	status, _, err := common.HTTPPOST(
		cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements",
		false,
		jsonAgreement)

	if err != nil {
		log.Error(pathLOG+"[createCOMPSsSLAAgreemnt] ERROR", err)
		return "Error creating the SLA Agreement", err
	}

	// Save QoS in DB
	if dbtaskqos != nil {
		err = db.SetTaskQoSValue(task.ID, *dbtaskqos)
		if err != nil {
			log.Error(pathLOG+"[createCOMPSsSLAAgreemnt] ERROR creating QoS:", err)
		}
	} else {
		log.Println(pathLOG + "[createCOMPSsSLAAgreemnt] ERROR creating QoS. dbtaskqos is nil")
	}

	log.Println(pathLOG + "[createCOMPSsSLAAgreemnt] RESPONSE: OK (" + strconv.Itoa(status) + ")")

	return agreementID, nil
}

/*
createSLAAgreemnt ...
*/
func createSLAAgreemnt(task structs.CLASS_TASK) (string, error) {
	if task.Qos.Name == "" || task.Qos.Name == "None" {
		log.Println(pathLOG + "[createSLAAgreemnt] No SLA agrement generation for task [" + task.ID + "]. QoS value is None.")
		return constants.SLANotDefined, nil
	}
	log.Println(pathLOG + "[createSLAAgreemnt] Generating SLA agrement for task [" + task.ID + "] ...")

	// Create uuid
	agreementID := strings.Replace(task.ID, "-", "_", -1)

	// Create SLA Agreement structure
	var jsonAgreement *structs.SLA_AGREEMENT
	jsonAgreement = new(structs.SLA_AGREEMENT)

	jsonAgreement.ID = agreementID
	jsonAgreement.Name = agreementID + "_" + task.Qos.Name
	jsonAgreement.State = "started"
	jsonAgreement.Location = task.Cluster
	jsonAgreement.Details.ID = agreementID
	jsonAgreement.Details.Name = agreementID + "_" + task.Qos.Name
	jsonAgreement.Details.Type = "agreement"
	jsonAgreement.Details.Provider.ID = "CLASSProvider"
	jsonAgreement.Details.Provider.Name = "CLASS Platform"
	jsonAgreement.Details.Client.ID = task.ID
	jsonAgreement.Details.Client.Name = "CLASS task " + task.ID
	jsonAgreement.Details.Creation = cfg.Config.SLAs.CreationDate     //"2019-01-01T00:00:00Z"
	jsonAgreement.Details.Expiration = cfg.Config.SLAs.ExpirationDate //"2025-01-01T00:00:00Z"

	resGuarantees, dbtaskqos, _ := CreateGuarantees(task)

	jsonAgreement.Details.Guarantees = resGuarantees

	// DEBUG
	strTxt, err := structs.CommSLAStructToString(*jsonAgreement)
	if err == nil {
		log.Println(pathLOG + "[createSLAAgreemnt] jsonAgreement: " + strTxt)
	}

	// Call to SLALite API to create the agreement
	// ==> curl -k -X POST -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements
	status, _, err := common.HTTPPOST(
		cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements",
		false,
		jsonAgreement)

	if err != nil {
		log.Error(pathLOG+"[createSLAAgreemnt] ERROR", err)
		return "Error creating the SLA Agreement", err
	}

	// Save QoS in DB
	if dbtaskqos != nil {
		err = db.SetTaskQoSValue(task.ID, *dbtaskqos)
		if err != nil {
			log.Error(pathLOG+"[createSLAAgreemnt] ERROR creating QoS:", err)
		}
	} else {
		log.Println(pathLOG + "[createSLAAgreemnt] ERROR creating QoS. dbtaskqos is nil")
	}

	log.Println(pathLOG + "[createSLAAgreemnt] RESPONSE: OK (" + strconv.Itoa(status) + ")")

	return agreementID, nil
}

/*
createInfrSLA ...
*/
func createInfrSLA(cluster string, qos structs.CLASS_TASK_QOS) (string, error) {
	log.Println(pathLOG + "[createInfrSLA] Generating SLA agrement for CLUSTER [" + cluster + "] ...")

	// Create uuid
	now := time.Now()      // current local time
	nsec := now.UnixNano() // number of nanoseconds since January 1, 1970 UTC
	agreementID := cluster + strconv.FormatInt(nsec, 10)

	// Create SLA Agreement structure
	var jsonAgreement *structs.SLA_AGREEMENT
	jsonAgreement = new(structs.SLA_AGREEMENT)

	jsonAgreement.ID = agreementID
	jsonAgreement.Name = cluster //agreementID + "_INFR_SLA"
	jsonAgreement.State = "started"
	jsonAgreement.Location = cluster
	jsonAgreement.Details.ID = agreementID
	jsonAgreement.Details.Name = cluster //agreementID + "_INFR_SLA"
	jsonAgreement.Details.Type = "agreement"
	jsonAgreement.Details.Provider.ID = "CLASSProvider"
	jsonAgreement.Details.Provider.Name = "CLASS Platform"
	jsonAgreement.Details.Client.ID = cluster
	jsonAgreement.Details.Client.Name = "CLASS Cluster " + cluster
	jsonAgreement.Details.Creation = cfg.Config.SLAs.CreationDate     //"2019-01-01T00:00:00Z"
	jsonAgreement.Details.Expiration = cfg.Config.SLAs.ExpirationDate //"2025-01-01T00:00:00Z"

	resGuarantees, dbtaskqos, _ := CreateInfrGuarantees(cluster, qos)

	jsonAgreement.Details.Guarantees = resGuarantees

	// DEBUG
	strTxt, err := structs.CommSLAStructToString(*jsonAgreement)
	if err == nil {
		log.Println(pathLOG + "[createInfrSLA] jsonAgreement: " + strTxt)
	}

	// Call to SLALite API to create the agreement
	// ==> curl -k -X POST -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements
	status, _, err := common.HTTPPOST(
		cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements",
		false,
		jsonAgreement)

	if err != nil {
		log.Error(pathLOG+"[createInfrSLA] ERROR", err)
		return "Error creating the SLA Agreement", err
	}

	// Save QoS in DB
	if dbtaskqos != nil {
		err = db.SetTaskQoSValue(cluster, *dbtaskqos)
		if err != nil {
			log.Error(pathLOG+"[createInfrSLA] ERROR creating QoS:", err)
		}
	} else {
		log.Println(pathLOG + "[createInfrSLA] ERROR creating QoS. dbtaskqos is nil")
	}

	log.Println(pathLOG + "[createInfrSLA] RESPONSE: OK (" + strconv.Itoa(status) + ")")

	return agreementID, nil
}

/*
startAgreemnt Call to SLALite API to start the agreement
*/
func startAgreemnt(agreementID string) (string, error) {
	log.Println(pathLOG + "[startAgreemnt] Starting SLA agrement [" + agreementID + "] ...")
	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements/a03/start
	data := url.Values{}
	status, _, err := common.HTTPPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreementID+"/start", false, data)

	if err != nil {
		log.Error(pathLOG+"[startAgreemnt] ERROR", err)
		return "Error starting the SLA Agreement", err
	}
	log.Println(pathLOG + "[startAgreemnt] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// stopSLA
func stopSLA(agreementID string, tries int) (string, error) {
	status, _, err := common.HTTPPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreementID+"/stop", false, url.Values{})
	if err != nil {
		log.Error(pathLOG+"[stopSLA] ERROR Trying to stop SLA ", err)

		if tries == 0 {
			return "Error stopping the SLA Agreement", err
		}

		log.Println(pathLOG + "[stopSLA] Trying to stop the SLA again in 15 seconds ...")

		time.Sleep(15 * time.Second)
		return stopSLA(agreementID, tries-1)
	}

	log.Println(pathLOG + "[stopSLA] RESPONSE: OK")
	return strconv.Itoa(status), nil
}

/*
stopAgreemnt Call to SLALite API to stop the agreement
*/
func stopAgreemnt(agreementID string) (string, error) {
	agreementID = strings.Replace(agreementID, "-", "_", -1)
	log.Println(pathLOG + "[stopAgreemnt] Stopping SLA agrement [" + agreementID + "] ...")
	return stopSLA(agreementID, 3)
}

/*
terminateAgreemnt Call to SLALite API to terminate the agreement
*/
func terminateAgreemnt(agreementID string) (string, error) {
	agreementID = strings.Replace(agreementID, "-", "_", -1)
	log.Println(pathLOG + "[terminateAgreemnt] Terminating SLA agrement [" + agreementID + "] ...")

	// ==> curl -k -X PUT -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/agreements/a03/terminate
	data := url.Values{}
	status, _, err := common.HTTPPUT(cfg.Config.Clusters[0].SLALiteEndPoint+"/agreements/"+agreementID+"/terminate", false, data)

	if err != nil {
		log.Error(pathLOG+"[terminateAgreemnt] ERROR", err)
		return "Error terminating the SLA Agreement", err
	}
	log.Println(pathLOG + "[terminateAgreemnt] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

///////////////////////////////////////////////////////////////////////////////

/*
CreateStartCOMPSsSLA creates and starts an SLA
*/
func CreateStartCOMPSsSLA(task structs.CLASS_TASK) error {
	log.Println(pathLOG + "[CreateStartCOMPSsSLA] Creating COMPSs SLA agreement ...")
	agreementID, err := createCOMPSsSLAAgreemnt(task, 1)
	if err != nil {
		log.Error(pathLOG + "[CreateStartCOMPSsSLA] ERROR creating COMPSs SLA Agreement")
	} else {
		log.Println(pathLOG + "[CreateStartCOMPSsSLA] Starting COMPSs SLA agreement " + agreementID + " ...")
		_, err = startAgreemnt(agreementID)
		if err != nil {
			log.Error(pathLOG + "[CreateStartCOMPSsSLA] ERROR starting COMPSs SLA Agreement")
		} else {
			log.Println(pathLOG + "[CreateStartCOMPSsSLA] COMPSs SLA agreement " + agreementID + " started")

			dbTask, err := db.ReadTaskValue(task.ID)
			if err == nil {
				dbTask.AgreementId = agreementID
				err = db.SetTaskValue(task.ID, *dbTask)
				if err == nil {
					log.Println(pathLOG + "[CreateStartCOMPSsSLA] Task updated")
				} else {
					log.Error(pathLOG+"[CreateStartCOMPSsSLA] ERROR setting 'id agreement' in task: ", err)
				}
			}
		}
	}

	return err
}

/*
CreateStartSLA creates and starts an SLA
*/
func CreateStartSLA(task structs.CLASS_TASK) error {
	log.Println(pathLOG + "[CreateStartSLA] Creating SLA agreement ...")
	agreementID, err := createSLAAgreemnt(task)
	if err != nil {
		log.Error(pathLOG+"[CreateStartSLA] ERROR creating SLA Agreement: ", err)
	} else {
		if agreementID == constants.SLANotDefined {
			log.Println(pathLOG + "[CreateStartSLA] Task is not using any SLA")
		} else {
			log.Println(pathLOG + "[CreateStartSLA] Starting SLA agreement " + agreementID + " ...")
			_, err = startAgreemnt(agreementID)
			if err != nil {
				log.Error(pathLOG+"[CreateStartSLA] ERROR starting SLA Agreement: ", err)
			} else {
				log.Println(pathLOG + "[CreateStartSLA] SLA agreement " + agreementID + " started")
			}
		}

		// update task
		dbTask, err := db.ReadTaskValue(task.ID)
		if err == nil {
			dbTask.AgreementId = agreementID
			err = db.SetTaskValue(task.ID, *dbTask)
			if err != nil {
				log.Error(pathLOG+"[CreateStartSLA] ERROR setting 'id agreement' in task ", err)
			}
		}
	}

	return err
}

/*
StopTerminateSLA stops and terminates an SLA agreement after TASK is successfully deployed
*/
func StopTerminateSLA(taskID string) {
	log.Println(pathLOG + "[StopTerminateSLA] Stopping SLA agreement from task [" + taskID + "] ...")

	agreementID := taskID

	// TODO agreement id to DB
	_, err := stopAgreemnt(agreementID)

	if err != nil {
		log.Error(pathLOG + "[StopTerminateSLA] ERROR Stopping SLA with id = " + agreementID)
	} else {
		log.Println(pathLOG + "[StopTerminateSLA] Terminating SLA agreement " + agreementID + " ...")
		_, err := terminateAgreemnt(agreementID)
		if err != nil {
			log.Println(pathLOG + "[StopTerminateSLA] ERROR Terminating SLA Agreement")
		} else {
			log.Println(pathLOG + "[StopTerminateSLA] SLA agreement " + agreementID + " terminated")
			_, err = db.DBDeleteTaskQos(taskID)
			if err != nil {
				log.Error(pathLOG+"[StopTerminateSLA] ERROR deleting QoS:", err)
			}
		}
	}
}

/*
AddPromMetric Call to SLALite API to add a new Prometheus metric
*/
func AddPromMetric(metricID string) error {
	log.Println(pathLOG + "[AddPromMetric] Adding new Prometheus metric [" + metricID + "] ...")
	// ==> curl -k -X POST -d @agreement.json http://rotterdam-slalite60.192.168.7.28.xip.io/metrics/metric_name
	data := url.Values{}
	_, _, err := common.HTTPPOST(cfg.Config.Clusters[0].SLALiteEndPoint+"/metrics/"+metricID, false, data)

	if err != nil {
		log.Error(pathLOG+"[AddPromMetric] ERROR", err)
		return err
	}
	log.Println(pathLOG + "[AddPromMetric] RESPONSE: OK")

	return nil
}
