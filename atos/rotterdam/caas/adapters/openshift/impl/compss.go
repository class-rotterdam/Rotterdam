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
	"strconv"
	"strings"
	"sync"
	"time"
)

// port number used to map compss applications
var (
	mu    sync.Mutex // guards balance
	rport int
)

// init
func init() {
	mu.Lock()
	rport = 25000
	mu.Unlock()
}

// set rport value
func setRPort() {
	mu.Lock()
	rport = rport + 1
	mu.Unlock()
}

// read rport value
func readRPort() int {
	mu.Lock()
	b := rport
	mu.Unlock()
	return b
}

// generate and read rport value
func newRPort() int {
	mu.Lock()
	rport = rport + 1
	b := rport
	mu.Unlock()
	return b
}

// checkPodsFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func checkPodsFromTask(cluster_index int, namespace string, name string, expected_replicas int) (string, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Getting info from pods ...")

	// get pods from task
	_, result, err := common.HttpGET_GenericStruct(cfg.Config.Clusters[cluster_index].KubernetesEndPoint + "/api/v1/namespaces/" +
		namespace + "/pods?labelSelector=app=" + name)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] ERROR", err)
		return "error", nil, err
	}

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Retrieved pods = " + strconv.Itoa(len(items)))
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Expected pods = " + strconv.Itoa(expected_replicas))

	if len(items) == expected_replicas {
		return "ready", result, err
	}

	return "not-ready", result, err
}

// get main ports of the application
func getMainPort(task structs.CLASS_TASK) (int, string) {
	// ports
	main_port := 0
	main_protocol := ""
	for _, contElement := range task.Containers {
		for _, portElement := range contElement.Ports {
			if main_port == 0 {
				main_port = portElement.ContainerPort
				main_protocol = strings.ToUpper(portElement.Protocol)
			}
		}
	}

	if main_port == 0 {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getMainPort] ERROR getting main port. Returning 80 ...")
		main_port = 80
		main_protocol = "TCP"
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getMainPort] main port = " + strconv.Itoa(main_port))
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getMainPort] main protocol = " + main_protocol)
	return main_port, main_protocol
}

// getPodsInfoFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func initPodsInfoFromTask(cluster_index int, namespace string, name string, result map[string]interface{},
	app_port int, app_protocol string) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

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
		podData.IP = cfg.Config.Clusters[cluster_index].ServerIP
		podData.Port = newRPort()
		podData.TargetPort = app_port
		podData.Protocol = app_protocol

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

		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] POD " + podData.Name)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - Status: " + podData.Status)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

// patchPods: k8s: patch pods names
func patchPods(cluster_index int, namespace string, pod_name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [patchPods] Generating 'pod patch' json ...")
	l_k8s_patch_pod := common.StructNewPodPatch(pod_name) // returns []structs.K8S_POD_PATCH_LINE

	str_txt, _ := common.CommPatchPodsListToString(l_k8s_patch_pod)
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [patchPods] [" + str_txt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [patchPods] Patching pod " + pod_name + " ...")
	status, _, err := common.HttpPATCH_GenericStruct(
		cfg.Config.Clusters[cluster_index].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/pods/"+pod_name,
		l_k8s_patch_pod)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [patchPods] ERROR", err)
		return "error", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [patchPods] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// podService: k8s: service
func podService(cluster_index int, namespace string, pod structs.DB_TASK_POD) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [podService] Generating 'service' json ...")
	k8s_serv := common.StructNewPodServiceTemplate(cluster_index, pod) // returns *K8S_SERVICE

	str_txt, _ := common.CommServiceStructToString(*k8s_serv)
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [podService] [" + str_txt + "]")

	// CALL to Kubernetes API to launch a new service
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [podService] Creating a new service in K8s cluster ...")
	status, _, err := common.HttpPOST_GenericStruct(
		cfg.Config.Clusters[cluster_index].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/services",
		k8s_serv)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [podService] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [podService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// Expose services and save task
func compssDeploymentBackgroundTasks(cluster_index int, task structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Executing background tasks ...")

	time.Sleep(10 * time.Second)
	main_port, main_protocol := getMainPort(task)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Checking status of task pods ...")
		str, result, err := checkPodsFromTask(cluster_index, task.Dock, task.Name, task.Replicas)

		if err == nil && str == "ready" {
			// get info from pods
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Getting all pods info (ids, names, IPs, ports) ...")
			lTaskPods := initPodsInfoFromTask(cluster_index, task.Dock, task.Name, result, main_port, main_protocol)

			dbTask, err := common.ReadTaskValue(task.Name)
			if err == nil {
				// update task with pods info
				dbTask.Pods = lTaskPods
				log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks]Saving info to db ...")
				err = common.SetTaskValue(task.Name, *dbTask)

				if err == nil {
					// update pods: pod-name (labels)
					log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Updating pods (/metadata/labels/pod-name) / Creating and exposing pods services ...")
					for _, pod := range lTaskPods {
						log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Updating pod " + pod.Name)
						status, err := patchPods(cluster_index, task.Dock, pod.Name)
						// TODO check errors
						log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] status: " + status)

						if err == nil {
							log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Creating and exposing pod's service ...")

							_, _ = podService(cluster_index, task.Dock, pod)
						}
					}

					log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] Finishing background process ...")

					break
				}
			}
		} else if err == nil {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] str = " + str)
		} else {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [deploymentBackgroundTasks] ERROR (1) ", err)
		}

		time.Sleep(30 * time.Second)
	}
}

