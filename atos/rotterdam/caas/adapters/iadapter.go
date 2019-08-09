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

package adapters

import (
	structs "atos/rotterdam/caas/common/structs"
)

// generic adapter
type Adapter interface {
	DeployTask(cluster int, namespace string, task structs.CLASS_TASK) (string, error)
	DeployTaskCompss(cluster int, namespace string, task structs.CLASS_TASK) (string, error)
	GetTask(cluster int, namespace string, name string) (structs.DB_TASK, error)
	GetConfig() (string, error)
	ScaleUpDown(cluster int, dbtask structs.DB_TASK, replicas int) (string, error)
	RemoveTask(cluster int, namespace string, name string) (string, error)
}
