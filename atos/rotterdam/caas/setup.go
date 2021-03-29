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

package caas

import (
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/database/imec"
	constants "atos/rotterdam/globals/constants"
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"
)

/*
configStructToString Parses a struct to a string
*/
func configStructToString(ct cfg.ConfigurationSetUp) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > Config [configStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
structCheckConfig Checks if Config struct is valid (from json)
*/
func structCheckConfig(req *http.Request) (*cfg.ConfigurationSetUp, error) {
	log.Println("Rotterdam > CAAS > Config [structCheckConfig] Checking json object ...")

	decoder := json.NewDecoder(req.Body)
	var t cfg.ConfigurationSetUp
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > CAAS > Config [structCheckConfig] ERROR (1)", err)
		return nil, err
	}

	tStr, err := configStructToString(t)
	log.Println("Rotterdam > CAAS > Config [structCheckConfig] Parsed object (string): " + tStr)
	log.Println("Rotterdam > CAAS > Config [structCheckConfig] Sending parsed object ...")

	return &t, nil
}

/*
validateJSONConfig validates input (json) and generates a valid ConfigurationSetUp struct
*/
func validateJSONConfig(r *http.Request) (*cfg.ConfigurationSetUp, error) {
	log.Println("Rotterdam > CAAS > Config [validateJSONConfig] Parsing default / old json definition ...")
	cfgSetUp, err := structCheckConfig(r)
	if err == nil {
		// change values -> default values

		log.Println("Rotterdam > CAAS > Config [validateJSONConfig] Returning new configuration ...")
		return cfgSetUp, nil
	}
	return nil, err
}

// parseVarStr
func parseVarStr(vVal string, vConfig string) string {
	if len(vVal) == 0 {
		return vConfig
	}
	return vVal
}

// trimLastChar
func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

// parseVarURL
func parseVarURL(sURL string, vConfig string) string {
	if len(sURL) == 0 {
		return vConfig
	}

	for {
		lu := sURL[len(sURL)-1:]
		if lu != "/" {
			break
		} else {
			sURL = trimLastChar(sURL)
		}
	}

	if len(sURL) == 0 {
		return vConfig
	}
	return sURL
}

// getVarInt
func parseVarInt(vVal int, vConfig int) int {
	if vVal == 0 {
		return vConfig
	}
	return vVal
}

// getVarFloat
func parseVarFloat(vVal float64, vConfig float64) float64 {
	if vVal == 0.0 {
		return vConfig
	}
	return vVal
}

// setMainClusterID sets the ID of the main cluster
func setMainClusterID(totalClusters int) {
	if totalClusters == 0 {
		cfg.Config.Clusters[0].ID = constants.MainClusterID
	} else {
		found := false
		for i := 0; i < totalClusters; i++ {
			if cfg.Config.Clusters[0].ID == constants.MainClusterID {
				break
			}
		}
		if !found {
			cfg.Config.Clusters[0].ID = constants.MainClusterID
		}
	}
}

