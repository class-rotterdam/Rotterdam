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

package caas

import (
	adapt_engine "atos/rotterdam/adaptation-engine/monitoring"
	adapters "atos/rotterdam/caas/adapters"
	adpcommon "atos/rotterdam/caas/adapters/common"
	k8s_adapter "atos/rotterdam/caas/adapters/kubernetes"
	ops_adapter "atos/rotterdam/caas/adapters/openshift"
	common "atos/rotterdam/caas/common"
	sla "atos/rotterdam/caas/sla"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	db "atos/rotterdam/database/caas"
	db_imec "atos/rotterdam/database/imec"
	constants "atos/rotterdam/globals/constants"
	structs "atos/rotterdam/globals/structs"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// adapter
var a adapters.Adapter

///////////////////////////////////////////////////////////////////////////////

// init
func init() {
	log.Println("Rotterdam > CAAS > adapter [init] Selected cluster index set to '0' ...")

	qosTemplatesPath := "./config/qos_templates.json"
	log.Println("Rotterdam > CAAS > adapter [init] Reading content of qos_templates.json file [" + qosTemplatesPath + "] ...")

	file, err := os.Open(qosTemplatesPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.QosTemplates)
	if err != nil {
		panic(err)
	}

	log.Println("Rotterdam > CAAS > adapter [init] List of available QoS elements ...")
	for i := range cfg.QosTemplates {
		log.Println("Rotterdam > CAAS > adapter [init] > [" + strconv.Itoa(i) + "] " + cfg.QosTemplates[i].GuaranteeName)
	}
}

///////////////////////////////////////////////////////////////////////////////

/*
errorMessage
*/
func errorMessage(w http.ResponseWriter, err error, method string) {
	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "error",
		Method:      method,
		Message:     err.Error(),
		CaaSVersion: cfg.Config.CaaSVersion})
}

/*
okTaskMessage
*/
func okTaskMessage(w http.ResponseWriter, message string, method string, dbTask structs.DB_TASK) {
	json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
		Resp:        "ok",
		Method:      method,
		Message:     message,
		ID:          dbTask.TaskDefinition.ID, //classTask.Id,
		CaaSVersion: cfg.Config.CaaSVersion,
		URL:         dbTask.Url,
		Task:        dbTask})
}

///////////////////////////////////////////////////////////////////////////////

/*
InitializeAdapter initialization function
*/
func InitializeAdapter() {
	log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Initializing adapter for '" + cfg.Config.Clusters[0].Type + "' ...")
	if cfg.Config.Clusters[0].Type == constants.TypeOpenshift {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Openshift adapter")
		a = ops_adapter.OpenshiftAdapter{}
	} else if cfg.Config.Clusters[0].Type == constants.TypeKubernetes || cfg.Config.Clusters[0].Type == constants.TypeMicroK8s {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Kubernetes adapter")
		a = k8s_adapter.KubernetesAdapter{}
	} else {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] ERROR No adapter selected")
	}
}

/*
GetAdapter initialization function
*/
func GetAdapter(idCluster string) (adapters.Adapter, error) {
	t := common.GetClusterType(idCluster)

	log.Println("Rotterdam > CAAS > adapter [GetAdapter] Getting adapter for '" + t + "' ...")
	if t == constants.TypeOpenshift {
		log.Println("Rotterdam > CAAS > adapter [GetAdapter] Using Openshift adapter")
		return ops_adapter.OpenshiftAdapter{}, nil
	} else if t == constants.TypeKubernetes || t == constants.TypeMicroK8s {
		log.Println("Rotterdam > CAAS > adapter [GetAdapter] Using Kubernetes adapter")
		return k8s_adapter.KubernetesAdapter{}, nil
	}

	log.Println("Rotterdam > CAAS > adapter [GetAdapter] ERROR No adapter selected")
	return nil, errors.New("adapter: No adapter selected")
}

///////////////////////////////////////////////////////////////////////////////

/*
NotImplementedFunc Default Function for not implemented calls
*/
func NotImplementedFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [NotImplementedFunc] -Not implemented-")

	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "ok",
		Method:      "NotImplementedFunc",
		Message:     "not implemented",
		CaaSVersion: cfg.Config.CaaSVersion})
}

/*
HomePath returns home path
*/
func HomePath(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [HomePath] Returning SwaggerUI path ...")

	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "ok",
		Method:      "HomePath",
		Message:     "UI URL: /swaggerui/",
		CaaSVersion: cfg.Config.CaaSVersion})
}

