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
package lognotifier

import (
	assessment_model "SLALite/assessment/model"
	"SLALite/model"

	"github.com/labstack/gommon/log"
)

type LogNotifier struct {
}

func (n LogNotifier) NotifyViolations(agreement *model.Agreement, result *assessment_model.Result) {
	log.Info("Violation of agreement: " + agreement.Id)
	for k, v := range *result {
		if len(v.Violations) > 0 {
			log.Info("Failed guarantee: " + k)
			for _, vi := range v.Violations {
				log.Infof("Failed guarantee %v of agreement %s at %s", vi.Guarantee, vi.AgreementId, vi.Datetime)
			}
		}
	}
}
