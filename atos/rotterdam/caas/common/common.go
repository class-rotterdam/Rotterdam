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

package common

import (
	structs "atos/rotterdam/caas/common/structs"
	"encoding/json"
	"log"
	cfg "atos/rotterdam/config"
)

/*
 * CommClassQoSTemplateListToString: Parses a struct to a string
 */
 func CommClassQoSTemplateListToString(ct cfg.CLASS_QOS_TEMPLATE_LIST) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommClassQoSTemplateListToString] json object / struct to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [CommClassQoSTemplateListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommClassStructToString: Parses a struct to a string
 */
func CommClassStructToString(ct structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommClassStructToString] json object / struct to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [CommClassStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommDeploymentStructToString: Parses a struct to a string
 */
func CommDeploymentStructToString(ct structs.K8S_DEPLOYMENT) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommDeploymentStructToString] json object / struct [K8S_DEPLOYMENT] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDeploymentStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommServiceStructToString: Parses a struct to a string
 */
func CommServiceStructToString(ct structs.K8S_SERVICE) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommServiceStructToString] json object / struct [K8S_SERVICE] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommServiceStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommRouteStructToString: Parses a struct to a string
 */
func CommRouteStructToString(ct structs.K8S_ROUTE) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommRouteStructToString] json object / struct [K8S_ROUTE] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommRouteStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}


/*
 * CommPatchPodsListToString: Parses a struct to a string
 */
 func CommPatchPodsListToString(ct []structs.K8S_POD_PATCH_LINE) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommPatchPodsListToString] json object / struct [K8S_POD_PATCH_LINE] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommPatchPodsListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}


/*
 * CommDbTaskStructToString: Parses a struct to a string
 */
func CommDbTaskStructToString(ct structs.DB_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommDbTaskStructToString] json object / struct [DB_TASK] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDbTaskStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommDbTaskQoSStructToString: Parses a struct to a string
 */
func CommDbTaskQoSStructToString(ct structs.DB_TASK_QOS) (string, error) {
	log.Println("Rotterdam > CAAS > common [CommDbTaskQoSStructToString] json object / struct [DB_TASK_QOS] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDbTaskQoSStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
 * CommStringToDbTaskStruct: Parses a string to a struct of type DB_TASK
 */
func CommStringToDbTaskStruct(ct string) (*structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > common [CommStringToDbTaskStruct] string tp json object / struct [DB_TASK]  ...")

	data := &structs.DB_TASK{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToDbTaskStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
 * CommStringToDbTaskQoSStruct: Parses a string to a struct of type DB_TASK_QOS
 */
func CommStringToDbTaskQoSStruct(ct string) (*structs.DB_TASK_QOS, error) {
	log.Println("Rotterdam > CAAS > common [CommStringToDbTaskQoSStruct] string tp json object / struct [DB_TASK_QOS]  ...")

	data := &structs.DB_TASK_QOS{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToDbTaskQoSStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
 * CommStringToDbTaskStruct: Parses a string to a struct of type K8S_SCALE
 */
func CommStringToK8S_SCALE(ct string) (*structs.K8S_SCALE, error) {
	log.Println("Rotterdam > CAAS > common [CommStringToK8S_SCALE] string tp json object / struct [K8S_SCALE]  ...")

	data := &structs.K8S_SCALE{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToK8S_SCALE] ERROR", err)
		return data, err
	}

	return data, nil
}