/*
DeployTask Deploys a task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTask(w http.ResponseWriter, classTask *structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > adapter [deploy] New Default Task Deployment")

	adpt, err := GetAdapter(classTask.Cluster)
	if err != nil {
		errorMessage(w, err, "deploy")
	}

	log.Println("Rotterdam > CAAS > adapter [deploy] (1.) Deploying task [" + classTask.Name + "] ...")
	_, err = adpt.DeployTask(*classTask)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [deploy] (2.) Creating and starting SLA(s) ...")
		err = sla.CreateStartSLA(*classTask)
		if err != nil {
			log.Error("Rotterdam > CAAS > adapter [deploy] ERROR when creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := db.ReadTaskValue(classTask.ID)
		if err == nil {
			okTaskMessage(w, "Task deployed", "deploy", *dbTask)
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "deploy",
				Message:     "Task deployed, but SLA not started",
				ID:          classTask.ID,
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		errorMessage(w, err, "deploy")
	}
}

/*
sendDataToPrometheusPushgateway data (workers) to prometheus via pushgateway
*/
func sendDataToPrometheusPushgateway(idCluster string, classTask structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] Sending data to Pushgateway (Prometheus) ... ")

	cl, err := db_imec.GetInfrByID(idCluster) // ([]DB_INFRASTRUCTURE_CLUSTER, error)
	if err == nil {
		pushgatewayURL := cl[0].PrometheusPushgatewayEndPoint

		log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] Sending metrics to prometheus pushgateway [" + pushgatewayURL + "] ...")

		s := `
		workers_` + classTask.ID + ` ` + strconv.Itoa(classTask.Replicas) + `
		`

		ioreader := strings.NewReader(s)

		//log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] POST " + pushgatewayURL + "/metrics/job/sla/instance/violations")
		//res, err := http.Post(pushgatewayURL+"/metrics/job/sla/instance/violations", "binary/octet-stream", ioreader)
		log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] POST " + pushgatewayURL + "/metrics/job/compss")
		res, err := http.Post(pushgatewayURL+"/metrics/job/compss", "binary/octet-stream", ioreader)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] Error (1): " + err.Error())
		}
		defer res.Body.Close()
		message, _ := ioutil.ReadAll(res.Body)
		log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] Response: " + string(message))
	} else {
		log.Println("Rotterdam > CAAS > adapter [sendDataToPrometheusPushgateway] Error (2): " + err.Error())
	}
}

/*
DeployTaskCOMPSs Deploys a COMPSs task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTaskCOMPSs(w http.ResponseWriter, classTask *structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] New COMPSs Task Deployment")

	adpt, err := GetAdapter(classTask.Cluster)
	if err != nil {
		errorMessage(w, err, "deployCOMPSs")
	}

	// call to predictive SLA
	log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (0.) Call to predictive SLA ...")
	classTask, _ = sla.Predict(classTask)

	// start deployment
	log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (1.) Deploying task " + classTask.Name + " ...")
	_, err = adpt.DeployTaskCompss(*classTask)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (2.) Creating and starting SLA ...")
		err = sla.CreateStartCOMPSsSLA(*classTask)
		if err != nil {
			log.Error("Rotterdam > CAAS > adapter [deployCOMPSs] ERROR creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := db.ReadTaskValue(classTask.ID)
		if err == nil {
			// send "workers" information to Prometheus / Pushgateway
			sendDataToPrometheusPushgateway(classTask.Cluster, *classTask)

			// response
			okTaskMessage(w, "Task deployed", "deployCOMPSs", *dbTask)
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "deployCOMPSs",
				Message:     "Task deployed, but SLA not started",
				ID:          classTask.ID,
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		errorMessage(w, err, "deployCOMPSs")
	}
}

/*
DeployRotterdamTask Deploys a task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployRotterdamTask(w http.ResponseWriter, classTask *structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] New Task Deployment (undefined type)")

	// Task (default, compss)
	classTaskStr, _ := structs.CommClassStructToString(*classTask)
	log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Parsed CLASS TASK object (string): " + classTaskStr)

	// Check cluster
	log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Checking cluster status ... ")
	status := adapt_engine.GetClusterStatus(classTask.Cluster)
	if status != "ok" {
		log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Looking for alternative clusters ... ")
		// change cluster if possible
		idNewCluster := adapt_engine.GetAvailableCluster()
		if idNewCluster != "Not-Found" {
			classTask.Cluster = idNewCluster
			log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] New cluster selected for task deployment: " + idNewCluster)
		} else {
			log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] No alternative cluster found. Using cluster " + classTask.Cluster)
		}
	}

	if classTask.Type == constants.TypeTaskCOMPSs {
		log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] COMPSs task ...")
		DeployTaskCOMPSs(w, classTask)
	} else {
		log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Default task ...")
		DeployTask(w, classTask)
	}
}

/*
RemoveRotterdamTask Deletes a task
*/
func RemoveRotterdamTask(w http.ResponseWriter, dbTask *structs.DB_TASK) {
	// Task (default, compss)
	log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (1.) Getting adapter for cluster [" + dbTask.ClusterId + "]  ...")
	adpt, err := GetAdapter(dbTask.ClusterId)
	if err != nil {
		errorMessage(w, err, "RemoveRotterdamTask")
	} else {
		log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (2.) Removing task " + dbTask.Id + " ...")
		resp, agreementID, err := adpt.RemoveTask(*dbTask)
		if err == nil {
			if agreementID == constants.SLANotDefined {
				log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (3.) SLA(s) not defined")
			} else {
				log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (3.) Stopping and terminating SLA(s) ...")
				sla.StopTerminateSLA(dbTask.Id)
			}

			log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (4.) Deleting DBTask ...")
			db.DeleteTask(dbTask.Id)

			json.NewEncoder(w).Encode(structs.ResponseCaaS{
				Resp:        "ok",
				Method:      "RemoveRotterdamTask",
				Message:     "Task removed",
				CaaSVersion: cfg.Config.CaaSVersion,
				Content:     resp})
		} else {
			errorMessage(w, err, "RemoveRotterdamTask")
		}
	}
}

