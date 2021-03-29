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
	db "atos/rotterdam/database/caas"
	imec_db "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
	"errors"
	"strconv"
)

// getK8sScale: k8s: get scale info
func getK8sScale(task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (*structs.K8S_SCALE, error) {
	log.Println(pathLOG + "Scaling [getK8sScale] Getting scaling info from task [" + task.Name + "] ...")
	_, data, err := common.HTTPGETString(
		urls.GetPathKubernetesScaleDeployment(cluster, task.Dock, task.ID),
		true)
	if err != nil {
		log.Error(pathLOG+"Scaling [getK8sScale] ERROR", err)
		return nil, err
	}
	log.Println(pathLOG + "Scaling [getK8sScale] RESPONSE: " + data)

	return structs.CommStringToK8S_SCALE(data)
}

// updateK8sScale: k8s: update scale info => scale up / down
func updateK8sScale(task structs.CLASS_TASK, scaleObj structs.K8S_SCALE, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println(pathLOG + "Scaling [updateK8sScale] Updating scaling info from task [" + task.Name + "] ...")
	status, _, err := common.HTTPPUTStruct(
		urls.GetPathKubernetesScaleDeployment(cluster, task.Dock, task.ID),
		true,
		scaleObj)
	if err != nil {
		log.Error(pathLOG+"Scaling [updateK8sScale] ERROR", err)
		return "", err
	}
	log.Println(pathLOG + "Scaling [updateK8sScale] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
ScaleUpDown Scale out / in a task
*/
func ScaleUpDown(dbTask structs.DB_TASK, replicas int) (string, error) {
	log.Println(pathLOG + "Scaling [ScaleUpDown] Scaling up/down task [" + dbTask.TaskDefinition.Name + "] ...")

	clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
	task := dbTask.TaskDefinition
	log.Println(pathLOG + "Scaling [ScaleUpDown] cluster=" + clusterInfr.ID)

	scaleObj, err := getK8sScale(task, clusterInfr)
	if err != nil {
		log.Error(pathLOG+"Scaling [ScaleUpDown] ERROR (1) ", err)
		return "", err
	}

	scaleObj.Spec.Replicas = replicas
	status, err := updateK8sScale(task, *scaleObj, clusterInfr)
	if err == nil {
		if err == nil {
			if dbTask.Type == structs.DB_TASK_TYPE_DEFAULT {
				// update task information: replicas
				dbTask.Replicas = replicas
				task.Replicas = replicas
				err = db.SetTaskValue(dbTask.Id, dbTask)
			} else if dbTask.Type == structs.DB_TASK_TYPE_COMPSS {
				// create / remove services nad update task info
				go func() {
					log.Println(pathLOG + "Scaling [ScaleUpDown] Starting background tasks...")
					dbTask = CompssScalingUpdateServices(clusterInfr, dbTask, replicas)
					_ = db.SetTaskValue(dbTask.Id, dbTask)
					log.Println(pathLOG + "Scaling [ScaleUpDown] Background tasks completed")
				}()
			} else {
				log.Println(pathLOG + "Scaling [ScaleUpDown] WARNING type of task is not defined: " + dbTask.Type)
			}
		}
		return status, nil
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Error(pathLOG+"Scaling [ScaleUpDown] ERROR (2) ", err)
	return "", err
}
