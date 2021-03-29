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
	log "atos/rotterdam/common/logs"
	db "atos/rotterdam/database/caas"
	structs "atos/rotterdam/globals/structs"
)

/*
GetAllTasks Get all tasks
*/
func GetAllTasks() ([]structs.DB_TASK, error) {
	log.Println(pathLOG + "[GetAllTasks] Getting tasks ...")

	dbtasks, err := db.ReadAllTasks()
	if err == nil {
		return dbtasks, nil
	}

	log.Error(pathLOG+"[GetAllTasks] ERROR ", err)
	dbtasks = make([]structs.DB_TASK, 0)
	return dbtasks, err
}

/*
GetDockTasks get all tasks from dock
*/
func GetDockTasks(dock string) ([]structs.DB_TASK, error) {
	log.Println(pathLOG + "[GetDockTasks] Getting tasks from dock [" + dock + "] ...")

	dbtasks, err := db.ReadAllDockTasks(dock)
	if err == nil {
		return dbtasks, nil
	}

	log.Error(pathLOG+"[GetDockTasks] ERROR ", err)
	dbtasks = make([]structs.DB_TASK, 0)
	return dbtasks, err
}

/*
GetAllTasksQoS Get all tasks qos
*/
func GetAllTasksQoS() ([]structs.DB_TASK_QOS, error) {
	log.Println(pathLOG + "[GetAllTasksQoS] Getting tasks QoS ...")

	dbtasks, err := db.DBReadAllTasksQos()
	if err == nil {
		return dbtasks, nil
	}

	log.Error(pathLOG+"[GetAllTasksQoS] ERROR ", err)
	dbtasks = make([]structs.DB_TASK_QOS, 0)
	return dbtasks, err
}

/*
GetTask Returns a task
*/
func GetTask(idTask string) (structs.DB_TASK, error) {
	log.Println(pathLOG + "[GetTask] Getting Task with id=" + idTask + " ...")

	// get task
	dbTask, err := db.ReadTaskValue(idTask)
	if err == nil {
		return *dbTask, nil
	}

	log.Error(pathLOG+"[GetTask] ERROR", err)
	return *dbTask, err
}
