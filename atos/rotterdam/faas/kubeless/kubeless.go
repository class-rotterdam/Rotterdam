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
	"atos/rotterdam/globals/structs"
	"encoding/json"
	"log"
)

///////////////////////////////////////////////////////////////////////////////
// MICROK8S FUNCTION DEPLOYMENT STRUCT

/*
MicroK8sDeploymentStruct represents the YAML JSON needed to deploy a function using MicroK8s REST API
YAML example:

	apiVersion: kubeless.io/v1beta1
	kind: Function
	metadata:
		name: get-python
		namespace: default
		label:
			created-by: kubeless
			function: get-python
	spec:
		runtime: python2.7
		timeout: "180"
		handler: helloget.foo
		deps: ""
		checksum: sha256:d251999dcbfdeccec385606fd0aec385b214cfc74ede8b6c9e47af71728f6e9a
		function-content-type: text
		function: |
			def foo(event, context):
				return "hello world"
*/
type MicroK8sDeploymentStruct struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name      string `json:"name,omitempty"`
		NameSpace string `json:"namespace,omitempty"`
		Label     struct {
			CreatedBy string `json:"created-by,omitempty"`
			Function  string `json:"function,omitempty"`
		} `json:"label,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Runtime             string `json:"runtime,omitempty"`
		Timeout             string `json:"timeout,omitempty"`
		Handler             string `json:"handler,omitempty"`
		Deps                string `json:"deps,omitempty"`
		Checksum            string `json:"checksum,omitempty"`
		FunctionContentType string `json:"function-content-type,omitempty"`
		Function            string `json:"function,omitempty"`
	} `json:"spec,omitempty"`
}

/*
NewFunctionDeployment creates a new MicroK8s Deployment json
*/
func newFunctionJSONDeployment(ftask structs.CLASS_FUNCTION_TASK, replicas int) *MicroK8sDeploymentStruct {
	log.Println("Rotterdam > FAAS > Adapters > Kubeless [NewFunctionJSONDeployment] Creating new serverless function definition (JSON) ...")

	var jsonDeployment *MicroK8sDeploymentStruct
	jsonDeployment = new(MicroK8sDeploymentStruct)

	jsonDeployment.APIVersion = "kubeless.io/v1beta1"
	jsonDeployment.Kind = "Function"
	jsonDeployment.Metadata.Name = ftask.ID
	jsonDeployment.Metadata.NameSpace = ftask.Dock
	jsonDeployment.Metadata.Label.CreatedBy = "kubeless"
	jsonDeployment.Metadata.Label.Function = ftask.ID
	jsonDeployment.Spec.Runtime = ftask.Runtime
	jsonDeployment.Spec.Timeout = ftask.Timeout
	jsonDeployment.Spec.Handler = ftask.Handler
	jsonDeployment.Spec.Deps = ""
	jsonDeployment.Spec.Checksum = ""
	jsonDeployment.Spec.FunctionContentType = ftask.FunctionType
	jsonDeployment.Spec.Function = ftask.Function

	return jsonDeployment
}

/*
structToString Parses a struct to a string
*/
func structToString(ct MicroK8sDeploymentStruct) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [structToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
mapStrStrToString Parses a struct to a string
*/
func mapStrStrToString(ct map[string]string) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > FAAS > Adapters > Kubeless [mapStrStrToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}
