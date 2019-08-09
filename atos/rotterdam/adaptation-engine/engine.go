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

package adaptation_engine

import (
	structs "atos/rotterdam/caas/common/structs"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Configuration
type Configuration struct {
	Adapter        string
	KubeConfigPath string
}

var Config Configuration

// CaaS Response
type ResponseCaaS struct {
	Resp    string `json:"resp,omitempty"`
	Message string `json:"message,omitempty"`
}

// Default Function for not implemented calls
func NotImplementedFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > Adaptation-Engine > engine [NotImplementedFunc] not implemented")
	json.NewEncoder(w).Encode(ResponseCaaS{Resp: "ok", Message: "not implemented"})
}

// Process Violations (from SLALite)
func ProcessViolation(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > Adaptation-Engine > engine [ProcessViolation] Reading params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > Adaptation-Engine > engine [ProcessViolation] guarantee..." + params["guarantee"])
	log.Println("Rotterdam > Adaptation-Engine > engine [ProcessViolation] name..." + params["name"]) // task name

	if r.Body == nil {
		log.Println("Rotterdam > Adaptation-Engine > engine [ProcessViolation] ERROR body is nil")
	} else {
		var u structs.ViolationInfo
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println("Rotterdam > Adaptation-Engine > engine [ProcessViolation] ERROR processing violation from SLA: ", err)
		} else {
			Process(w, u)
		}
	}
}
