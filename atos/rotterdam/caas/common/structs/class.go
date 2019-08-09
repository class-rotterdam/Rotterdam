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

package structs

import (
	cfg "atos/rotterdam/config"
)

///////////////////////////////////////////////////////////////////////////////
// COMMON STRUCTS

// Configuration
type Configuration struct {
	Adapter        string
	KubeConfigPath string
}

// Http Response content
type RespHttp struct {
	RespStr string
}

// CaaS Response
type ResponseCaaS struct {
	Resp               string `json:"resp,omitempty"`
	Method             string `json:"method,omitempty"`
	Message            string `json:"message,omitempty"`
	CaaSVersion        string `json:"caasversion,omitempty"`
	RulesEngineVersion string `json:"rulesengineversion,omitempty"`
	RestApiVersion     string `json:"restapi,omitempty"`
	SLALiteVersion     string `json:"slalite,omitempty"`
	Content            string `json:"content,omitempty"`
}

// CaaS Response
type ResponseCaaSTask struct {
	Resp        string  `json:"resp,omitempty"`
	Method      string  `json:"method,omitempty"`
	Message     string  `json:"message,omitempty"`
	CaaSVersion string  `json:"caasversion,omitempty"`
	Content     string  `json:"content,omitempty"`
	URL         string  `json:"url,omitempty"`
	Task        DB_TASK `json:"task,omitempty"`
}

// ResponseCaaSTasks
type ResponseCaaSTasks struct {
	Resp        string    `json:"resp,omitempty"`
	Method      string    `json:"method,omitempty"`
	Message     string    `json:"message,omitempty"`
	CaaSVersion string    `json:"caasversion,omitempty"`
	Content     string    `json:"content,omitempty"`
	Tasks       []DB_TASK `json:"tasks,omitempty"`
}

// ResponseCaaSTasksQoS
type ResponseCaaSTasksQoS struct {
	Resp        string        `json:"resp,omitempty"`
	Method      string        `json:"method,omitempty"`
	Message     string        `json:"message,omitempty"`
	CaaSVersion string        `json:"caasversion,omitempty"`
	Content     string        `json:"content,omitempty"`
	TasksQoS    []DB_TASK_QOS `json:"tasksqos,omitempty"`
}

// ResponseQoSDefinition
type ResponseQoSDefinition struct {
	Resp        string                       `json:"resp,omitempty"`
	Method      string                       `json:"method,omitempty"`
	Message     string                       `json:"message,omitempty"`
	CaaSVersion string                       `json:"caasversion,omitempty"`
	Content     string                       `json:"content,omitempty"`
	QoSDef      cfg.CLASS_QOS_TEMPLATE__ELEM `json:"qosdef,omitempty"`
}

// ResponseQoSDefinitions
type ResponseQoSDefinitions struct {
	Resp        string                         `json:"resp,omitempty"`
	Method      string                         `json:"method,omitempty"`
	Message     string                         `json:"message,omitempty"`
	CaaSVersion string                         `json:"caasversion,omitempty"`
	Content     string                         `json:"content,omitempty"`
	QoSDefs     []cfg.CLASS_QOS_TEMPLATE__ELEM `json:"qosdefs,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// CLASS TASK (JSON) DEFINITION
//
// JSON FORMAT - EXAMPLE:
// {
// 	"name": "adas-ped-detection",
// 	"dock": "adas-pro",
// 	"qos": {},
// 	"replicas": 1,
// 	"containers": [{
// 			"name": "adas-ped-distance",
// 			"image": "docker.class.com/adas/adas_task_1:0.1.1",
// 			"ports": [
// 				{
// 					"containerPort": "80",
// 					"hostPort": "80",
// 					"protocol": "tcp",
// 				}
// 			],
// 			"volumes": [
// 				{
// 					"name": "workdir",
// 					"mounthPath": "/usr/share/class/adas/"
// 				}
// 			],
// 			  "environment": [
// 				   {
// 				   "name": "COMPS_VERSION",
// 				   "value":"1.2.1"
// 				   },
// 				   {
// 				   "name": "COMPS_MASTER_URL",
// 				   "value":"<ES-URL>"
// 				   }
// 			  ]
// 		 }
// 	   ]
// }

type CLASS_TASK_CONTAINER_PORTS struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type CLASS_TASK_CONTAINER_VOLUMES struct {
	Name      string `json:"name,omitempty"`
	MountPath string `json:"mounthPath,omitempty"`
}

type CLASS_TASK_CONTAINER_ENV struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type CLASS_TASK_CONTAINER struct {
	Name        string                         `json:"name,omitempty"`
	Image       string                         `json:"image,omitempty"`
	Ports       []CLASS_TASK_CONTAINER_PORTS   `json:"ports,omitempty"`
	Volumes     []CLASS_TASK_CONTAINER_VOLUMES `json:"volumes,omitempty"`
	Environment []CLASS_TASK_CONTAINER_ENV     `json:"environment,omitempty"`
}

type CLASS_TASK_QOS_CUSTOM_GUARANTEES struct {
	Metric    string `json:"metric,omitempty"`
	Condition string `json:"condition,omitempty"`
	Value     string `json:"value,omitempty"`
}

type CLASS_TASK_QOS_CUSTOM struct {
	Type        string                             `json:"type,omitempty"` // type: "infr" or "task"
	Name        string                             `json:"name,omitempty"`
	Description string                             `json:"description,omitempty"`
	Guarantees  []CLASS_TASK_QOS_CUSTOM_GUARANTEES `json:"guarantees,omitempty"`
	ScaleFactor int                                `json:"scalefactor,omitempty"`
	Action      string                             `json:"action,omitempty"`
	Max         int                                `json:"max,omitempty"`
	Min         int                                `json:"min,omitempty"`
}

type CLASS_TASK_QOS struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Custom      []CLASS_TASK_QOS_CUSTOM `json:"custom,omitempty"`
}

type CLASS_TASK struct {
	Name       string                 `json:"name,omitempty"`
	Id         string                 `json:"id,omitempty"`
	Dock       string                 `json:"dock,omitempty"`
	Qos        CLASS_TASK_QOS         `json:"qos,omitempty"`
	Replicas   int                    `json:"replicas,omitempty"`
	Containers []CLASS_TASK_CONTAINER `json:"containers,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// ViolationInfo
type ViolationInfo struct {
	Agreement_id string
	Status       string
	Client_id    string
	Client_name  string
	Guarantee    string
}

///////////////////////////////////////////////////////////////////////////////
// SLA_AGREEMENT

type SLA_AGREEMENT_DETAILS_GUARANTIES struct {
	Name       string `json:"name,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

type SLA_AGREEMENT_DETAILS_ST struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SLA_AGREEMENT_DETAILS struct {
	Id         string                             `json:"id,omitempty"`
	Type       string                             `json:"type,omitempty"`
	Name       string                             `json:"name,omitempty"`
	Provider   SLA_AGREEMENT_DETAILS_ST           `json:"provider,omitempty"`
	Client     SLA_AGREEMENT_DETAILS_ST           `json:"client,omitempty"`
	Creation   string                             `json:"creation,omitempty"`
	Expiration string                             `json:"expiration,omitempty"`
	Guarantees []SLA_AGREEMENT_DETAILS_GUARANTIES `json:"guarantees,omitempty"`
}

type SLA_AGREEMENT struct {
	Id      string                `json:"id,omitempty"`
	Name    string                `json:"name,omitempty"`
	State   string                `json:"state,omitempty"`
	Details SLA_AGREEMENT_DETAILS `json:"details,omitempty"`
}
