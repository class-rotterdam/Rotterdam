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

package monitoring

import (
	log "atos/rotterdam/common/logs"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// path used in logs
const pathLOG string = "Rotterdam >  Adaptation-Engine > Monitoring : "

/*
METRIC ...
*/
const METRIC = "node_memory_MemFree_bytes"

/*
MAXVALUE ...
*/
const MAXVALUE = 800 // MB

/*
ClusterMonitoring ...
*/
type ClusterMonitoring struct {
	DbIfrClusterID     string  `json:"dbIfrClusterID,omitempty"` // DB_INFRASTRUCTURE_CLUSTER from database/imec/db
	PrometheusEndPoint string  `json:"prometheusEndPoint,omitempty"`
	Metric1            string  `json:"metric1,omitempty"`
	Value1             int     `json:"value1,omitempty"`
	Time1              float64 `json:"time1,omitempty"`
	Status             string  `json:"status,omitempty"`
}

// Prometheus Queries
/*	{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				"metric": {},
				"value": [
					"0": 1787216387123 		// TIME
					"1": "973213123123"		// VALUE
				]
			]
		}
	}

*/

/*
Metric ...
*/
type Metric struct {
	Name     string `json:"__name__"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

/*
Result ...
*/
type Result struct {
	Metrics Metric        `json:"metric"`
	Values  []interface{} `json:"value"`
}

/*
Data ...
*/
type Data struct {
	ResultType string   `json:"resultType"`
	Results    []Result `json:"result"`
}

/*
Query ...
*/
type Query struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

/*
ClustersMonitoringStatus status of clusters
*/
var ClustersMonitoringStatus []ClusterMonitoring

/*
mutex ...
*/
var mutex = &sync.Mutex{}

/*
quit chan used to start-stop ticks used while checking clusters
*/
var quit chan struct{}

/*
Init fucntion
*/
func init() {
	log.Println(pathLOG + "[init] Initializing 'quit' channel ...")
	quit = make(chan struct{})

	log.Println(pathLOG + "[init] Initializing 'ClustersMonitoringStatus' struct ...")
	ClustersMonitoringStatus = make([]ClusterMonitoring, 0)
}

/*
checkClusters check infrastructure metrics from all clusters
*/
func checkClusters() {
	mutex.Lock()
	for _, c := range ClustersMonitoringStatus {
		// Check metric
		res, err := checkMetric(&c)
		if err != nil {
			log.Error(pathLOG+"[CheckClusters] ERROR checking metric: ", err)
		}
		// Change status
		c.Status = res
	}
	mutex.Unlock()
}

/*
checkMetric get metric from Prometheus and evaluate it
*/
func checkMetric(c *ClusterMonitoring) (string, error) {
	log.Println(pathLOG + "[checkMetric] Checking metrics from cluster [" + c.DbIfrClusterID + "] [" + c.PrometheusEndPoint + "] ...")
	// GET: URLPrometheus+"/api/v1/query?query=sum("+c.Metric1+")"
	i := strings.Index(c.PrometheusEndPoint, "https")
	isHTTPS := false
	if i > -1 {
		isHTTPS = true
	}

	_, result, err := GETHttpString(c.PrometheusEndPoint+"/api/v1/query?query=sum("+c.Metric1+")", isHTTPS)

	if err != nil {
		log.Error(pathLOG+"[checkMetric] Metric: '"+c.Metric1+"', ERROR executing query [GET "+c.PrometheusEndPoint+"/api/v1/query?query=sum("+c.Metric1+")]", err)
	} else {
		// RESPONSE:
		log.Debug(pathLOG + "[checkMetric] Metric: '" + c.Metric1 + "', Results: " + result) //str)

		res := Query{}
		json.Unmarshal([]byte(result), &res)

		/*	{
				"status": "success",
				"data": {
					"resultType": "vector",
					"result": [
						"metric": {},
						"value": [
							"0": 1787216387123 		// TIME
							"1": "973213123123"		// VALUE
						]
					]
				}
			}
		*/
		if len(res.Data.Results) > 0 {
			// Time:
			timeFloat := res.Data.Results[0].Values[0].(float64)
			// Value:
			value := 0

			strVal := res.Data.Results[0].Values[1].(string)
			n, err := strconv.ParseInt(strVal, 10, 64)
			if err == nil {
				value = int(n / 1024 / 1024) // to MB
			}

			c.Value1 = value
			c.Time1 = timeFloat
		}

		log.Debug(pathLOG + "[checkMetric] Metric: '" + c.Metric1 + "', c.Value1: " + strconv.Itoa(c.Value1))
	}

	// process value => c.Value1 > MAXVALUE
	if c.Value1 > 0 && c.Value1 < MAXVALUE {
		log.Println(pathLOG + "[checkMetric] c.Value1 > 0 && c.Value1 [" + strconv.Itoa(c.Value1) + "] < MAXVALUE [" + strconv.Itoa(MAXVALUE) + "] ==> Status=NO")
		c.Status = "no"
	} else {
		log.Println(pathLOG + "[checkMetric] Status=OK c.Value1 [" + strconv.Itoa(c.Value1) + "] > MAXVALUE [" + strconv.Itoa(MAXVALUE) + "]")
		c.Status = "ok"
	}

	return "", nil
}

/*
StartCheckingClusters ...
*/
func StartCheckingClusters() {
	log.Println(pathLOG + "[StopCheckingClusters] Starting assessment process ...")
	ticker := time.NewTicker(150 * time.Second)
	//quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkClusters()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

/*
StopCheckingClusters ...
*/
func StopCheckingClusters() {
	log.Println(pathLOG + "[StopCheckingClusters] Finishing assessment process ...")
	close(quit)
}

/*
GetClusterStatus ...
*/
func GetClusterStatus(idCluster string) string {
	log.Println(pathLOG + "[GetClusterStatus] Gettin cluster status [" + METRIC + "] ...")

	for _, c := range ClustersMonitoringStatus {
		if c.DbIfrClusterID == idCluster {
			return c.Status
		}
	}

	return "Not-Found"
}

/*
GetAvailableCluster ...
*/
func GetAvailableCluster() string {
	log.Println(pathLOG + "[GetAvailableCluster] Gettin cluster status [" + METRIC + "] ...")

	for _, c := range ClustersMonitoringStatus {
		if c.Status == "ok" {
			return c.DbIfrClusterID
		}
	}

	return "Not-Found"
}

/*
ClusterStatusInfoResponse ...
*/
type ClusterStatusInfoResponse struct {
	Resp    string `json:"resp"`
	Method  string `json:"method"`
	Message string `json:"message"`
}

/*
GetClusterStatusInfo ...
*/
func GetClusterStatusInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println(pathLOG + "[GetClusterStatusInfo] Getting infrastructure " + params["id"] + " ...")
	idCluster := params["id"]

	log.Debug(pathLOG + "[GetClusterStatusInfo] Gettin cluster status [" + METRIC + "] ...")

	res := "Not-Found"
	for _, c := range ClustersMonitoringStatus {
		if c.DbIfrClusterID == idCluster {
			res = c.Status
		}
	}

	json.NewEncoder(w).Encode(ClusterStatusInfoResponse{
		Resp:    "ok",
		Method:  "GetClusterStatusInfo",
		Message: res})
}

/*
AddCluster ...
*/
func AddCluster(idCluster string, prometheusURL string) {
	log.Println(pathLOG + "[AddCluster] Adding cluster [" + idCluster + "] to monitoring process ...")

	var c *ClusterMonitoring
	c = new(ClusterMonitoring)
	c.DbIfrClusterID = idCluster
	c.PrometheusEndPoint = prometheusURL
	c.Metric1 = "node_memory_MemFree_bytes"
	//c.Value1 = ""
	//c.Time1 = ""
	c.Status = "ok"

	mutex.Lock()
	ClustersMonitoringStatus = append(ClustersMonitoringStatus, *c)
	mutex.Unlock()

	log.Println(pathLOG + "[AddCluster] cluster [" + idCluster + "] added to monitoring process.")
}

/*
RemoveCluster ...
*/
func RemoveCluster(idCluster string) error {
	log.Println(pathLOG + "[RemoveCluster] Adding intial infrastructures / orchestrators to DB ...")

	return nil
}
