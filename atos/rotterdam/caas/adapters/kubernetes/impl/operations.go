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
	common "atos/rotterdam/caas/common"
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	db "atos/rotterdam/database/caas"
	imec_db "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
)

// getK8sDeployment: k8s: get deployment info
func getK8sDeployment(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, name string) (string, error) {
	log.Println(pathLOG + "Operations [getK8sDeployment] Getting deployment from K8s cluster ...")
	// map[string]interface{}, error
	_, data, err := common.HTTPGETString(
		urls.GetPathKubernetesDeployment(cluster, namespace, name),
		false)
	if err != nil {
		log.Error(pathLOG+"Operations [getK8sDeployment] ERROR", err)
		return "", err
	}
	log.Debug(pathLOG + "Operations [getK8sDeployment] RESPONSE: " + data)

	return data, nil
}

/*
GetTaskAllInfo Returns a task (including deployment info)
*/
func GetTaskAllInfo(idTask string) (structs.DB_TASK, error) {
	log.Println(pathLOG + "Operations [GetTaskAllInfo] Getting deployment from Task with id=" + idTask + " ...")

	// get task
	dbTask, err := db.ReadTaskValue(idTask)
	if err == nil {
		clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
		task := dbTask.TaskDefinition
		namespace := task.Dock

		// get deployment information
		deploymentData, err := getK8sDeployment(clusterInfr, namespace, idTask)
		if err == nil {
			dbTask.Deployment = deploymentData
			return *dbTask, nil
		}
	}

	log.Error(pathLOG+"Operations [GGetTaskAllInfoetTask] ERROR", err)
	return *dbTask, err
}

/*
GetConfig Get k8s configuration
*/
func GetConfig() (string, error) {
	log.Println(pathLOG + "Operations [GetConfig] Getting configuration from K8s cluster ...")

	_, data, err := common.HTTPGETString(cfg.Config.Clusters[0].KubernetesEndPoint+"/api", false)
	if err != nil {
		log.Error(pathLOG+"Operations [GetConfig] ERROR", err)
		return "", err
	}
	log.Println(pathLOG + "Operations [GetConfig] RESPONSE: " + data)

	return data, err
}
