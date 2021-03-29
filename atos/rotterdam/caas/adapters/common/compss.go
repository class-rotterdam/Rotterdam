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

package common

import (
	urls "atos/rotterdam/caas/adapters"
	common "atos/rotterdam/caas/common"
	log "atos/rotterdam/common/logs"
	db "atos/rotterdam/database/caas"
	imec_db "atos/rotterdam/database/imec"
	structs "atos/rotterdam/globals/structs"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
)

/*
Get Pods from current task. Retrieved Pods = Expected Pods
*/
// checkPodsFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func checkPodsFromTask(namespace string, id string, expectedReplicas int, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, map[string]interface{}, error) {
	log.Println(pathLOG + "COMPSs [checkPodsFromTask] Getting pods ...")

	// get pods from task
	_, result, err := common.HTTPGETStruct(
		urls.GetPathKubernetesPodsApp(cluster, namespace, id),
		true)
	if err != nil {
		log.Error(pathLOG+"COMPSs [checkPodsFromTask] WARNING (1)", err)
		_, result, err = common.HTTPGETStruct(
			urls.GetPathKubernetesPodsApp(cluster, namespace, id),
			false)
		if err != nil {
			log.Error(pathLOG+"COMPSs [checkPodsFromTask] ERROR (2)", err)
			return "error", nil, err
		}
	}

	// items
	items := result["items"].([]interface{})
	log.Debug(pathLOG + "COMPSs [checkPodsFromTask] Retrieved pods = " + strconv.Itoa(len(items)))
	log.Debug(pathLOG + "COMPSs [checkPodsFromTask] Expected pods = " + strconv.Itoa(expectedReplicas))

	if len(items) == expectedReplicas {
		return "ready", result, err
	}

	return "not-ready", result, err
}

/*
Get Pods from current task. Retrieved Pods = Expected Pods
*/
// curl -X GET -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9"
// http://192.168.7.28:8001/api/v1/namespaces/class/services/serv-g7bv9fu6d0rtaius7okjnk1580403141125383512-5c465cbb4b-xqdl8
func getServiceFromPod(namespace string, idPod string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, SERV_POD, error) {
	log.Println(pathLOG + "COMPSs [getServiceFromPod] Getting service from Pod [" + idPod + "] ...")

	// get pods from task
	_, result, err := common.HTTPGETString(
		urls.GetPathKubernetesService(cluster, namespace, idPod),
		true)
	if err != nil {
		log.Error(pathLOG+"COMPSs [getServiceFromPod] WARNING (1)", err)
		_, result, err = common.HTTPGETString(
			urls.GetPathKubernetesService(cluster, namespace, idPod),
			false)
		if err != nil {
			log.Error(pathLOG+"COMPSs [getServiceFromPod] ERROR (2)", err)
			return "error", SERV_POD{}, err
		}
	}

	// metadata
	sp, err := StringToServPodStruct(result)
	if err != nil {
		log.Error(pathLOG+"COMPSs [getServiceFromPod] ERROR (3)", err)
		return "error", SERV_POD{}, err
	}

	if sp.Metadata.Name == "" {
		return "not-ready", SERV_POD{}, err
	}

	log.Println("Rotterdam > CAAS > Adapters > Common > COMPSs [getServiceFromPod] Retrieved metadata.name = " + sp.Metadata.Name)

	return "ready", *sp, nil
}

// getPodsInfoFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func getPodsInfoFromTask(result map[string]interface{}, appPort int, appProtocol string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) []structs.DB_TASK_POD {
	log.Println(pathLOG + "COMPSs [getPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println(pathLOG + "COMPSs [getPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

	// final res
	var lres []structs.DB_TASK_POD

	// iterate json (response)
	for _, item := range items { // "items" element ==> PODS INFO
		/* DB_TASK_POD:
		{
			Name       string `json:"name,omitempty"`
			IP         string `json:"ip,omitempty"`         // IP accessed by external apps
			HostIP     string `json:"hostIp,omitempty"`     // node IP
			PodIP      string `json:"podIp,omitempty"`      // internal IP created by Kubernetes / Openshift
			Status     string `json:"status,omitempty"`     // running, unknown
			Port       int    `json:"port,omitempty"`       // port exposed in Kubernetes / Openshift
			TargetPort int    `json:"targetPort,omitempty"` // application port
		}
		*/
		podData := &structs.DB_TASK_POD{}
		podData.IP = cluster.HostIP
		podData.Port = NewRPort()
		podData.TargetPort = appPort
		podData.Protocol = appProtocol

		// "metadata" element
		for key1, value1 := range item.(map[string]interface{}) {
			if key1 == "metadata" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "name" {
						podData.Name = value2.(string) // Pod name
					}
				}
			} else if key1 == "status" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "podIP" {
						podData.PodIP = value2.(string) // IP address from pod
					} else if key2 == "hostIP" {
						podData.HostIP = value2.(string) // IP address from node
					} else if key2 == "phase" {
						podData.Status = value2.(string) // Status / Phase from pod
					}
				}
			}
		}

		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask] POD " + podData.Name)
		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask]    - Status: " + podData.Status)
		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Debug(pathLOG + "COMPSs [getPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println(pathLOG + "COMPSs [getPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

// patchPods: k8s: patch pods names
func patchPods(namespace string, podName string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println(pathLOG + "COMPSs [patchPods] Generating 'pod patch' json ...")
	lK8sPatchPod := structs.StructNewPodPatch(podName) // returns []structs.K8S_POD_PATCH_LINE

	strTxt, _ := structs.CommPatchPodsListToString(lK8sPatchPod)
	log.Println(pathLOG + "COMPSs [patchPods] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println(pathLOG + "COMPSs [patchPods] Patching pod " + podName + " ...")
	status, _, err := common.HTTPPATCHStruct(
		urls.GetPathKubernetesPod(cluster, namespace, podName),
		true,
		lK8sPatchPod)
	if err != nil {
		log.Error(pathLOG+"COMPSs [patchPods] WARNING (1)", err)
		status, _, err = common.HTTPPATCHStruct(
			urls.GetPathKubernetesPod(cluster, namespace, podName),
			false,
			lK8sPatchPod)
		if err != nil {
			log.Error(pathLOG+"COMPSs [patchPods] ERROR (2)", err)
			return "error", err
		}
	}
	log.Println(pathLOG + "COMPSs [patchPods] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// podService: k8s: service
func podService(namespace string, pod structs.DB_TASK_POD, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println(pathLOG + "COMPSs [podService] Generating 'service' json ...")
	k8sServ := structs.StructNewPodServiceTemplate(cluster.HostIP, pod) // returns *K8S_SERVICE

	strTxt, _ := structs.CommServiceStructToString(*k8sServ)
	log.Println(pathLOG + "COMPSs [podService] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new service
	log.Println(pathLOG + "COMPSs [podService] Creating a new service in K8s cluster ...")
	status, _, err := common.HTTPPOST(
		urls.GetPathKubernetesCreateService(cluster, namespace),
		true,
		k8sServ)
	if err != nil {
		log.Error(pathLOG+"COMPSs [podService] WARNING (1)", err)
		status, _, err = common.HTTPPOST(
			urls.GetPathKubernetesCreateService(cluster, namespace),
			false,
			k8sServ)
		if err != nil {
			log.Error(pathLOG+"COMPSs [podService] ERROR (2)", err)
			return "", err
		}
	}
	log.Println(pathLOG + "COMPSs [podService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// podServiceTrash: k8s: service
func podServiceTrash(podname string, namespace string, pod structs.DB_TASK_POD, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println(pathLOG + "COMPSs [podServiceTrash] Generating trash 'service' json ...")

	id := shortuuid.NewWithAlphabet("0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx")
	id = strings.ToLower(id)
	podname = podname + id
	log.Println(pathLOG + "COMPSs [podServiceTrash] [" + podname + "]")

	pod.Port = 20013 + rand.Intn(100)
	k8sServ := structs.StructNewPodWithNameServiceTemplate(cluster.HostIP, pod, podname) // returns *K8S_SERVICE

	// CALL to Kubernetes API to launch a new service
	log.Println(pathLOG + "COMPSs [podServiceTrash] Creating a new service in K8s cluster ...")
	_, _, err := common.HTTPPOST(
		urls.GetPathKubernetesCreateService(cluster, namespace),
		true,
		k8sServ)
	if err != nil {
		log.Error(pathLOG+"COMPSs [podServiceTrash] WARNING (1)", err)
		_, _, err = common.HTTPPOST(
			urls.GetPathKubernetesCreateService(cluster, namespace),
			false,
			k8sServ)
		if err != nil {
			log.Error(pathLOG+"COMPSs [podServiceTrash] ERROR (2)", err)
			return "", err
		}
	}
	log.Println(pathLOG + "COMPSs [podServiceTrash] 'trash' service created")

	time.Sleep(10 * time.Second)

	log.Println(pathLOG + "COMPSs [podServiceTrash] Deleting trash 'service' ...")
	statusS, err := DelK8sService(namespace, podname, cluster, true)
	if err != nil {
		log.Error(pathLOG+"COMPSs [removeCOMPSsTask] ERROR (2)", err)
		_, _ = DelK8sService(namespace, podname, cluster, false)
	} else if statusS == "200" {
		log.Println(pathLOG + "COMPSs [removeCOMPSsTask] Task removed with success")
		//return "ok", nil
	}

	return statusS, nil
}

/*
CompssDeploymentBackgroundTasks Deployment of a COMPSs task
*/
func CompssDeploymentBackgroundTasks(task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) {
	log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Executing background tasks ...")

	time.Sleep(10 * time.Second)
	mainPort, mainProtocol := GetMainPort(task)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Checking status of task pods ...")
		str, result, err := checkPodsFromTask(task.Dock, task.ID, task.Replicas, cluster)

		if err == nil && str == "ready" { // Expected Pods = Retrieved Pods
			// get info from pods
			log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Getting all pods info (ids, names, IPs, ports) ...")
			lTaskPods := getPodsInfoFromTask(result, mainPort, mainProtocol, cluster)

			dbTask, err := db.ReadTaskValue(task.ID)
			if err == nil {
				// update task with pods info
				dbTask.Pods = lTaskPods
				log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks]Saving info to db ...")
				err = db.SetTaskValue(task.ID, *dbTask)

				if err == nil {
					// update pods: pod-name (labels)
					log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Updating pods (/metadata/labels/pod-name) / Creating and exposing pods services ...")
					for i, pod := range lTaskPods {
						log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Updating pod " + pod.Name)
						status, err := patchPods(task.Dock, pod.Name, cluster)
						// TODO check errors
						log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] status: " + status)

						if err == nil {
							log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Creating and exposing pod's service ...")

							_, _ = podService(task.Dock, pod, cluster)

							// solve Openshift "problem" where the last service is created but cannot be accessed
							if i == len(lTaskPods)-1 {
								_, _ = podServiceTrash("trash", task.Dock, pod, cluster)
							}
						}
					}

					log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] Finishing background process ...")

					break
				}
			}
		} else if err == nil {
			log.Println(pathLOG + "COMPSs [compssDeploymentBackgroundTasks] str = " + str)
		} else {
			log.Error(pathLOG+"COMPSs [compssDeploymentBackgroundTasks] ERROR (1) ", err)
		}

		time.Sleep(30 * time.Second)
	}
}

// checkCurrentPodsInfoFromTask
func checkCurrentPodsInfoFromTask(result map[string]interface{}, appPort int, appProtocol string, dbTask structs.DB_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) []structs.DB_TASK_POD {
	log.Println(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

	// final res
	var lres []structs.DB_TASK_POD

	// iterate json (response)
	for _, item := range items { // "items" element
		/* DB_TASK_POD:
		{
			Name       string `json:"name,omitempty"`
			IP         string `json:"ip,omitempty"`         // IP accessed by external apps
			HostIP     string `json:"hostIp,omitempty"`     // node IP
			PodIP      string `json:"podIp,omitempty"`      // internal IP created by Kubernetes / Openshift
			Status     string `json:"status,omitempty"`     // running, unknown
			Port       int    `json:"port,omitempty"`       // port exposed in Kubernetes / Openshift
			TargetPort int    `json:"targetPort,omitempty"` // application port
		}
		*/
		podData := &structs.DB_TASK_POD{}
		podData.IP = cluster.HostIP
		podData.TargetPort = appPort
		podData.Protocol = appProtocol

		// "metadata" element
		for key1, value1 := range item.(map[string]interface{}) {
			if key1 == "metadata" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "name" {
						podData.Name = value2.(string) // Pod name
					}
				}
			} else if key1 == "status" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "podIP" {
						podData.PodIP = value2.(string) // IP address from pod
					} else if key2 == "hostIP" {
						podData.HostIP = value2.(string) // IP address from node
					} else if key2 == "phase" {
						podData.Status = value2.(string) // Status / Phase from pod
					}
				}
			}
		}

		// GET SERVICE, if exists
		ready, srvPod, err := getServiceFromPod(dbTask.TaskDefinition.Dock, podData.Name, cluster)
		if err == nil && ready == "ready" {
			log.Println(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] POD " + podData.Name + " already exists")
			podData.Port = srvPod.Spec.Ports[0].Port
		} else {
			log.Println(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] Generating new port for POD " + podData.Name)
			podData.Port = NewRPort()
		}

		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] POD " + podData.Name)
		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask]    - Status: " + podData.Status)
		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Debug(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println(pathLOG + "COMPSs [checkCurrentPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

/*
CompssScalingOutBackgroundTasks Expose new services and save task
*/
func CompssScalingOutBackgroundTasks(dbTask structs.DB_TASK, items map[string]interface{}, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) []structs.DB_TASK_POD {
	log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] Exposing new services ...")

	// get main port and protocol
	mainPort, mainProtocol := GetMainPort(dbTask.TaskDefinition)

	// get list of pods - all pods
	lTaskPods := checkCurrentPodsInfoFromTask(items, mainPort, mainProtocol, dbTask, cluster)

	// add new services
	log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] Updating new pods (/metadata/labels/pod-name) / Creating and exposing new pods services ...")
	for i, pod := range lTaskPods { // iterate all pods (including new pods)
		notfound := true
		for _, podold := range dbTask.Pods { // check if pod is new or old
			if pod.Name == podold.Name {
				notfound = false
				break
			}
		}

		if notfound == true {
			log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] Updating pod " + pod.Name)
			status, err := patchPods(dbTask.TaskDefinition.Dock, pod.Name, cluster)
			// TODO check errors
			log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] status: " + status)

			if err == nil {
				log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] Creating and exposing pod's service ...")
				_, _ = podService(dbTask.TaskDefinition.Dock, pod, cluster)
			}
		}

		// solve Openshift "problem" where the last service is created but cannot be accessed
		if i == len(lTaskPods)-1 {
			_, _ = podServiceTrash("trash", dbTask.TaskDefinition.Dock, pod, cluster)
		}
	}

	log.Println(pathLOG + "COMPSs [CompssScalingOutBackgroundTasks] Update finished")
	return lTaskPods
}

