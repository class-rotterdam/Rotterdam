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

// getK8sScale: k8s: get scale info
func getK8sScale(task structs.CLASS_TASK) (*structs.K8S_SCALE, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [getK8sScale] Getting scaling info from task [" + task.Name + "] ...")
	_, data, err := common.HttpGET_String(
		cfg.Config.Clusters[0].KubernetesEndPoint + "/apis/apps/v1/namespaces/" + task.Dock + "/deployments/" + task.Name + "/scale")
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [getK8sScale] ERROR", err)
		return nil, err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [getK8sScale] RESPONSE: " + data)

	return common.CommStringToK8S_SCALE(data)
}

// updateK8sScale: k8s: update scale info => scale up / down
func updateK8sScale(task structs.CLASS_TASK, scale_obj structs.K8S_SCALE) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [updateK8sScale] Updating scaling info from task [" + task.Name + "] ...")
	status, _, err := common.HttpPUT_GenericStruct(
		cfg.Config.Clusters[0].KubernetesEndPoint+"/apis/apps/v1/namespaces/"+task.Dock+"/deployments/"+task.Name+"/scale",
		scale_obj)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [updateK8sScale] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [updateK8sScale] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// ScaleUpDown: Scale up task
func ScaleUpDown(dbtask structs.DB_TASK, replicas int) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [ScaleUpDown] Scaling up/down task [" + dbtask.TaskDefinition.Name + "] ...")

	scale_obj, err := getK8sScale(dbtask.TaskDefinition)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [ScaleUpDown] ERROR (1) ", err)
		return "", err
	} else {
		scale_obj.Spec.Replicas = replicas
		status, err := updateK8sScale(dbtask.TaskDefinition, *scale_obj)
		if err == nil {
			return status, nil
		}
		err = errors.New("Task creation failed. status = [" + status + "]")
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Scaling [ScaleUpDown] ERROR (2) ", err)
		return "", err
	}
}