func checkCurrentPodsInfoFromTask(cluster_index int, namespace string, name string, result map[string]interface{},
	app_port int, app_protocol string) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

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
		podData.IP = cfg.Config.Clusters[cluster_index].ServerIP
		podData.Port = newRPort()
		podData.TargetPort = app_port
		podData.Protocol = app_protocol

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

		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] POD " + podData.Name)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - Status: " + podData.Status)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [getPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

// Expose new services and save task
func compssScalingOutBackgroundTasks(cluster_index int, dbTask structs.DB_TASK, items map[string]interface{}) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Exposing new services ...")

	// get main port and protocol
	main_port, main_protocol := getMainPort(dbTask.TaskDefinition)

	// get list of pods
	lTaskPods := checkCurrentPodsInfoFromTask(cluster_index, dbTask.TaskDefinition.Dock, dbTask.TaskDefinition.Name, items, main_port, main_protocol)

	// add new services
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Updating new pods (/metadata/labels/pod-name) / Creating and exposing new pods services ...")
	for _, pod := range lTaskPods {
		notfound := true
		for _, podold := range dbTask.Pods {
			if pod.Name == podold.Name {
				notfound = false
				break
			}
		}

		if notfound == true {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Updating pod " + pod.Name)
			status, err := patchPods(cluster_index, dbTask.TaskDefinition.Dock, pod.Name)
			// TODO check errors
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] status: " + status)

			if err == nil {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Creating and exposing pod's service ...")
				_, _ = podService(cluster_index, dbTask.TaskDefinition.Dock, pod)
			}
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Update finished")

	return lTaskPods
}

// Remove unused services and save task
func compssScalingInBackgroundTasks(cluster_index int, dbTask structs.DB_TASK, items map[string]interface{}) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingInBackgroundTasks] Removing unused services ...")

	// get main port and protocol
	main_port, main_protocol := getMainPort(dbTask.TaskDefinition)

	// get list of pods
	lTaskPods := checkCurrentPodsInfoFromTask(cluster_index, dbTask.TaskDefinition.Dock, dbTask.TaskDefinition.Name, items, main_port, main_protocol)

	// remove 'old' services
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Removing unused services ...")
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
			log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Removing unused service of old pod [" + podold.Name + "] ...")
			_, _ = delK8sService(cluster_index, dbTask.TaskDefinition.Dock, podold.Name)
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [compssScalingOutBackgroundTasks] Services update finished")

	return lTaskPods
}

// Update (remove unused / add new) services used to expose pods
func CompssScalingUpdateServices(cluster_index int, dbTask structs.DB_TASK, new_replicas_value int) structs.DB_TASK {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [CompssScalingUpdateServices] Executing background tasks ...")

	time.Sleep(20 * time.Second)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [CompssScalingUpdateServices] Checking status of task pods [new_replicas_value=" + strconv.Itoa(new_replicas_value) + "] ...")
		str, items, err := checkPodsFromTask(cluster_index, dbTask.TaskDefinition.Dock, dbTask.TaskDefinition.Name, new_replicas_value)
		if err == nil && str == "ready" {
			if dbTask.Replicas < new_replicas_value {
				dbTask.Pods = compssScalingOutBackgroundTasks(cluster_index, dbTask, items)
				break
			} else if dbTask.Replicas > new_replicas_value {
				dbTask.Pods = compssScalingInBackgroundTasks(cluster_index, dbTask, items)
				break
			} else {
				// don't update the task
				log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [CompssScalingUpdateServices] WARNING type of task is not defined: " + dbTask.Type)
				return dbTask
			}
		}

		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [CompssScalingUpdateServices] Result: " + str)
		log.Println("Rotterdam > CAAS > Adapters > Openshift > COMPSs [CompssScalingUpdateServices] Trying again in 20s ...")
		time.Sleep(20 * time.Second)
	}

	dbTask.TaskDefinition.Replicas = new_replicas_value
	dbTask.Replicas = new_replicas_value
	return dbTask
}

// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
/*
{
  "items": [{
      "metadata": {
        "name": "nginx-app-68b6bd8f76-c6hvq",
        "generateName": "nginx-app-68b6bd8f76-",
        "namespace": "class",
        "selfLink": "/api/v1/namespaces/class/pods/nginx-app-68b6bd8f76-c6hvq",
        "uid": "c29cda84-8c45-11e9-8947-005056986059",
        "resourceVersion": "30855861",
        "creationTimestamp": "2019-06-11T12:38:05Z",
        "labels": {
          "app": "nginx-app",
        },
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-jlw6j",
            "secret": {
              "secretName": "default-token-jlw6j",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "nginx",
            "image": "nginx",
            "ports": [
              {
                "containerPort": 80,
                "protocol": "TCP"
              }
            ],
            "env": [
              {
                "name": "TEST_VALUE",
                "value": "1.2.3"
              }
            ],
            "resources": {},
            "volumeMounts": [
              {
                "name": "default-token-jlw6j",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ]
            }
          }
        ],
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "kube5"
      },
      "status": {
		"phase": "Running",
		...
        "hostIP": "192.168.7.29",
        "podIP": "10.129.0.34",
        "startTime": "2019-06-11T12:38:05Z",
        "containerStatuses": [
          {
            "name": "nginx",
            "state": {
              "running": {
                "startedAt": "2019-06-11T12:38:14Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "docker.io/nginx:latest",
            "imageID": "docker-pullable://docker.io/nginx@sha256:bdbf36b7f1f77ffe7bd2a32e59235dff6ecf131e3b6b5b96061c652f30685f3a",
            "containerID": "docker://363d946ee75efc19a053c1320a5cf20395fe4e0412ce8173b44c69d2a166d936"
          }
        ],
        "qosClass": "BestEffort"
      }
    },
*/
