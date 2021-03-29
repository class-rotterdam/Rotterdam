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

package rest_api

import (
	"atos/rotterdam/caas/common"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	"atos/rotterdam/globals/constants"
	"atos/rotterdam/globals/structs"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"

	caas "atos/rotterdam/caas"
	db "atos/rotterdam/database/caas"
	faas "atos/rotterdam/faas"
)

/*
generateID generates an Identifier based on a random string and the current time in nanoseconds
*/
func generateID() string {
	// 1. ID generation
	id := shortuuid.NewWithAlphabet("0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx") //id := txt + shortuuid.New()
	id = strings.ToLower(id)                                                                     // Kubernetes doesnt allow uppercase names for deployments

	// 2. time in nanoseconds is appended to ID
	now := time.Now()      // current local time
	nsec := now.UnixNano() // number of nanoseconds since January 1, 1970 UTC
	id = id + strconv.FormatInt(nsec, 10)

	return id
}

/*
setClassStructDataDefaultValues
*/
func setClassStructDataDefaultValues(task *structs.CLASS_TASK, isCOMPs bool) {
	if len(task.Cluster) == 0 {
		task.Cluster = constants.MainClusterID
	}
	if len(task.Dock) == 0 {
		task.Dock = common.GetClusterDefaultDock(task.Cluster)
	}

	if isCOMPs == true {
		task.Type = constants.TypeTaskCOMPSs
	} else if len(task.Type) == 0 || (task.Type != constants.TypeTaskCOMPSs && task.Type != constants.TypeFTaskDefault) {
		task.Type = constants.TypeTaskDefault
	}
}

/*
setClassStructDataAll
*/
func setClassStructDataAll(task *structs.CLASS_TASK, isCOMPs bool) {
	task.ID = generateID()                         // ID
	task.Created = time.Now().Format(time.RFC3339) // current local time
	setClassStructDataDefaultValues(task, isCOMPs)
}

/*
validateJSONTask validates input (json) and generates a valid CLASS Task struct
*/
func validateJSONTask(r *http.Request) (*structs.CLASS_TASK, error) {
	log.Println(pathLOG + "[validateJsonTask] 'Duplicating' r.Body to use it multiple times with decoders ...")
	rBody := r.Body
	var buf bytes.Buffer
	tee := io.TeeReader(rBody, &buf)

	log.Println(pathLOG + "[validateJsonTask] Parsing default / old json definition ...")
	decoder := json.NewDecoder(tee)
	classTask, err := structs.StructCheckClassTask(decoder)
	if err == nil {
		setClassStructDataAll(classTask, false)
		log.Println(pathLOG + "[validateJsonTask] Task with name " + classTask.Name + " received. Task ID is [" + classTask.ID + "]. Type = " + classTask.Type)
		return classTask, nil
	} else {
		log.Println(pathLOG + "[validateJsonTask] Parsing new json definition (COMPSs) ...")

		decoder = json.NewDecoder(&buf)
		classCOMPSsTask, err := structs.StructCheckClassCOMPSsTask(decoder)
		if err == nil {
			classCOMPSsTask.ID = generateID() // ID
			classTask = structs.TransfCOMPSSTASKtoTASK(classCOMPSsTask)
			log.Println(pathLOG + "[validateJsonTask] Task with name " + classCOMPSsTask.Name + " received. Task ID is [" + classCOMPSsTask.ID + "]. Type = " + classTask.Type + ". Using 'COMPSs' JSON format.")

			if classTask.Type == constants.TypeFTaskDefault {
				setClassStructDataDefaultValues(classTask, false)
			} else {
				setClassStructDataDefaultValues(classTask, true)
			}

			return classTask, nil
		}
	}

	return nil, err
}

///////////////////////////////////////////////////////////////////////////////

/*
Deploy Deploys a task
*/
func Deploy(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println(pathLOG + "[Deploy] Task/Function deployment ...")

	log.Println(pathLOG + "[Deploy] Validating json and extracting Task information...")
	classTask, err := validateJSONTask(r)
	if err == nil {
		if classTask.Type == structs.DB_TASK_TYPE_FUNCTION {
			// Task is a function
			log.Println(pathLOG + "[Deploy] JSON input should be a function")
			log.Println(pathLOG + "[Deploy] Getting function information (class task to class function transformation) ...")
			classFunctionTask := structs.TransfTASKtoFunctionTASK(classTask)

			faas.DeployFunction(w, classFunctionTask)
		} else {
			// Task (default, compss)
			log.Println(pathLOG + "[Deploy] JSON input is a task")
			caas.DeployRotterdamTask(w, classTask)
		}
	} else {
		log.Println(pathLOG + "[Deploy] ERROR JSON not valid")
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "Deploy",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

/*
DeployTask Deploys a task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println(pathLOG + "[Deploy] Default Task deployment ...")

	classTask, err := validateJSONTask(r)
	if err != nil {
		log.Println(pathLOG + "[Deploy] ERROR JSON not valid")
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTask",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		caas.DeployTask(w, classTask)
	}
}

/*
DeployTaskCOMPSs Deploys a COMPSs task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTaskCOMPSs(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println(pathLOG + "[Deploy] COMPSs Task deployment ...")

	classTask, err := validateJSONTask(r)
	if err != nil {
		log.Println(pathLOG + "[Deploy] ERROR JSON not valid")
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "DeployTaskCOMPSs",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		caas.DeployTaskCOMPSs(w, classTask)
	}
}

/*
Remove Removes a task
*/
func Remove(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println(pathLOG + "[Remove] Task/Function termination ...")

	log.Println(pathLOG + "[Remove] Reading input params ...")
	params := mux.Vars(r)

	log.Println(pathLOG + "[Remove] Getting task [" + params["id"] + "] from DB ...")
	dbTask, err := db.ReadTaskValue(params["id"])
	if err == nil {
		if dbTask.Type == structs.DB_TASK_TYPE_FUNCTION {
			log.Println(pathLOG + "[Remove] Input parameter is a function ID")
			faas.RemoveFunction(w, dbTask)
		} else {
			log.Println(pathLOG + "[Remove] Input parameter is a task ID")
			caas.RemoveRotterdamTask(w, dbTask)
		}
	} else {
		log.Println(pathLOG + "[Remove] ERROR")
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "Remove",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}
