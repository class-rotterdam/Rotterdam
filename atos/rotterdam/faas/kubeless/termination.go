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

package kubeless

import (
	urls "atos/rotterdam/caas/adapters"
	"atos/rotterdam/caas/common"
	imec_db "atos/rotterdam/database/imec"
	"atos/rotterdam/globals/structs"
	"errors"
	"log"
	"strconv"
)

/*
kubelessDeployment Function termination
Examples:

	curl -X DELETE http://10.0.2.15:8001/apis/kubeless.io/v1beta1/namespaces/default/functions/get-python
*/
func kubelessTermination(namespace string, id string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, error) {
	// CALL to Kubernetes API to delete a deployment
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessTermination] Deleting function from K8s cluster ...")

	status, _, err := common.HTTPDELETEStruct(
		"http://"+urls.GetHostIP(cluster)+":8001/apis/kubeless.io/v1beta1/namespaces/"+namespace+"/functions/"+id,
		sec)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessTermination] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessTermination] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// removeDefaultFunc deletes a function
func removeDefaultFunc(namespace string, id string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [removeDefaultFunc] Deleting Function [" + id + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := kubelessTermination(namespace, id, cluster, false)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [removeDefaultFunc] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [removeDefaultFunc] Function removed with success")
		return "ok", nil
	}

	err = errors.New("Task termination failed. Internal status = [" + status + "]")
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [removeDefaultFunc] ERROR (2)", err)
	return "", err
}

/*
Remove deletes a function
*/
func Remove(dbTask structs.DB_TASK) (string, string, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Remove] Deleting function [" + dbTask.Id + "] ...")

	clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
	function := dbTask.FunctionDefinition
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Remove] cluster id = " + dbTask.ClusterId + ", dock = " + function.Dock + "")

	// remove task
	if dbTask.Type == structs.DB_TASK_TYPE_FUNCTION {
		res, err := removeDefaultFunc(function.Dock, dbTask.Id, clusterInfr)
		return res, dbTask.AgreementId, err
	}

	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Remove] WARNING type of task is not defined: " + dbTask.Type)
	return "", "", nil
}
