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

package caas

import (
	adapters "atos/rotterdam/caas/adapters"
	adpcommon "atos/rotterdam/caas/adapters/common"
	k8s_adapter "atos/rotterdam/caas/adapters/kubernetes"
	ops_adapter "atos/rotterdam/caas/adapters/openshift"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	sla "atos/rotterdam/caas/sla"
	cfg "atos/rotterdam/config"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

// adapter
var a adapters.Adapter

// selected cluster
var selected_cluster int

// init
func init() {
	log.Println("Rotterdam > CAAS > adapter [init] Selected cluster index set to '0' ...")
	selected_cluster = 0

	qos_templates_path := "./config/qos_templates.json"
	log.Println("Rotterdam > CAAS > adapter [init] Reading content of qos_templates.json file [" + qos_templates_path + "] ...")

	file, err := os.Open(qos_templates_path)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.QosTemplates)
	if err != nil {
		panic(err)
	}

	log.Println("Rotterdam > CAAS > adapter [init] List of available QoS elements ...")
	for i, _ := range cfg.QosTemplates {
		log.Println("Rotterdam > CAAS > adapter [init] > [" + strconv.Itoa(i) + "] " + cfg.QosTemplates[i].GuaranteeName)
	}
}

// initialization function
func InitializeAdapter() {
	log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Initializing adapter for '" + cfg.Config.Clusters[0].Mode + "' ...")
	if cfg.Config.Clusters[0].Mode == "Openshift" {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Openshift adapter")
		a = ops_adapter.OpenshiftAdapter{}
	} else if cfg.Config.Clusters[0].Mode == "Kubernetes" {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Kubernetes adapter")
		a = k8s_adapter.KubernetesAdapter{}
	} else {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] ERROR No adapter selected")
	}
}

///////////////////////////////////////////////////////////////////////////////

// Deploy a task (k8s: deployment & service & volumes ...)
// Input example: view CLASS_TASK in structs.go
func DeployTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [DeployTask] Validating json ...")

	class_task, err := common.StructCheckClassTask(r)
	if err == nil {
		// check if id already exists
		// TODO
		// Create uuid
		alphabet := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx"
		task_id := shortuuid.NewWithAlphabet(alphabet)
		class_task.Id = class_task.Name + "-" + task_id

		log.Println("Rotterdam > CAAS > adapter [DeployTask] Task with name " + class_task.Name + " received")
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTask",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}

	log.Println("Rotterdam > CAAS > adapter [DeployTask] Deploying task " + class_task.Name + " ...")
	_, err = a.DeployTask(selected_cluster, "class", *class_task)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [DeployTask] Creating and starting SLA ...")
		err = sla.CreateStartSLA(*class_task)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [DeployTask] ERROR when creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := common.ReadTaskValue(class_task.Name)
		if err == nil {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "DeployTask",
				Message:     "Task deployed",
				CaaSVersion: cfg.Config.CaaSVersion,
				URL:         dbTask.Url,
				Task:        *dbTask})
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "DeployTask",
				Message:     "Task deployed, but SLA not started",
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTask",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

///////////////////////////////////////////////////////////////////////////////

