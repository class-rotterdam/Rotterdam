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

package structs

import (
	imec_db "atos/rotterdam/database/imec"
	constants "atos/rotterdam/globals/constants"
	"log"
	"time"
)

// getClusterDefaultDock Returns the type of a cluster identified by its name / id value
func getClusterDefaultDock(idCluster string) string {
	resp, err := imec_db.GetInfrByID(idCluster)
	if err == nil && len(resp) > 0 {
		return resp[0].DefaultDock
	}

	log.Println("Rotterdam > GLOBALS > structs [getClusterDefaultDock] WARNING Returning default cluster dock ...")
	return constants.DefaultDock
}

/*
TransfCOMPSSTASKtoTASK transforms a CLASS_COMPSS_TASK object into CLASS_TASK object
*/
func TransfCOMPSSTASKtoTASK(classCOMPSsTask *CLASS_COMPSS_TASK) *CLASS_TASK {
	log.Println("Rotterdam > GLOBALS > structs [TransfCOMPSSTASKtoTASK] CLASS_COMPSS_TASK to CLASS_TASK object ...")

	var classTask *CLASS_TASK
	classTask = new(CLASS_TASK)

	// main fields:
	classTask.ID = classCOMPSsTask.ID
	classTask.Name = classCOMPSsTask.Name
	// cluster
	if len(classCOMPSsTask.Cluster) == 0 {
		classTask.Cluster = constants.MainClusterID
	} else {
		classTask.Cluster = classCOMPSsTask.Cluster
	}
	// namespace / dock
	if len(classCOMPSsTask.Dock) == 0 {
		classTask.Dock = getClusterDefaultDock(classCOMPSsTask.Cluster)
	} else {
		classTask.Dock = classCOMPSsTask.Dock
	}
	classTask.Replicas = classCOMPSsTask.Replicas
	classTask.Created = time.Now().Format(time.RFC3339) // current local time

	// QoS:
	classTask.QoSCOMPSs = classCOMPSsTask.Qos
	// classTask.Qos // empty

	// Containers:
	classTask.Containers = make([]CLASS_TASK_CONTAINER, 1)
	classTask.Containers[0].Name = classCOMPSsTask.Name
	classTask.Containers[0].Image = classCOMPSsTask.Image
	classTask.Containers[0].Command = classCOMPSsTask.Command
	classTask.Containers[0].Args = classCOMPSsTask.Args
	// Containers - Ports:
	totalPorts := len(classCOMPSsTask.Ports)
	classTask.Containers[0].Ports = make([]CLASS_TASK_CONTAINER_PORTS, totalPorts)
	for i := 0; i < totalPorts; i++ {
		classTask.Containers[0].Ports[i].ContainerPort = classCOMPSsTask.Ports[i]
		classTask.Containers[0].Ports[i].HostPort = classCOMPSsTask.Ports[i]
		classTask.Containers[0].Ports[i].Protocol = "tcp"
	}

	log.Println("Rotterdam > GLOBALS > structs [TransfCOMPSSTASKtoTASK] Returning CLASS_TASK object ...")

	return classTask
}

/*
TransfTASKtoFunctionTASK transforms a CLASS_TASK object into CLASS_FUNCTION_TASK object
*/
func TransfTASKtoFunctionTASK(classTask *CLASS_TASK) *CLASS_FUNCTION_TASK {
	log.Println("Rotterdam > GLOBALS > structs [TransfCOMPSSTASKtoTASK] CLASS_TASK to CLASS_FUNCTION_TASK object ...")

	var classFunctionTask *CLASS_FUNCTION_TASK
	classFunctionTask = new(CLASS_FUNCTION_TASK)

	// main fields:
	classFunctionTask.ID = classTask.ID
	classFunctionTask.Name = classTask.Name
	// cluster
	if len(classTask.Cluster) == 0 {
		classFunctionTask.Cluster = constants.MainClusterID
	} else {
		classFunctionTask.Cluster = classTask.Cluster
	}
	// namespace / dock
	if len(classTask.Dock) == 0 {
		classFunctionTask.Dock = getClusterDefaultDock(classTask.Cluster)
	} else {
		classFunctionTask.Dock = classTask.Dock
	}

	classFunctionTask.Created = time.Now().Format(time.RFC3339) // current local time

	// QoS:
	classFunctionTask.Qos = classTask.Qos

	// Function:
	classFunctionTask.Runtime = classTask.Runtime
	classFunctionTask.Timeout = classTask.Timeout
	classFunctionTask.Handler = classTask.Handler
	classFunctionTask.FunctionType = classTask.FunctionType
	classFunctionTask.Function = classTask.Function

	log.Println("Rotterdam > GLOBALS > structs [TransfCOMPSSTASKtoTASK] Returning CLASS_FUNCTION_TASK object ...")

	return classFunctionTask
}
