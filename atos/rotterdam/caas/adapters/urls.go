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

package adapters

import (
	log "atos/rotterdam/common/logs"
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/database/imec"
)

// Paths

/*
PathKubernetesDeployment ...
*/
const PathKubernetesDeployment = "/apis/apps/v1/namespaces/default/deployments"

/*
PathKubernetesService ...
*/
const PathKubernetesService = "/api/v1/namespaces/default/services"

/*
PathKubernetesPodsApp ...
*/
const PathKubernetesPodsApp = "/api/v1/namespaces/default/pods?labelSelector=app="

/*
PathKubernetesPod ...
*/
const PathKubernetesPod = "/api/v1/namespaces/default/pods"

/*
PathOpenshiftRoute ...
*/
const PathOpenshiftRoute = "/apis/route.openshift.io/v1/namespaces/default/routes/route-"

/*
PathOpenshiftRoutes "/apis/route.openshift.io/v1/namespaces/"+namespace+"/routes"
*/
const PathOpenshiftRoutes = "/apis/route.openshift.io/v1/namespaces/default/routes"

// URLs
var defaultK8sEndPoint string
var defaultOpenshiftEndPoint string

/*
Initialize ...
*/
func Initialize() {
	defaultK8sEndPoint = cfg.Config.Clusters[0].KubernetesEndPoint
	defaultOpenshiftEndPoint = cfg.Config.Clusters[0].OpenshiftEndPoint
}

///////////////////////////////////////////////////////////////////////////////
// HOSTIP

/*
GetHostIP Return the Host IP or the value of  cfg.Config.Clusters[0].HostIP if cluster is nil
*/
func GetHostIP(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) string {
	ip := ""
	if cluster == nil {
		ip = cfg.Config.Clusters[0].HostIP
		log.Println(pathLOG + "urls [GetHostIP] ERROR Cluster is nil. Returning cluster[0] Host IP [" + ip + "] ...")
	} else {
		ip = cluster.HostIP
		log.Println(pathLOG + "urls [GetHostIP] Host IP: " + ip)
	}
	return ip
}

// KUBERNETES

// Deployments

/*
GetPathKubernetesDeployment Path = "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
*/
func GetPathKubernetesDeployment(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeployment] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesDeployment + "/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeployment] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeployment] URL: " + url)
	}
	return url
}

/*
GetPathKubernetesCreateDeployment Path = "/apis/apps/v1/namespaces/" + namespace + "/deployments"
*/
func GetPathKubernetesCreateDeployment(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments"
		log.Println(pathLOG + "urls [GetPathKubernetesCreateDeployment] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesDeployment
		log.Println(pathLOG + "urls [GetPathKubernetesCreateDeployment] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments"
		log.Println(pathLOG + "urls [GetPathKubernetesCreateDeployment] URL: " + url)
	}
	return url
}

/*
GetPathKubernetesDeleteDeployment Path = "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
*/
func GetPathKubernetesDeleteDeployment(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeleteDeployment] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesDeployment + "/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeleteDeployment] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task
		log.Println(pathLOG + "urls [GetPathKubernetesDeleteDeployment] URL: " + url)
	}
	return url
}

// Services