// Deploy a COMPSs task (k8s: deployment & service & volumes ...)
// Input example: view CLASS_TASK in structs.go
func DeployTaskCOMPSs(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] Validating json ...")

	class_task, err := common.StructCheckClassTask(r)
	if err == nil {
		// Create uuid
		alphabet := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx"
		task_id := shortuuid.NewWithAlphabet(alphabet)
		class_task.Id = class_task.Name + "-" + task_id

		log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] Task with name " + class_task.Name + " received")
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTaskCOMPSs",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}

	log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] Deploying task " + class_task.Name + " ...")
	_, err = a.DeployTaskCompss(selected_cluster, "class", *class_task)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] Creating and starting SLA ...")
		err = sla.CreateStartSLA(*class_task)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] ERROR when creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := common.ReadTaskValue(class_task.Name)
		if err == nil {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "DeployTaskCOMPSs",
				Message:     "Task deployed",
				CaaSVersion: cfg.Config.CaaSVersion,
				URL:         dbTask.Url,
				Task:        *dbTask})
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "DeployTaskCOMPSs",
				Message:     "Task deployed, but SLA not started",
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTaskCOMPSs",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Deletes a task
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [RemoveTask] Reading params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [RemoveTask] dock..." + params["dock"])
	log.Println("Rotterdam > CAAS > adapter [RemoveTask] name..." + params["name"])

	log.Println("Rotterdam > CAAS > adapter [RemoveTask] Removing task " + params["name"] + " ...")
	resp, err := a.RemoveTask(selected_cluster, params["dock"], params["name"])
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [RemoveTask] Stopping and terminating SLA ...")
		sla.StopTerminateSLA(params["name"])

		log.Println("Rotterdam > CAAS > adapter [RemoveTask] Deleting DBTask ...")
		common.DBDeleteTask(params["name"])

		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "RemoveTask",
			Message:     "Task removed",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     resp})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "RemoveTask",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Get all tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetAllTasks] Getting tasks ...")

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
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetAllTasks",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Get all tasks QoS
func GetAllTasksQoS(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetAllTasksQoS] Getting tasks ...")

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
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetAllTasksQoS",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Load QoS definitions list
func LoadQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Loading QoS definitions list ...")

	res, err := common.StructCheckClasQoSTemplateList(r)

	if err == nil {
		cfg.QosTemplates = *res

		msg := strconv.Itoa(len(*res)) + " tasks qos loaded"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasksQoS{
			Resp:        "ok",
			Method:      "LoadQoSDefinitions",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "LoadQoSDefinitions",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Get QoS defintions
func GetQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Retrieving QoS definitions list ...")

	json.NewEncoder(w).Encode(structs.ResponseQoSDefinitions{
		Resp:        "ok",
		Method:      "GetQoSDefinitions",
		Message:     "QoS definition retrieved",
		CaaSVersion: cfg.Config.CaaSVersion,
		QoSDefs:     cfg.QosTemplates})
}

// Get QoS defintion by name
func GetQoSDefinition(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] Reading params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] Getting QoS definition [" + params["name"] + "] ...")

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

// Get all tasks from Dock
func GetDockTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] Reading params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] Getting tasks from dock [" + params["dock"] + "] ...")

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
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetDockTasks",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Gets a task
func GetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetTask] Reading params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetTask] dock..." + params["dock"])
	log.Println("Rotterdam > CAAS > adapter [GetTask] name..." + params["name"])

	log.Println("Rotterdam > CAAS > adapter [GetTask] Getting task " + params["name"] + " ...")
	resp, err := a.GetTask(selected_cluster, params["dock"], params["name"])
	if err == nil {
		json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
			Resp:        "ok",
			Method:      "GetTask",
			Message:     "Task retrieved",
			CaaSVersion: cfg.Config.CaaSVersion,
			Task:        resp})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetTask",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Get configuration
func GetConfig(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetConfig",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

// Get Rotterdam version
func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [GetVersion] Getting Rotterdam components versions ...")
	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:               "ok",
		Method:             "GetVersion",
		Message:            "Rotterdam components versions retrieved",
		CaaSVersion:        cfg.Config.CaaSVersion,
		RulesEngineVersion: cfg.Config.RulesEngineVersion,
		RestApiVersion:     cfg.Config.RestApiVersion,
		SLALiteVersion:     cfg.Config.SLALiteVersion})
}

// ScaleUpDown
func ScaleUpDown(dbtask structs.DB_TASK, replicas int) {
	log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling up / down ...")
	dbtask.TaskDefinition.Dock = "class"
	resp, err := a.ScaleUpDown(selected_cluster, dbtask, replicas)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling up / down finalized: " + resp)
	} else {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] ERROR Scaling up / down task " + dbtask.TaskDefinition.Name)
	}
}
