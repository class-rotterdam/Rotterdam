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
	db "atos/rotterdam/database/caas"
	imec_db "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
	"errors"
	"log"
	"strconv"
)

// getK8sScale gets the scale info
func getK8sScale(task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (*structs.K8S_SCALE, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [getK8sScale] Getting scaling info from task [" + task.ID + "] ...")
	_, data, err := common.HTTPGETString(
		urls.GetPathKubernetesScaleDeployment(cluster, task.Dock, task.ID),
		true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [getK8sScale] ERROR", err)
		return nil, err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [getK8sScale] RESPONSE: " + data)
	return structs.CommStringToK8S_SCALE(data)
}

// updateK8sScale updates the scale info => scale up / down
func updateK8sScale(task structs.CLASS_TASK, scaleObj structs.K8S_SCALE, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [updateK8sScale] Updating scaling info from task [" + task.ID + "] ...")
	status, _, err := common.HTTPPUTStruct(
		urls.GetPathKubernetesScaleDeployment(cluster, task.Dock, task.ID),
		true,
		scaleObj)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [updateK8sScale] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [updateK8sScale] RESPONSE: OK")
	return strconv.Itoa(status), nil
}

/*
ScaleUpDown scales in / out a task in Openshift.
*/
func ScaleUpDown(dbTask structs.DB_TASK, replicas int) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] Scaling up/down task [" + dbTask.TaskDefinition.Name + "] ...")

	clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
	task := dbTask.TaskDefinition
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] cluster=" + clusterInfr.ID)

	// 1. get scale object (from Kubernetes) needed to scale out / in the application
	scaleObj, err := getK8sScale(task, clusterInfr)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] ERROR (1) ", err)
		return "", err
	}

	// set new value for scale object's replicas
	scaleObj.Spec.Replicas = replicas
	// update application
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
					log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] Starting background tasks...")
					dbTask = adapt_common.CompssScalingUpdateServices(dbTask, replicas, clusterInfr, true)
					_ = db.SetTaskValue(dbTask.Id, dbTask)
					log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] Background tasks completed")
				}()
			} else {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] WARNING type of task is not defined: " + dbTask.Type)
			}
		}
		return status, nil
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Scaling [ScaleUpDown] ERROR (2) ", err)
	return "", err
}