/*
UpdateConfig Update configuration values
*/
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### Update Rotterdam Configuration")

	log.Println("Rotterdam > CAAS > Config [UpdateConfig] (1.) Validating json ...")
	cfgSetUp, err := validateJSONConfig(r)
	if err == nil {
		log.Println("Rotterdam > CAAS > Config [UpdateConfig] Updating configuration values from ConfigurationSetUp ...")

		totalClusters := len(cfgSetUp.Clusters)
		cfg.Config.Clusters = make([]cfg.ConfigurationCluster, totalClusters)
		for i := 0; i < totalClusters; i++ {
			/*
				ID                            string
				Name                          string
				Description                   string
				Type                          string
				SO                            string
				DefaultDock                   string
				OpenshiftOauthToken           string
				KubernetesEndPoint            string
				OpenshiftEndPoint             string
				SLALiteEndPoint               string
				PrometheusPushgatewayEndPoint string
				HostIP                        string
				HostPort                      int
				User                          string
				Password                      string
				KeyFile                       string
			*/
			cfg.Config.Clusters[i].ID = parseVarStr(cfgSetUp.Clusters[i].ID, "undefined")
			cfg.Config.Clusters[i].Name = parseVarStr(cfgSetUp.Clusters[i].Name, "undefined")
			cfg.Config.Clusters[i].Description = parseVarStr(cfgSetUp.Clusters[i].Description, "undefined")
			cfg.Config.Clusters[i].Type = parseVarStr(cfgSetUp.Clusters[i].Type, constants.DefaultType)
			cfg.Config.Clusters[i].DefaultDock = parseVarStr(cfgSetUp.Clusters[i].DefaultDock, constants.DefaultDock)
			cfg.Config.Clusters[i].SO = parseVarStr(cfgSetUp.Clusters[i].SO, "undefined")
			cfg.Config.Clusters[i].HostIP = cfgSetUp.Clusters[i].HostIP
			cfg.Config.Clusters[i].KubernetesEndPoint = parseVarURL(cfgSetUp.Clusters[i].KubernetesEndPoint, "undefined")
			cfg.Config.Clusters[i].OpenshiftEndPoint = parseVarURL(cfgSetUp.Clusters[i].OpenshiftEndPoint, "undefined")
			cfg.Config.Clusters[i].SLALiteEndPoint = parseVarURL(cfgSetUp.Clusters[i].SLALiteEndPoint, "undefined")
			cfg.Config.Clusters[i].PrometheusPushgatewayEndPoint = parseVarURL(cfgSetUp.Clusters[i].PrometheusPushgatewayEndPoint, "undefined")
			cfg.Config.Clusters[i].OpenshiftOauthToken = parseVarStr(cfgSetUp.Clusters[i].OpenshiftOauthToken, "undefined")
			cfg.Config.Clusters[i].HostPort = parseVarInt(cfgSetUp.Clusters[i].HostPort, 80)
			cfg.Config.Clusters[i].User = parseVarStr(cfgSetUp.Clusters[i].User, "undefined")
			cfg.Config.Clusters[i].Password = parseVarStr(cfgSetUp.Clusters[i].Password, "undefined")
			cfg.Config.Clusters[i].KeyFile = parseVarStr(cfgSetUp.Clusters[i].KeyFile, "undefined")
		}

		// checks / sets the ID of the main cluster
		setMainClusterID(totalClusters)

		/*
			MaxReplicas int
			MinReplicas int
			MaxAllowed  int
			ScaleFactor float64
			Value       int
			Comparator  string
			Action      string
		*/
		cfg.Config.Tasks.MaxReplicas = parseVarInt(cfgSetUp.Tasks.MaxReplicas, constants.DefaultTasksMaxReplicas)
		cfg.Config.Tasks.MinReplicas = parseVarInt(cfgSetUp.Tasks.MinReplicas, constants.DefaultTasksMinReplicas)
		cfg.Config.Tasks.MaxAllowed = parseVarInt(cfgSetUp.Tasks.MaxAllowed, constants.DefaultTasksMaxAllowed)
		cfg.Config.Tasks.ScaleFactor = parseVarFloat(cfgSetUp.Tasks.ScaleFactor, constants.DefaultTasksScaleFactor)
		cfg.Config.Tasks.Value = parseVarInt(cfgSetUp.Tasks.Value, constants.DefaultTasksValue)
		cfg.Config.Tasks.Comparator = parseVarStr(cfgSetUp.Tasks.Comparator, constants.DefaultTasksComparator)
		cfg.Config.Tasks.Action = parseVarStr(cfgSetUp.Tasks.Action, constants.DefaultTasksAction)

		// SLA
		cfg.Config.SLAs.DefaultInfrQoSRule = parseVarStr(cfgSetUp.SLAs.DefaultInfrQoSRule, constants.DefaultInfrQoSRule)

		// reset infrs from imec database
		imec_db.ResetDB()

		// save to DB: adding clusters to DB
		_, err := imec_db.AddConfigInfrsToDB()
		if err != nil {
			// error
			json.NewEncoder(w).Encode(cfg.ResponseConfig{
				Resp:        "error",
				Method:      "UpdateConfig",
				Message:     "An error occurred while saving new configuration values to IMEC database",
				CaaSVersion: cfg.Config.CaaSVersion})
		} else {
			// return response
			json.NewEncoder(w).Encode(cfg.ResponseConfig{
				Resp:        "ok",
				Method:      "UpdateConfig",
				Message:     "Configuration updated",
				CaaSVersion: cfg.Config.CaaSVersion,
				Content:     cfg.Config})
		}
	} else {
		// error
		json.NewEncoder(w).Encode(cfg.ResponseConfig{
			Resp:        "error",
			Method:      "UpdateConfig",
			Message:     "An error occurred while reading / setting the new configuration values",
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

/*
GetCurrentConfig Get current configuration values
*/
func GetCurrentConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > Config > Getting configuration values from Configuration ...")

	// return response
	json.NewEncoder(w).Encode(cfg.ResponseConfig{
		Resp:        "ok",
		Method:      "GetCurrentConfig",
		Message:     "Current configuration",
		CaaSVersion: cfg.Config.CaaSVersion,
		Content:     cfg.Config})
}