/*
CompssScalingInBackgroundTasks Removes unused services and save task
*/
func CompssScalingInBackgroundTasks(dbTask structs.DB_TASK, items map[string]interface{}, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) []structs.DB_TASK_POD {
	log.Println(pathLOG + "COMPSs [CompssScalingInBackgroundTasks] Removing unused services ...")

	// get main port and protocol
	mainPort, mainProtocol := GetMainPort(dbTask.TaskDefinition)

	// get list of pods
	lTaskPods := checkCurrentPodsInfoFromTask(items, mainPort, mainProtocol, dbTask, cluster)

	// remove 'old' services
	log.Println(pathLOG + "COMPSs [CompssScalingInBackgroundTasks] Removing unused services ...")
	for _, podold := range dbTask.Pods {
		notfound := true
		for _, pod := range lTaskPods {
			if pod.Name == podold.Name {
				notfound = false
				break
			}
		}

		// remove service if old pod was not found in the list of current pods
		if notfound == true {
			log.Println(pathLOG + "COMPSs [CompssScalingInBackgroundTasks] Removing unused service of old pod [" + podold.Name + "] ...")
			_, _ = DelK8sService(dbTask.TaskDefinition.Dock, podold.Name, cluster, sec)
		}
	}

	log.Println(pathLOG + "COMPSs [CompssScalingInBackgroundTasks] Services update finished")
	return lTaskPods
}

