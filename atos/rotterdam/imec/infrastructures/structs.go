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

package infrastructures

import (
	log "atos/rotterdam/common/logs"
	db "atos/rotterdam/database/imec"
	"encoding/json"
	"net/http"
)

// path used in logs
const pathLOG string = "Rotterdam > Imec > infrastructures : "

/*
INFRASTRUCTURE represents a host / edge device / cluster access point
*/
type INFRASTRUCTURE struct {
	ID                            string `json:"id,omitempty"`
	Name                          string `json:"name,omitempty"`
	Description                   string `json:"description,omitempty"`
	Type                          string `json:"type,omitempty"`        // EDGE: 'microk8s', 'k3s'; CLUSTER: 'openshift', 'kubernetes'
	SO                            string `json:"so,omitempty"`          // S.O.: 'ubuntu18', 'ubuntu16', 'centos'
	DefaultDock                   string `json:"defaultDock,omitempty"` // 'class', 'default'
	HostIP                        string `json:"hostIP,omitempty"`
	HostPort                      int    `json:"hostPort,omitempty"`
	User                          string `json:"user,omitempty"`
	Password                      string `json:"password,omitempty"`
	KeyFile                       string `json:"keyFile,omitempty"`
	OpenshiftOauthToken           string `json:"openshiftOauthToken,omitempty"`
	KubernetesEndPoint            string `json:"kubernetesEndPoint,omitempty"`
	OpenshiftEndPoint             string `json:"openshiftEndPoint,omitempty"`
	SLALiteEndPoint               string `json:"slaliteEndPoint,omitempty"`
	PrometheusPushgatewayEndPoint string `json:"prometheusPushgatewayEndPoint,omitempty"`
	PrometheusEndPoint            string `json:"prometheusEndPoint,omitempty"`
	Status                        string `json:"status,omitempty"`
}

/*
infrJSONToString Parses a struct to a string
*/
func infrJSONToString(ct INFRASTRUCTURE) (string, error) {
	log.Println(pathLOG + "[infrJSONToString] json object / struct [INFRASTRUCTURE] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Error(pathLOG+"[infrJSONToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
ValidateStructInfrJSON Validate structure INFRASTRUCTURE
*/
func ValidateStructInfrJSON(req *http.Request) (*INFRASTRUCTURE, error) {
	log.Println(pathLOG + "[ValidateStructInfrJSON] Checking json object ...")

	decoder := json.NewDecoder(req.Body)
	var t INFRASTRUCTURE
	err := decoder.Decode(&t)
	if err != nil {
		log.Error(pathLOG+"[ValidateStructInfrJSON] ERROR (1)", err)
		return nil, err
	}

	tStr, err := infrJSONToString(t)
	log.Debug(pathLOG + "[ValidateStructInfrJSON] Parsed object (string): " + tStr)
	log.Debug(pathLOG + "[ValidateStructInfrJSON] Sending parsed object ...")

	return &t, nil
}

/*
InfrJSONToDBObj Transforms a INFRASTRUCTURE into a DB_INFRASTRUCTURE_CLUSTER object
*/
func InfrJSONToDBObj(i INFRASTRUCTURE, newID string) *db.DB_INFRASTRUCTURE_CLUSTER {
	log.Println(pathLOG + "[InfrJSONToDBObj] Transforming INFRASTRUCTURE object ...")

	var dbObj *db.DB_INFRASTRUCTURE_CLUSTER
	dbObj = new(db.DB_INFRASTRUCTURE_CLUSTER)

	dbObj.ID = newID
	dbObj.Name = i.Name
	dbObj.Description = i.Description
	dbObj.HostIP = i.HostIP
	dbObj.HostPort = i.HostPort
	dbObj.Type = i.Type
	dbObj.SO = i.SO
	dbObj.User = i.User
	dbObj.Password = i.Password
	dbObj.KeyFile = i.KeyFile
	dbObj.OpenshiftOauthToken = i.OpenshiftOauthToken
	dbObj.KubernetesEndPoint = i.KubernetesEndPoint
	dbObj.OpenshiftEndPoint = i.OpenshiftEndPoint
	dbObj.SLALiteEndPoint = i.SLALiteEndPoint
	dbObj.PrometheusPushgatewayEndPoint = i.PrometheusPushgatewayEndPoint
	dbObj.PrometheusEndPoint = i.PrometheusEndPoint
	dbObj.Status = i.Status

	dbObj.TableID = db.DB_TABLE_INFRASTRUCTURE_CLUSTER
	if len(i.DefaultDock) == 0 {
		dbObj.DefaultDock = "default"
	} else {
		dbObj.DefaultDock = i.DefaultDock
	}

	log.Println(pathLOG + "[InfrJSONToDBObj] Returning DB_INFRASTRUCTURE_CLUSTER object ...")
	return dbObj
}
