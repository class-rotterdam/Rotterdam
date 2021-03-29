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

package common

import (
	log "atos/rotterdam/common/logs"
	structs "atos/rotterdam/globals/structs"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

// path used in logs
const pathLOG string = "Rotterdam > CAAS > Adapters > Common : "

/*
SERV_POD_PORT struct
*/
type SERV_POD_PORT struct {
	Name       string `json:"name,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Port       int    `json:"port,omitempty"`
	TargetPort int    `json:"targetPort,omitempty"`
}

/*
SERV_POD struct
*/
type SERV_POD struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Ports []SERV_POD_PORT `json:"ports,omitempty"`
	} `json:"spec,omitempty"`
}

// port number used to map compss applications
var (
	mu    sync.Mutex // guards balance
	rport int
)

// init
func init() {
	mu.Lock()
	rport = 25000
	mu.Unlock()
}

// set rport value
func setRPort() {
	mu.Lock()
	rport = rport + 1
	mu.Unlock()
}

// read rport value
func readRPort() int {
	mu.Lock()
	b := rport
	mu.Unlock()
	return b
}

/*
NewRPort generate and read rport value
*/
func NewRPort() int {
	mu.Lock()
	log.Debug(pathLOG + "[newRPort] Getting new PORT ...")
	rport = rport + 1
	b := rport
	log.Debug(pathLOG + "[newRPort] New PORT = " + strconv.Itoa(b))
	mu.Unlock()
	return b
}

/*
StringToServPodStruct Parses a string to a struct of type SERV_POD
*/
func StringToServPodStruct(ct string) (*SERV_POD, error) {
	log.Println(pathLOG + "[StringToServPodStruct] string tp json object / struct [SERV_POD]  ...")

	data := &SERV_POD{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Error(pathLOG+"[StringToServPodStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
GetMainPort get main ports of the application
*/
func GetMainPort(task structs.CLASS_TASK) (int, string) {
	// ports
	mainPort := 0
	mainProtocol := ""
	for _, contElement := range task.Containers {
		for _, portElement := range contElement.Ports {
			if mainPort == 0 {
				mainPort = portElement.ContainerPort
				mainProtocol = strings.ToUpper(portElement.Protocol)
			}
		}
	}

	if mainPort == 0 {
		log.Debug(pathLOG + "[getMainPort] ERROR getting main port. Returning 80 ...")
		mainPort = 80
		mainProtocol = "TCP"
	}

	log.Debug(pathLOG + "[getMainPort] main port = " + strconv.Itoa(mainPort))
	log.Debug(pathLOG + "[getMainPort] main protocol = " + mainProtocol)
	return mainPort, mainProtocol
}