/*
CompssScalingUpdateServices Updates (remove unused / add new) services used to expose pods. Replicas are scaled in / out before calling this method
*/
func CompssScalingUpdateServices(dbTask structs.DB_TASK, newReplicasValue int, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) structs.DB_TASK {
	log.Println(pathLOG + "COMPSs [CompssScalingUpdateServices] Executing background tasks ...")

	time.Sleep(20 * time.Second)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println(pathLOG + "COMPSs [CompssScalingUpdateServices] Checking status of task pods [new_replicas_value=" + strconv.Itoa(newReplicasValue) + "] ...")
		str, items, err := checkPodsFromTask(dbTask.TaskDefinition.Dock, dbTask.TaskDefinition.ID, newReplicasValue, cluster)
		if err == nil && str == "ready" {
			if dbTask.Replicas < newReplicasValue {
				dbTask.Pods = CompssScalingOutBackgroundTasks(dbTask, items, cluster)
				break
			} else if dbTask.Replicas > newReplicasValue {
				dbTask.Pods = CompssScalingInBackgroundTasks(dbTask, items, cluster, sec)
				break
			} else {
				// don't update the task
				log.Println(pathLOG + "COMPSs [CompssScalingUpdateServices] WARNING type of task is not defined: " + dbTask.Type)
				return dbTask
			}
		}

		log.Debug(pathLOG + "COMPSs [CompssScalingUpdateServices] Result: " + str)
		log.Debug(pathLOG + "COMPSs [CompssScalingUpdateServices] Trying again in 20s ...")
		time.Sleep(20 * time.Second)
	}

	dbTask.TaskDefinition.Replicas = newReplicasValue
	dbTask.Replicas = newReplicasValue
	return dbTask
}