/*
RemoveTask Deletes a task
*/
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### Task termination (DEPRECATED)")

	log.Println("Rotterdam > CAAS > adapter [RemoveTask] (1.) Reading input params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [RemoveTask] dock=" + params["dock"] + ", id=" + params["id"])

	log.Println("Rotterdam > CAAS > adapter [RemoveTask] (2.) Getting task [" + params["id"] + "] from DB ...")
	dbTask, err := db.ReadTaskValue(params["id"])
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [RemoveTask] (2.) Removing task " + params["id"] + " ...")
		resp, agreementID, err := a.RemoveTask(*dbTask)
		if err == nil {
			if agreementID == constants.SLANotDefined {
				log.Println("Rotterdam > CAAS > adapter [RemoveTask] (3.) SLA(s) not defined")
			} else {
				log.Println("Rotterdam > CAAS > adapter [RemoveTask] (3.) Stopping and terminating SLA(s) ...")
				sla.StopTerminateSLA(dbTask.Id)
			}

			log.Println("Rotterdam > CAAS > adapter [RemoveTask] (4.) Deleting DBTask ...")
			db.DeleteTask(dbTask.Id)

			json.NewEncoder(w).Encode(structs.ResponseCaaS{
				Resp:        "ok",
				Method:      "RemoveTask",
				Message:     "Task removed",
				CaaSVersion: cfg.Config.CaaSVersion,
				Content:     resp})
		} else {
			errorMessage(w, err, "RemoveTask")
		}
	} else {
		errorMessage(w, err, "RemoveTask")
	}
}

/*
GetAllTasks Gets all tasks
*/
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/tasks")

	log.Println("Rotterdam > CAAS > adapter [GetAllTasks] Getting all tasks ...")
	resp, err := adpcommon.GetAllTasks()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks retrieved"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasks{
			Resp:        "ok",
			Method:      "GetAllTasks",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetAllTasks")
	}
}

/*
GetAllTasksQoS Gets all tasks QoS
*/
func GetAllTasksQoS(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetAllTasksQoS] Getting QoS tasks ...")

	resp, err := adpcommon.GetAllTasksQoS()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks qos retrieved"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasksQoS{
			Resp:        "ok",
			Method:      "GetAllTasksQoS",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			TasksQoS:    resp})
	} else {
		errorMessage(w, err, "GetAllTasksQoS")
	}
}

