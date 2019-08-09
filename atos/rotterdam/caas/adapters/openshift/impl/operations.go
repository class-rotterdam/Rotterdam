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
	"log"
)


// getK8sDeployment: k8s: get deployment info
func getK8sDeployment(cluster_index int, namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] Getting deployment from K8s cluster ...")
	// map[string]interface{}, error
	_, data, err := common.HttpGET_String(
		cfg.Config.Clusters[cluster_index].KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] RESPONSE: " + data)

	return data, nil
}

// GetTask: Gets a deployment
func GetTask(cluster_index int, namespace string, name string) (structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetTask] Getting deployment [" + name + "] from [" + namespace + "] ...")

	// get task
	dbTask, err := common.ReadTaskValue(name)
	if err == nil {
		// get deployment information
		deployment_data, err := getK8sDeployment(cluster_index, namespace, name)
		if err == nil {
			dbTask.Deployment = deployment_data
			dbTask.Deployment = ""
			return *dbTask, nil
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetTask] ERROR", err)
	return *dbTask, err
}

// GetConfig: k8s configuration
func GetConfig() (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] Getting configuration from K8s cluster ...")

	_, data, err := common.HttpGET_String(cfg.Config.Clusters[0].KubernetesEndPoint + "/api")
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] RESPONSE: " + data)

	return data, err
}
