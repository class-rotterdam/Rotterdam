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
func delK8sDeployment(cluster_index int, namespace string, name string) (string, error) {
	// CALL to Kubernetes API to delete a deployment
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sDeployment] Deleting deployment from K8s cluster ...")
	// map[string]interface{}, error
	status, _, err := common.HttpDELETE_GenericStruct(
		cfg.Config.Clusters[cluster_index].KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sDeployment] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// delK8sService: k8s: service
func delK8sService(cluster_index int, namespace string, name string) (string, error) {
	// CALL to Kubernetes API to delete a service
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sService] Deleting service from K8s cluster ...")
	// map[string]interface{}, error
	status, _, err := common.HttpDELETE_GenericStruct(
		cfg.Config.Clusters[cluster_index].KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/services/serv-" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sService] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// delK8sRoute: k8s: route
func delK8sRoute(cluster_index int, namespace string, name string) (string, error) {
	// CALL to Kubernetes API to delete a route
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] Deleting route from K8s cluster ...")
	// map[string]interface{}, error
	status, _, err := common.HttpDELETE_GenericStruct(
		cfg.Config.Clusters[cluster_index].OpenshiftEndPoint + "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes/route-" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [delK8sRoute] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// RemoveTask: Deletes a task
func RemoveTaskOLD(cluster_index int, namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := delK8sDeployment(cluster_index, namespace, name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. ROUTE /////
		status, err := delK8sRoute(cluster_index, namespace, name)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] ERROR (2)", err)
			return "", err
		} else if status == "200" {
			// 3. SERVICE /////
			status, err := delK8sService(cluster_index, namespace, name)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] ERROR (3)", err)
				return "", err
			} else if status == "200" {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] Task removed with success")
				return "ok", nil
			}
		}
	}

	err = errors.New("Task termination failed. Internal status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] ERROR (4)", err)
	return "", err
}

// removeDefaultTask: Deletes a task
func removeDefaultTask(cluster_index int, namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := delK8sDeployment(cluster_index, namespace, name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. ROUTE /////
		status, err := delK8sRoute(cluster_index, namespace, name)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeDefaultTask] ERROR (2)", err)
			return "", err
		} else if status == "200" {
			// 3. SERVICE /////
			status, err := delK8sService(cluster_index, namespace, name)
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

// removeCOMPSsTask: Deletes a COMPSs task
func removeCOMPSsTask(cluster_index int, namespace string, name string, dbtask structs.DB_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] Deleting COMPSs task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := delK8sDeployment(cluster_index, namespace, name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [removeCOMPSsTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICES /////
		for _, pod := range dbtask.Pods {
			// DB_TASK_POD
			status, err := delK8sService(cluster_index, namespace, pod.Name)
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

// RemoveTask: Deletes a task
func RemoveTask(cluster_index int, namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// get type of task
	dbTask, err := common.ReadTaskValue(name) // (*structs.DB_TASK, error)
	if err == nil {
		// remove task
		if dbTask.Type == structs.DB_TASK_TYPE_DEFAULT {
			removeDefaultTask(cluster_index, namespace, name)
		} else if dbTask.Type == structs.DB_TASK_TYPE_COMPSS {
			removeCOMPSsTask(cluster_index, namespace, name, *dbTask)
		} else {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Termination [RemoveTask] WARNING type of task is not defined: " + dbTask.Type)
		}
	}

	return "", err
}
