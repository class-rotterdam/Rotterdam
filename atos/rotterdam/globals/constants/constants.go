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

package constants

/*
MainClusterID Main cluster identifier
*/
const MainClusterID = "maincluster"

/*
DefaultDock default dock / namespace
*/
const DefaultDock = "class"

/*
TypeOpenshift "Openshift"
*/
const TypeOpenshift = "Openshift"

/*
TypeKubernetes "Kubernetes"
*/
const TypeKubernetes = "Kubernetes"

/*
TypeMicroK8s "MicroK8s"
*/
const TypeMicroK8s = "MicroK8s"

/*
TypeKubeless "Kubeless"
*/
const TypeKubeless = "Kubeless"

/*
TypeDocker "Docker"
*/
const TypeDocker = "Docker"

/*
DefaultType default type: Kubernetes / Openshift / Docker / MicroK8s
*/
const DefaultType = TypeOpenshift

/*
DefaultTasksMaxReplicas default max replicas per service
*/
const DefaultTasksMaxReplicas = 10

/*
DefaultTasksMinReplicas default min replicas per service
*/
const DefaultTasksMinReplicas = 1

/*
DefaultTasksMaxAllowed default max allowed
*/
const DefaultTasksMaxAllowed = 2

/*
DefaultTasksScaleFactor default scale factor
*/
const DefaultTasksScaleFactor = 2.5

/*
DefaultTasksValue default value
*/
const DefaultTasksValue = 2

/*
DefaultTasksComparator default Comparator
*/
const DefaultTasksComparator = "<"

/*
DefaultTasksAction default action
*/
const DefaultTasksAction = "scale_out"

/*
TypeTaskDefault default dock / namespace
*/
const TypeTaskDefault = "default"

/*
TypeTaskCOMPSs default dock / namespace
*/
const TypeTaskCOMPSs = "compss"

/*
TypeFTaskDefault default function task
*/
const TypeFTaskDefault = "function"

/*
SLANotDefined No SLA defined in the task
*/
const SLANotDefined = "NO_SLA"

/*
ClusterRUNNING ...
*/
const ClusterRUNNING = "RUNNING"

/*
ClusterERROR ...
*/
const ClusterERROR = "ERROR"

/*
ClusterDEPLOYING ...
*/
const ClusterDEPLOYING = "DEPLOYING"

/*
DefaultInfrQoSRule ...
*/
const DefaultInfrQoSRule = "Infr_Mem_2GB"
