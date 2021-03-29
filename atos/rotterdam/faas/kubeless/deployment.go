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

package kubeless

import (
	urls "atos/rotterdam/caas/adapters"
	"atos/rotterdam/caas/common"
	db "atos/rotterdam/database/caas"
	imec_db "atos/rotterdam/database/imec"
	"atos/rotterdam/globals/structs"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
kubelessDeployment Function deployment
Examples:

	curl -H "Content-Type: application/yaml" -X POST http://10.0.2.15:8001/apis/kubeless.io/v1beta1/namespaces/default/functions --data-binary @function.yaml
	curl -H "Content-Type: application/json" -X POST http://10.0.2.15:8001/apis/kubeless.io/v1beta1/namespaces/default/functions -d @function.json
*/
func kubelessDeployment(namespace string, f structs.CLASS_FUNCTION_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [K8sDeployment] Generating 'deployment' json ...")
	jsonDepl := newFunctionJSONDeployment(f, 1) // returns *MicroK8sDeploymentStruct

	strTxt, _ := structToString(*jsonDepl)
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [K8sDeployment] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [K8sDeployment] Creating a new deployment in K8s cluster ...")
	status, _, err := common.HTTPPOSTStruct(
		"http://"+urls.GetHostIP(cluster)+":8001/apis/kubeless.io/v1beta1/namespaces/"+namespace+"/functions",
		sec,
		jsonDepl)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [K8sDeployment] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [K8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
Deploy Deploys a function
*/
func Deploy(w http.ResponseWriter, f *structs.CLASS_FUNCTION_TASK) (*structs.DB_TASK, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Deploy] Deploying new serverless function ...")

	clusterInfr, _ := imec_db.GetCluster(f.Cluster)
	clusterID := ""
	clusterHost := ""
	if clusterInfr != nil {
		clusterID = clusterInfr.ID
		clusterHost = clusterInfr.HostIP
	}
	namespace := f.Dock
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Deploy] cluster id = " + clusterID + ", dock = " + namespace + ", host = " + clusterHost + "")

	// 1. DEPLOYMENT /////
	status, err := kubelessDeployment(namespace, *f, clusterInfr, false)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [Deploy] ERROR (1)", err)
		return nil, err
	} else if status == "200" || status == "201" {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [Deploy] Function deployed with success")
		// 2. save to DB
		dbtask := &structs.DB_TASK{
			DbId:               structs.DB_TABLE_TASK,
			Id:                 f.ID,
			Name:               f.Name,
			NameSpace:          namespace,
			Type:               structs.DB_TASK_TYPE_FUNCTION,
			ClusterId:          clusterID,
			AgreementId:        strings.Replace(f.ID, "-", "_", -1),
			Url:                "http://" + f.ID + "." + clusterHost + ".xip.io",
			Status:             "Deployed",
			Replicas:           1,
			FunctionDefinition: *f}
		db.SetTaskValue(f.ID, *dbtask)

		return dbtask, nil
	}

	err = errors.New("Function creation failed. status = [" + status + "]")
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Deploy] ERROR (2)", err)
	return nil, err
}