/*
GetPathKubernetesCreateService Path = "/api/v1/namespaces/" + namespace + "/services"
*/
func GetPathKubernetesCreateService(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string) string {
	url := ""
	if cluster == nil {
		log.Println(pathLOG + "urls [GetPathKubernetesCreateService] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
		url = defaultK8sEndPoint + "/api/v1/namespaces/" + namespace + "/services"
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesService
		log.Println(pathLOG + "urls [GetPathKubernetesCreateService] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/services"
		log.Println(pathLOG + "urls [GetPathKubernetesCreateService] URL: " + url)
	}
	return url
}

/*
GetPathKubernetesService Path = "/api/v1/namespaces/" + namespace + "/services/serv-" + task
*/
func GetPathKubernetesService(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/api/v1/namespaces/" + namespace + "/services/serv-" + task
		log.Println(pathLOG + "urls [GetPathKubernetesService] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesService + "/serv-" + task
		log.Println(pathLOG + "urls [GetPathKubernetesService] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/services/serv-" + task
		log.Println(pathLOG + "urls [GetPathKubernetesService] URL: " + url)
	}
	return url
}

// Scale

/*
GetPathKubernetesScaleDeployment Path = "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task + "/scale"
*/
func GetPathKubernetesScaleDeployment(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task + "/scale"
		log.Println(pathLOG + "urls [GetPathKubernetesScaleDeployment] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesDeployment + "/" + task + "/scale"
		log.Println(pathLOG + "urls [GetPathKubernetesScaleDeployment] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/apis/apps/v1/namespaces/" + namespace + "/deployments/" + task + "/scale"
		log.Println(pathLOG + "urls [GetPathKubernetesScaleDeployment] URL: " + url)
	}
	return url
}

// Pods

/*
GetPathKubernetesPodsApp Path = "/api/v1/namespaces/" + namespace + "/pods?labelSelector=app=" + task
*/
func GetPathKubernetesPodsApp(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/api/v1/namespaces/" + namespace + "/pods?labelSelector=app=" + task
		log.Println(pathLOG + "urls [GetPathKubernetesPodsApp] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		// 'default' namespace
		url = cluster.KubernetesEndPoint + PathKubernetesPodsApp + "/" + task + "/scale"
		log.Println(pathLOG + "urls [GetPathKubernetesPodsApp] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/pods?labelSelector=app=" + task
		log.Println(pathLOG + "urls [GetPathKubernetesPodsApp] URL: " + url)
	}
	return url
}

/*
GetPathKubernetesPod Path = "/api/v1/namespaces/" + namespace + "/pods/" + podName
*/
func GetPathKubernetesPod(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, podName string) string {
	url := ""
	if cluster == nil {
		url = defaultK8sEndPoint + "/api/v1/namespaces/" + namespace + "/pods/" + podName
		log.Println(pathLOG + "urls [GetPathKubernetesPod] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.KubernetesEndPoint + PathKubernetesPod + "/" + podName
		log.Println(pathLOG + "urls [GetPathKubernetesPod] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.KubernetesEndPoint + "/api/v1/namespaces/" + namespace + "/pods/" + podName
		log.Println(pathLOG + "urls [GetPathKubernetesPod] URL: " + url)
	}
	return url
}

///////////////////////////////////////////////////////////////////////////////
// OPENSHIFT

/*
GetPathOpenshiftRoute Path = "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes/route-" + task
PathOpenshiftRoute = "/apis/route.openshift.io/v1/namespaces/default/routes/route-"
*/
func GetPathOpenshiftRoute(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, task string) string {
	url := ""
	if cluster == nil {
		url = defaultOpenshiftEndPoint + "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes/route-" + task
		log.Println(pathLOG + "urls [GetPathOpenshiftRoute] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.OpenshiftEndPoint + PathOpenshiftRoute + task
		log.Println(pathLOG + "urls [GetPathOpenshiftRoute] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.OpenshiftEndPoint + "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes/route-" + task
		log.Println(pathLOG + "urls [GetPathOpenshiftRoute] URL: " + url)
	}
	return url
}

/*
GetPathOpenshiftRoutes Path = "/apis/route.openshift.io/v1/namespaces/"+namespace+"/routes"
PathOpenshiftRoutes "/apis/route.openshift.io/v1/namespaces/"+namespace+"/routes"
*/
func GetPathOpenshiftRoutes(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string) string {
	url := ""
	if cluster == nil {
		url = defaultOpenshiftEndPoint + "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes"
		log.Println(pathLOG + "urls [GetPathOpenshiftRoutes] ERROR Cluster is nil. Returning cluster[0] URL [" + url + "] ...")
	} else if len(namespace) == 0 {
		url = cluster.OpenshiftEndPoint + PathOpenshiftRoutes
		log.Println(pathLOG + "urls [GetPathOpenshiftRoutes] Using 'default' namespace. URL: " + url)
	} else {
		url = cluster.OpenshiftEndPoint + "/apis/route.openshift.io/v1/namespaces/" + namespace + "/routes"
		log.Println(pathLOG + "urls [GetPathOpenshiftRoutes] URL: " + url)
	}
	return url
}
