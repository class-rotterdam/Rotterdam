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
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/database/imec"
	constants "atos/rotterdam/globals/constants"
	"log"
	"strconv"
)

/*
GetClusterIndex Returns the index of a cluster identified by its name / id value
*/
func GetClusterIndex(idCluster string) int {
	for index := range cfg.Config.Clusters {
		if idCluster == cfg.Config.Clusters[index].ID {
			log.Println("Rotterdam > CAAS > common [GetClusterIndex] Returning index (" + strconv.Itoa(index) + ") of cluster [" + idCluster + "] ...")
			return index
		}
	}

	log.Println("Rotterdam > CAAS > common [GetClusterIndex] WARNING Returning default cluster index (0) ...")
	return 0
}

/*
GetClusterType Returns the type of a cluster identified by its name / id value
*/
func GetClusterType(idCluster string) string {
	resp, err := imec_db.GetInfrByID(idCluster)
	if err == nil && len(resp) > 0 {
		return resp[0].Type
	}

	log.Println("Rotterdam > CAAS > common [GetClusterType] WARNING Returning default cluster type ...")
	return "Openshift"
}

/*
GetClusterDefaultDock Returns the type of a cluster identified by its name / id value
*/
func GetClusterDefaultDock(idCluster string) string {
	resp, err := imec_db.GetInfrByID(idCluster)
	if err == nil && len(resp) > 0 {
		return resp[0].DefaultDock
	}

	log.Println("Rotterdam > CAAS > common [GetClusterDefaultDock] WARNING Returning default cluster dock ...")
	return constants.DefaultDock
}
