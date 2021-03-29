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

package faas

import (
	adapters "atos/rotterdam/caas/adapters"
	adpcommon "atos/rotterdam/caas/adapters/common"
	cfg "atos/rotterdam/config"
	db "atos/rotterdam/database/caas"
	kubeless "atos/rotterdam/faas/kubeless"
	constants "atos/rotterdam/globals/constants"
	structs "atos/rotterdam/globals/structs"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// adapter
var a adapters.Adapter

///////////////////////////////////////////////////////////////////////////////

// init
func init() {
	log.Println("Rotterdam > FAAS > adapter [init] Selected cluster index set to '0' ...")

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
		FaaSVersion: cfg.Config.FaaSVersion})
}

/*
okFTaskMessage
*/
func okFTaskMessage(w http.ResponseWriter, message string, method string, dbTask structs.DB_TASK) {
	json.NewEncoder(w).Encode(structs.ResponseFaaSTask{
		Resp:        "ok",
		Method:      method,
		Message:     message,
		ID:          dbTask.Id,
		FaaSVersion: cfg.Config.FaaSVersion,
		Task:        dbTask})
}

/*
okFTaskCallMessage
*/
func okFTaskCallMessage(w http.ResponseWriter, message string, method string, id string, result string) {
	json.NewEncoder(w).Encode(structs.ResponseFaaSTask{
		Resp:        "ok",
		Method:      method,
		Message:     message,
		ID:          id,
		Result:      result,
		FaaSVersion: cfg.Config.FaaSVersion})
}

///////////////////////////////////////////////////////////////////////////////

/*
DeployFunction Deploys a function
*/
func DeployFunction(w http.ResponseWriter, classFTask *structs.CLASS_FUNCTION_TASK) { //w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/functions    <<Function Deployment>>")

	dbTask, err := kubeless.Deploy(w, classFTask)
	if err == nil {
		okFTaskMessage(w, "Function deployed", "DeployFunction", *dbTask)
	} else {
		errorMessage(w, err, "DeployFunction")
	}
}

/*
CallFunction Calls a function
*/
func CallFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/functions/{id}    <<Call to a function>>")

	log.Println("Rotterdam > FAAS > adapter [CallFunction] (1.) Reading input params ...")
	params := mux.Vars(r)

	log.Println("Rotterdam > FAAS > adapter [CallFunction] (2.) Getting task/function [" + params["id"] + "] from DB ...")
	dbTask, err := db.ReadTaskValue(params["id"])
	if err == nil && dbTask.Type == constants.TypeFTaskDefault {
		classFTask := dbTask.FunctionDefinition
		if err != nil {
			errorMessage(w, err, "CallFunction")
		} else {
			res, resut, err := kubeless.Call(w, r, &dbTask.FunctionDefinition)
			if err == nil {
				okFTaskCallMessage(w, res, "CallFunction", classFTask.ID, resut)
			} else {
				errorMessage(w, err, "CallFunction")
			}
		}
	}

}

/*
RemoveFunction Deletes a function
*/
func RemoveFunction(w http.ResponseWriter, dbTask *structs.DB_TASK) {
	log.Println("####################################################################################")
	log.Println("### DELETE /api/v1/functions/{id}    <<Function termination>>")

	log.Println("Rotterdam > FAAS > adapter [RemoveFunction] (3.) Removing function " + dbTask.Id + " ...")
	resp, _, err := kubeless.Remove(*dbTask)
	if err == nil {
		// TODO
		/*if agreementID == constants.SLANotDefined {
			log.Println("Rotterdam > FAAS > adapter [RemoveFunction] (5.) SLA(s) not defined")
		} else {
			log.Println("Rotterdam > FAAS > adapter [RemoveFunction] (5.) Stopping and terminating SLA(s) ...")
			sla.StopTerminateSLA(dbTask.Id)
		}*/

		log.Println("Rotterdam > FAAS > adapter [RemoveFunction] (6.) Deleting DBTask object ...")
		db.DeleteTask(dbTask.Id)

		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "RemoveFunction",
			Message:     "Function removed",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     resp})
	} else {
		errorMessage(w, err, "RemoveFunction")
	}
}

/*
GetAllFunctions Gets all functions
*/
func GetAllFunctions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/functions")

	log.Println("Rotterdam > FAAS > adapter [GetAllFunctions] (1.) Getting all functions ...")
	// TODO
	resp, err := adpcommon.GetAllTasks()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " functions retrieved"
		json.NewEncoder(w).Encode(structs.ResponseFaaSTasks{
			Resp:        "ok",
			Method:      "GetAllFunctions",
			Message:     msg,
			FaaSVersion: cfg.Config.FaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetAllFunctions")
	}
}

/*
GetFunction Gets a function
*/
func GetFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/functions/{id}")

	log.Println("Rotterdam > FAAS > adapter [GetFunction] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > FAAS > adapter [GetFunction] id=" + params["id"])

	log.Println("Rotterdam > FAAS > adapter [GetFunction] (1.) Getting function " + params["id"] + " ...")
	// TODO
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		okFTaskMessage(w, "Function retrieved", "GetFunction", dbTask)
	} else {
		errorMessage(w, err, "GetFunction")
	}
}