/*
LoadQoSDefinitions Loads QoS definitions list
*/
func LoadQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/qos/definitions")

	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Loading QoS definitions list ...")
	res, err := structs.StructCheckClasQoSTemplateList(r)

	if err == nil {
		cfg.QosTemplates = *res

		msg := strconv.Itoa(len(*res)) + " tasks qos loaded"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasksQoS{
			Resp:        "ok",
			Method:      "LoadQoSDefinitions",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		errorMessage(w, err, "LoadQoSDefinitions")
	}
}

/*
GetQoSDefinitions Gets QoS defintions
*/
func GetQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/qos/definitions")

	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Retrieving QoS definitions list ...")
	json.NewEncoder(w).Encode(structs.ResponseQoSDefinitions{
		Resp:        "ok",
		Method:      "GetQoSDefinitions",
		Message:     "QoS definition retrieved",
		CaaSVersion: cfg.Config.CaaSVersion,
		QoSDefs:     cfg.QosTemplates})
}

/*
GetQoSDefinition Gets QoS defintion by name
*/
func GetQoSDefinition(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/qos/definitions/{name}")

	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] name=" + params["name"])

	qosTemplate, found := sla.GetQoSElem(params["name"])
	if found {
		json.NewEncoder(w).Encode(structs.ResponseQoSDefinition{
			Resp:        "ok",
			Method:      "GetQoSDefinition",
			Message:     "QoS definition retrieved",
			CaaSVersion: cfg.Config.CaaSVersion,
			QoSDef:      qosTemplate})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetQoSDefinition",
			Message:     "QoS definition not found. Check logs",
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

/*
GetDockTasks Gets all tasks from Dock
*/
func GetDockTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] Reading input params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] dock=" + params["dock"])

	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] (1.) Getting all tasks from dock (" + params["dock"] + ") ...")
	resp, err := adpcommon.GetDockTasks(params["dock"])
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks retrieved [dock=" + params["dock"] + "]"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasks{
			Resp:        "ok",
			Method:      "GetDockTasks",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetDockTasks")
	}
}

/*
GetTask Gets a task
*/
func GetTask(w http.ResponseWriter, r *http.Request) {
	log.Trace("####################################################################################")
	log.Trace("### GET /api/v1/tasks/{id}")

	log.Trace("Rotterdam > CAAS > adapter [GetTask] Reading input params ...")
	params := mux.Vars(r)
	log.Trace("Rotterdam > CAAS > adapter [GetTask] dock=" + params["dock"] + ", id=" + params["id"])

	log.Trace("Rotterdam > CAAS > adapter [GetTask] (1.) Getting task " + params["id"] + " ...")
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		okTaskMessage(w, "Task retrieved", "GetTask", dbTask)
	} else {
		errorMessage(w, err, "GetTask")
	}
}

/*
GetTaskAllInfo Gets a task including deployment info
*/
func GetTaskAllInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/tasks/{id}/all")

	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] id=" + params["id"])

	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] (1.) Getting task " + params["id"] + " ...")
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		adapter, err := GetAdapter(dbTask.ClusterId)
		if err == nil {
			adapter.GetTaskAllInfo(params["id"])
			okTaskMessage(w, "Task retrieved", "GetTaskAllInfo", dbTask)
		}
	}
	errorMessage(w, err, "GetTaskAllInfo")
}

/*
GetConfig Gets configuration
*/
func GetConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetConfig] Getting configuration ...")
	resp, err := a.GetConfig()
	if err == nil {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "GetConfig",
			Message:     "Configuration retrieved",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     resp})
	} else {
		errorMessage(w, err, "GetConfig")
	}
}

/*
GetVersion Gets Rotterdam version
*/
func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetVersion] Getting Rotterdam components versions ...")
	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:               "ok",
		Method:             "GetVersion",
		Message:            "Rotterdam components versions retrieved",
		CaaSVersion:        cfg.Config.CaaSVersion,
		RulesEngineVersion: cfg.Config.RulesEngineVersion,
		RestAPIVersion:     cfg.Config.RestApiVersion,
		SLALiteVersion:     cfg.Config.SLALiteVersion,
		IMECVersion:        cfg.Config.IMECVersion})
}

/*
ScaleUpDown ...
*/
func ScaleUpDown(dbtask structs.DB_TASK, replicas int) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling out / in the task ...")

	resp, err := a.ScaleUpDown(dbtask, replicas)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling up / down finalized: " + resp)
	} else {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] ERROR Scaling up / down task " + dbtask.TaskDefinition.Name + " [" + dbtask.Id + "] ...")
	}
}
