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
	imec_db "atos/rotterdam/database/imec"
	"atos/rotterdam/globals/structs"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// parseMap
func parseMap(aMap map[string]interface{}, prefix string, b *strings.Builder) error {
	var sol string
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(concreteVal, prefix+key+".", b)
		case string:
			sol = key + "=" + concreteVal + " "
			b.WriteString(prefix)
			b.WriteString(sol)
			b.WriteRune(',')
		case float64:
			sol = key + "=" + strconv.FormatFloat(concreteVal, 'f', -1, 64) + " "
			b.WriteString(prefix)
			b.WriteString(sol)
			b.WriteRune(',')
		default:
			return errors.New("Unsupported value")
		}
	}

	return nil
}

// parsecustom
func parsecustom(aMap map[string]interface{}) string {
	b := &strings.Builder{}
	err := parseMap(aMap, "", b)
	if err == nil {
		return strings.TrimSuffix(b.String(), ",")
	}
	return "error"
}

/*
kubelessCall Function call
Examples:

	curl -L --data '{"Another": "Echo"}' --header "Content-Type:application/json"  http://10.0.2.15:8001/api/v1/namespaces/default/services/hello:http-function-port/proxy/
	curl -L --data '{"Another": "Echo"}' --header "Content-Type:application/json"  http://10.0.2.15:8001/api/v1/namespaces/default/services/get-python:http-function-port/proxy/
*/
func kubelessCall(namespace string, body string, f structs.CLASS_FUNCTION_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, string, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessCall] Generating 'deployment' json ...")

	jsonCall := map[string]string{"Event_Data": body}
	strTxt, _ := mapStrStrToString(jsonCall)
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessCall] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessCall] Calling function ...")
	status, results, err := common.HTTPPOSTStruct(
		"http://"+urls.GetHostIP(cluster)+":8001/api/v1/namespaces/"+namespace+"/services/"+f.ID+":http-function-port/proxy/",
		sec,
		jsonCall)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessCall] ERROR", err)
		return "", "", err
	}
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [kubelessCall] RESPONSE: OK")

	// results to string
	s := parsecustom(results)

	return strconv.Itoa(status), s, nil
}

/*
Call Callss a function
*/
func Call(w http.ResponseWriter, r *http.Request, f *structs.CLASS_FUNCTION_TASK) (string, string, error) {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Call] Calling serverless function ...")

	clusterInfr, _ := imec_db.GetCluster(f.Cluster)
	clusterID := ""
	clusterHost := ""
	if clusterInfr != nil {
		clusterID = clusterInfr.ID
		clusterHost = clusterInfr.HostIP
	}
	namespace := f.Dock
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Call] cluster id = " + clusterID + ", dock = " + namespace + ", host = " + clusterHost + "")

	// 1. body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	// 2. call
	status, res, err := kubelessCall(namespace, bodyString, *f, clusterInfr, false)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [Call] ERROR (1)", err)
		return "", res, err
	} else if status == "200" || status == "201" {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [Call] Function called with success")
		return "Function called with success", res, nil
	}

	err = errors.New("Function call failed. status = [" + status + "]")
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [Call] ERROR (2)", err)
	return "", "", err
}
