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
// Created on 11 June 2019
// @author: Roi Sucasas - ATOS
//

package impl

import (
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"errors"
	"log"
	"strconv"
)

// delK8sDeployment: k8s: deployment
func delK8sDeployment(namespace string, name string) (string, error) {
	// CALL to Kubernetes API to delete a deployment
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sDeployment] Deleting deployment from K8s cluster ...")
	// map[string]interface{}, error
	status, _, err := common.HttpDELETE_GenericStruct(
		cfg.Config.Clusters[0].KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sDeployment] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// delK8sService: k8s: service
func delK8sService(namespace string, name string) (string, error) {
	// CALL to Kubernetes API to delete a service
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sService] Deleting service from K8s cluster ...")
	// map[string]interface{}, error
	status, _, err := common.HttpDELETE_GenericStruct(
		cfg.Config.Clusters[0].KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/services/serv-" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sService] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [delK8sService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// removeDefaultTask: Deletes a task
func removeDefaultTask(namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := delK8sDeployment(namespace, name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICE /////
		status, err := delK8sService(namespace, name)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (2)", err)
			return "", err
		} else if status == "200" {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] Task removed with success")
			return "ok", nil
		}
	}

	err = errors.New("Task termination failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (4)", err)
	return "", err
}

// removeCOMPSsTask: Deletes a COMPSs task
func removeCOMPSsTask(namespace string, name string, dbtask structs.DB_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] Deleting COMPSs task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := delK8sDeployment(namespace, name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICES /////
		for _, pod := range dbtask.Pods {
			// DB_TASK_POD
			status, err := delK8sService(namespace, pod.Name)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (2)", err)
				return "", err
			} else if status == "200" {
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] Task removed with success")
				return "ok", nil
			}
		}
	}

	err = errors.New("Task termination failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (4)", err)
	return "", err
}

// RemoveTask: Deletes a task
func RemoveTask(namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// get type of task
	dbTask, err := common.ReadTaskValue(name) // (*structs.DB_TASK, error)
	if err == nil {
		// remove task
		if dbTask.Type == "default" {
			removeDefaultTask(namespace, name)
		} else if dbTask.Type == "compss" {
			removeCOMPSsTask(namespace, name, *dbTask)
		} else {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] WARNING type of task is not defined: " + dbTask.Type)
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] ERROR", err)
	return "", err
}
