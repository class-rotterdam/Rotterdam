/*
Copyright 2017 Atos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

CLASS Project: https://class-project.eu/

@author: ATOS
*/
package rotterdamnotifier

import (
	assessment_model "SLALite/assessment/model"
	"SLALite/model"
	"bytes"
	"encoding/json"
	"os"
	"strconv"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// ViolationInfo
type ViolationInfo struct {
	Agreement_id string
	Status       string
	Client_id    string
	Client_name  string
	Guarantee    string
}

// ExportNotifier
type ExportNotifier struct {
}

// Prometheus (ENV) variables
var UrlRotterdam string

// initialization function
func init() {
	log.SetLevel(log.DebugLevel)
	log.Println("SLALite > rotterdamnotifier > init > Initializing Prometheus notifier ...")

	// ENVIRONMENT VARIABLES
	// UrlRotterdam
	if os.Getenv("UrlRotterdam") != "" {
		log.Println("SLALite > rotterdamnotifier > init > Setting 'UrlRotterdam' value ... " + os.Getenv("UrlRotterdam"))
		UrlRotterdam = os.Getenv("UrlRotterdam")
	} else {
		UrlRotterdam = "http://rotterdam-caas.192.168.7.28.xip.io"
	}
	log.Println("SLALite > rotterdamnotifier > init > 'UrlRotterdam' = " + UrlRotterdam)
}

// NotifyViolations: send violations to Rotterdam Rules Engine
func (n ExportNotifier) NotifyViolations(agreement *model.Agreement, result *assessment_model.Result) {
	log.Info("SLALite > rotterdamnotifier > NotifyViolations > Checking agreement [" + agreement.Id + "] ...")
	log.Info("SLALite > rotterdamnotifier > NotifyViolations > Total violations = " + strconv.Itoa(len(result.GetViolations())))
	for k, v := range *result {
		if len(v.Violations) > 0 {
			log.Info("SLALite > rotterdamnotifier > NotifyViolations > Failed guarantee [" + k + "] from client [" + agreement.Details.Client.Id + "]")
			for _, vi := range v.Violations {
				log.Infof("SLALite > rotterdamnotifier > NotifyViolations > Failed guarantee %v of agreement %s at %s", vi.Guarantee, vi.AgreementId, vi.Datetime)
			}

			u := ViolationInfo{Agreement_id: agreement.Id,
				Status:      "violation",
				Client_id:   agreement.Details.Client.Id,
				Client_name: agreement.Details.Client.Name,
				Guarantee:   k}
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(u)
			_, err := http.Post(UrlRotterdam+"/api/v1/sla/tasks/"+agreement.Details.Client.Id+"/guarantee/"+k,
				"application/json; charset=utf-8",
				b)

			if err != nil {
				log.Errorf("SLALite > rotterdamnotifier > NotifyViolations > Error sending violations")
			} else {
				log.Infof("SLALite > rotterdamnotifier > NotifyViolations > Violation sent to [" + UrlRotterdam + "/api/v1/sla/tasks/" + agreement.Details.Client.Id + "/guarantee/" + k + "]")
			}
		}
	}
}
