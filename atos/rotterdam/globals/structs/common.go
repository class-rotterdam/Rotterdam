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
	cfg "atos/rotterdam/config"
	"encoding/json"
	"log"
)

/*
CommClassQoSTemplateListToString Parses a struct to a string
*/
func CommClassQoSTemplateListToString(ct cfg.CLASS_QOS_TEMPLATE_LIST) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common [CommClassQoSTemplateListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommClassStructToString Parses a struct to a string
*/
func CommClassStructToString(ct CLASS_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommClassStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommClassCOMPSsStructToString Parses a struct to a string
*/
func CommClassCOMPSsStructToString(ct CLASS_COMPSS_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommClassCOMPSsStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommFuncClassStructToString Parses a struct to a string
*/
func CommFuncClassStructToString(ct CLASS_FUNCTION_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommFuncClassStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDeploymentStructToString Parses a struct to a string
*/
func CommDeploymentStructToString(ct K8S_DEPLOYMENT) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommDeploymentStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommServiceStructToString Parses a struct to a string
*/
func CommServiceStructToString(ct K8S_SERVICE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommServiceStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommRouteStructToString Parses a struct to a string
*/
func CommRouteStructToString(ct K8S_ROUTE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommRouteStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommPatchPodsListToString Parses a struct to a string
*/
func CommPatchPodsListToString(ct []K8S_POD_PATCH_LINE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommPatchPodsListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDbTaskStructToString Parses a struct to a string
*/
func CommDbTaskStructToString(ct DB_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommDbTaskStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDbTaskQoSStructToString Parses a struct to a string
*/
func CommDbTaskQoSStructToString(ct DB_TASK_QOS) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommDbTaskQoSStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommSLAStructToString Parses a struct to a string
*/
func CommSLAStructToString(ct SLA_AGREEMENT) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommSLAStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommStringToDbTaskStruct Parses a string to a struct of type DB_TASK
*/
func CommStringToDbTaskStruct(ct string) (*DB_TASK, error) {
	data := &DB_TASK{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommStringToDbTaskStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
CommStringToDbTaskQoSStruct Parses a string to a struct of type DB_TASK_QOS
*/
func CommStringToDbTaskQoSStruct(ct string) (*DB_TASK_QOS, error) {
	data := &DB_TASK_QOS{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommStringToDbTaskQoSStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
CommStringToK8S_SCALE Parses a string to a struct of type K8S_SCALE
*/
func CommStringToK8S_SCALE(ct string) (*K8S_SCALE, error) {
	data := &K8S_SCALE{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > GLOBALS > structs > common[CommStringToK8S_SCALE] ERROR", err)
		return data, err
	}

	return data, nil
}
