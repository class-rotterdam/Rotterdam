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

package impl

import (
	urls "atos/rotterdam/caas/adapters"
	adapt_common "atos/rotterdam/caas/adapters/common"
	common "atos/rotterdam/caas/common"
	imec_db "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
	"errors"
	"log"
	"strconv"
)

// delK8sRoute deletes a route
func delK8sRoute(namespace string, name string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] Deleting route from K8s cluster ...")

	status, _, err := common.HTTPDELETEStruct(
		urls.GetPathOpenshiftRoute(cluster, namespace, name),
		true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] ERROR", err)
		return strconv.Itoa(status), err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] RESPONSE: OK")
	return strconv.Itoa(status), nil
}

// removeDefaultTask deletes a task
func removeDefaultTask(namespace string, name string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := adapt_common.DelK8sDeployment(namespace, name, cluster, true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. ROUTE /////
		status, err := delK8sRoute(namespace, name, cluster)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (2)", err)
			return "", err
		} else if status == "200" {
			// 3. SERVICE /////
			status, err := adapt_common.DelK8sService(namespace, name, cluster, true)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (3)", err)
				return "", err
			} else if status == "200" {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] Task removed with success")
				return "ok", nil
			}
		}
	}

	err = errors.New("Task termination failed. Internal status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (4)", err)
	return "", err
}

// removeCOMPSsTask deletes a COMPSs task
func removeCOMPSsTask(namespace string, name string, dbtask structs.DB_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] Deleting COMPSs task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := adapt_common.DelK8sDeployment(namespace, name, cluster, true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICES /////
		for _, pod := range dbtask.Pods {
			// DB_TASK_POD
			status, err := adapt_common.DelK8sService(namespace, pod.Name, cluster, true)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] ERROR (2)", err)
				//return "", err
			} else if status == "200" {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] Task removed with success")
				//return "ok", nil
			}
		}
		return "ok", nil
	}

	err = errors.New("Task termination failed. Internal status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] ERROR (4)", err)
	return "", err
}

/*
RemoveTask deletes a task
*/
func RemoveTask(dbTask structs.DB_TASK) (string, string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] Deleting task [" + dbTask.Id + "] ...")

	clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
	task := dbTask.TaskDefinition
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] cluster id = " + dbTask.ClusterId + ", dock = " + task.Dock + "")

	// remove task
	if dbTask.Type == structs.DB_TASK_TYPE_DEFAULT {
		res, err := removeDefaultTask(task.Dock, dbTask.Id, clusterInfr)
		return res, dbTask.AgreementId, err
	} else if dbTask.Type == structs.DB_TASK_TYPE_COMPSS {
		res, err := removeCOMPSsTask(task.Dock, dbTask.Id, dbTask, clusterInfr)
		return res, dbTask.AgreementId, err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] WARNING type of task is not defined: " + dbTask.Type)
	return "", "", nil
}
