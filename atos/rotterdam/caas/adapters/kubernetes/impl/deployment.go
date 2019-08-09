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

// k8sDeployment: k8s: deployment
func k8sDeployment(namespace string, task structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sDeployment] Generating 'deployment' json ...")
	k8s_depl := common.StructNewDeploymentTemplate(task, 1) // returns *K8S_DEPLOYMENT

	str_txt, _ := common.CommDeploymentStructToString(*k8s_depl)
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sDeployment] [" + str_txt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sDeployment] Creating a new deployment in K8s cluster ...")
	status, _, err := common.HttpPOST_GenericStruct( //common.MOCKUP_HttpPOST_GenericStruct(
		cfg.Config.Clusters[0].KubernetesEndPoint+"/apis/apps/v1/namespaces/"+namespace+"/deployments",
		k8s_depl)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sDeployment] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// k8sService: k8s: service
func k8sService(namespace string, task structs.CLASS_TASK) (string, int, string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sService] Generating 'service' json ...")
	k8s_serv, main_port, main_port_name := common.StructNewServiceTempalte(task) // returns *K8S_SERVICE

	str_txt, _ := common.CommServiceStructToString(*k8s_serv)
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sService] [" + str_txt + "]")

	// CALL to Kubernetes API to launch a new service
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sService] Creating a new service in K8s cluster ...")
	status, _, err := common.HttpPOST_GenericStruct( //MOCKUP_HttpPOST_GenericStruct( //.HttpPOST_GenericStruct(
		cfg.Config.Clusters[0].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/services",
		k8s_serv)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sService] ERROR", err)
		return "", -1, "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [k8sService] RESPONSE: OK")

	return strconv.Itoa(status), main_port, main_port_name, nil
}

// DeployTask: Deploy a task (k8s: deployment & service & volumes ...)
func DeployTask(namespace string, task structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTask] Deploying new task ...")

	// 1. DEPLOYMENT /////
	status, err := k8sDeployment(namespace, task)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTask] ERROR (1)", err)
		return "", err
	} else if status == "200" || status == "201" {
		// 2. SERVICE /////
		status, _, _, err := k8sService(namespace, task)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTask] ERROR (2)", err)
			return "", err
		} else if status == "200" || status == "201" {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTask] Task deployed with success")
			// save to DB
			dbtask := &structs.DB_TASK{
				Id:             task.Name,
				Name:           task.Name,
				NameSpace:      namespace,
				Type:           "default",
				Url:            "http://" + task.Name + "." + cfg.Config.Clusters[0].ServerIP + ".xip.io",
				Status:         "Deployed",
				Replicas:       1,
				TaskDefinition: task}
			common.SetTaskValue(task.Name, *dbtask)

			return "ok", nil
		}
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTask] ERROR (4)", err)
	return "", err
}

// DeployTaskCompss: Deploy a COMPSs task (k8s: deployment & service & volumes ...)
func DeployTaskCompss(namespace string, task structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTaskCompss] Deploying new task ...")

	// 1. DEPLOYMENT /////
	status, err := k8sDeployment(namespace, task)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTaskCompss] ERROR (1)", err)
		return "", err
	} else if status == "200" || status == "201" {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTaskCompss] Task deployed with success")
		// save to DB
		dbtask := &structs.DB_TASK{
			Id:             task.Name,
			Name:           task.Name,
			NameSpace:      namespace,
			Type:           "compss",
			Url:            "http://" + task.Name + "." + cfg.Config.Clusters[0].ServerIP + ".xip.io",
			Status:         "Deployed",
			Replicas:       1,
			TaskDefinition: task}
		common.SetTaskValue(task.Name, *dbtask)

		return "ok", nil
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Deployment [DeployTaskCompss] ERROR (4)", err)
	return "", err
}
