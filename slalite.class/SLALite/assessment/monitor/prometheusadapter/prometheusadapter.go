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

/*
Package prometheusadapter provides an integration of SLALite in prometheus.
It just returns a query to see if the service is up.

Usage:
	ma := prometheusadapter.New()
	ma.Initialize(&agreement)
	for _, gt := range gts {
		for values := ma.NextValues(gt); values != nil; values = ma.NextValues(gt) {
			...
		}
	}
*/
package prometheusadapter

import (
	assessment_model "SLALite/assessment/model"
	"SLALite/assessment/monitor"
	"SLALite/model"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Prometheus (ENV) variables

/*
URLPrometheus Main  Prometheus Endpoint
*/
var URLPrometheus string // endpoint

/*
MetricsPrometheus list of Prometheus metrics to be analyzed
*/
var MetricsPrometheus []string // metrics

/*
PrometheusEndPoints list of Prometheus Endponts
*/
var PrometheusEndPoints []string // endpoint

// initialization function
func init() {
	log.SetLevel(log.DebugLevel)
	log.Println("SLALite > prometheusadapter > init > Initializing Prometheus adapter ...")

	// ENVIRONMENT VARIABLES
	// UrlPrometheus
	if os.Getenv("UrlPrometheus") != "" {
		log.Println("SLALite > prometheusadapter > init > Setting 'UrlPrometheus' value ... " + os.Getenv("UrlPrometheus"))
		URLPrometheus = os.Getenv("UrlPrometheus")
	} else {
		URLPrometheus = "http://prometheus.192.168.7.28.xip.io"
	}
	log.Println("SLALite > prometheusadapter > init > 'UrlPrometheus' = " + URLPrometheus)

	// MetricsPrometheus
	if os.Getenv("MetricsPrometheus") != "" {
		log.Println("SLALite > prometheusadapter > init > Setting 'MetricsPrometheus' value ... " + os.Getenv("MetricsPrometheus"))
		MetricsPrometheus = strings.Split(os.Getenv("MetricsPrometheus"), ",")
	} else {
		MetricsPrometheus = []string{"go_memstats_frees_total"}
	}
	log.Println("SLALite > prometheusadapter > init > 'MetricsPrometheus' = " + strings.Join(MetricsPrometheus, " / "))

	// PrometheusEndPoints (used to listen from multiple Prometheus sources)
	if os.Getenv("PrometheusEndPoints") != "" {
		log.Println("SLALite > prometheusadapter > init > Setting 'PrometheusEndPoints' value ... " + os.Getenv("PrometheusEndPoints"))
		PrometheusEndPoints = strings.Split(os.Getenv("PrometheusEndPoints"), ",")
	} else {
		PrometheusEndPoints = []string{URLPrometheus}
	}
	log.Println("SLALite > prometheusadapter > init > 'PrometheusEndPoints' = " + strings.Join(PrometheusEndPoints, " / "))
}

// JSON Structs
type Metric struct {
	Name     string `json:"__name__"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

type Result struct {
	Metrics Metric        `json:"metric"`
	Values  []interface{} `json:"value"`
}

type Data struct {
	Resulttype string   `json:"resultType"`
	Results    []Result `json:"result"`
}

type Query struct {
	Status     string `json:"status"`
	Data_Array Data   `json:"data"`
}

// MonitoringAdapter struct
type monitoringAdapter struct {
	agreement *model.Agreement
	size      int
}

// New returns a new Prometheus Monitoring Adapter.
func New(size int) monitor.MonitoringAdapter {
	log.Println("SLALite > prometheusadapter [New] Creating a new Prometheus adapter ...")

	return &monitoringAdapter{
		agreement: nil,
		size:      size,
	}
}

// Initialize function: called before GetValues
func (ma *monitoringAdapter) Initialize(a *model.Agreement) {
	log.Println("SLALite > prometheusadapter [Initialize] Call to Initialize method ...")

	ma.agreement = a
}

// GetValues function: collect information from prometheus and send violations
//
// returns assessment_model.GuaranteeData:
//
// 		GuaranteeData represents the list of values needed to evaluate an expression at several points in time
//		==> type GuaranteeData []ExpressionData
//
// 		A ExpressionData represents the set of values needed to evaluate an expression at a single time
// 		==> type ExpressionData map[string]model.MetricValue
//
// 		A MetricValue is the SLALite representation of a metric value.
// 		==>	type MetricValue struct {
// 				Key      string      `json:"key"`
// 				Value    interface{} `json:"value"`
// 				DateTime time.Time   `json:"datetime"`
// 			}
//
// 	==> GuaranteeData[0]["algo"].Key
//  ==> GuaranteeData[0]["algo"].Value
func (ma *monitoringAdapter) GetValues(gt model.Guarantee, vars []string) assessment_model.GuaranteeData {
	log.Println("SLALite > prometheusadapter > [GetValues] Getting values from Prometheus for agreement [" + ma.agreement.Id + "] ...")

	if len(PrometheusEndPoints) > 1 {
		log.Println("SLALite > prometheusadapter [GetValues] getting values from multiple sources ...")
		return getValuesFromMultipleSources(ma, gt, vars)
	} else {
		//log.Println("SLALite > prometheusadapter [GetValues] Checking metrics ...")
		result := make(assessment_model.GuaranteeData, 0)

		// log.Println("SLALite > prometheusadapter [GetValues] Checking metrics ...")
		for _, metricName := range vars {
			// log.Println("SLALite > prometheusadapter [GetValues] Checking metric: " + metricName)

			for _, pm := range MetricsPrometheus {
				if pm == metricName {
					//log.Println("SLALite > prometheusadapter [GetValues] Getting metric '" + pm + "' values from PROMETHEUS ...")
					Up := ""
					var TimeFloat float64
					var TimeInt int

					i := strings.Index(URLPrometheus, "https")
					log.Println("SLALite > prometheusadapter [GetValues] HTTPS? ", i)
					var stdout, stderr bytes.Buffer
					err := errors.New("no error found")
					if i > -1 {
						c := exec.Command("curl", "-g", "-k", URLPrometheus+"/api/v1/query?query="+pm)
						c.Stdout = &stdout
						c.Stderr = &stderr
						err = c.Run()
					} else {
						c := exec.Command("curl", "-g", URLPrometheus+"/api/v1/query?query="+pm)
						c.Stdout = &stdout
						c.Stderr = &stderr
						err = c.Run()
					}

					if err != nil {
						log.Println("SLALite > prometheusadapter [GetValues] Metric: '"+pm+"', ERROR executing query [GET "+URLPrometheus+"/api/v1/query?query="+pm+"]", err)
					} else {
						str := string(stdout.Bytes())
						log.Println("SLALite > prometheusadapter [GetValues] Metric: '" + pm + "', Results: " + str)

						res := Query{}
						json.Unmarshal([]byte(str), &res)

						for i := range res.Data_Array.Results {
							job := res.Data_Array.Results[i].Metrics.Job

							if job == "kubernetes-nodes-cadvisor" || job == "compss" || job == "deadlines_missed" || strings.Index(job, "deadlines_missed") != -1 {
								// Metric and Time:
								Up = res.Data_Array.Results[i].Metrics.Name
								TimeFloat = res.Data_Array.Results[i].Values[0].(float64)

								// Value: "kubernetes-nodes-cadvisor": get stacked values
								strVal := res.Data_Array.Results[i].Values[1].(string)

								n, err := strconv.ParseInt(strVal, 10, 64)
								if err == nil {
									if Up == "go_memstats_frees_total" {
										TimeInt = TimeInt + int(n/1024/1024)
									} else {
										// "kubelet_running_pod_count"
										TimeInt = TimeInt + int(n)
									}
								} else {
									log.Println("SLALite > prometheusadapter [GetValues] ERROR: Returning [TimeInt = int(TimeFloat)] ... ", err)
									TimeInt = int(TimeFloat)
								}
							}
						}

						//log.Printf("SLALite > prometheusadapter [GetValues] Result: Key: %s, Time: %f, Value: %d", Up, TimeFloat, TimeInt)
						for i := 0; i < ma.size; i++ {
							val := make(assessment_model.ExpressionData)

							val[metricName] = model.MetricValue{
								DateTime: time.Now(),
								Key:      metricName,
								Value:    TimeInt,
							}

							result = append(result, val)
						}
					}
				}
			}
		}

		//log.Println("SLALite > prometheusadapter [GetValues] Returning results [size=" + strconv.Itoa(len(result)) + "] ...")
		return result
	}
}

/*
GetValuesFromMultipleSources
*/
func getValuesFromMultipleSources(ma *monitoringAdapter, gt model.Guarantee, vars []string) assessment_model.GuaranteeData {
	result := make(assessment_model.GuaranteeData, 0)

	for _, endpoint := range PrometheusEndPoints {
		// log.Println("SLALite > prometheusadapter [GetValues] Checking metrics ...")
		for _, metricName := range vars {
			// log.Println("SLALite > prometheusadapter [GetValues] Checking metric: " + metricName)

			for _, pm := range MetricsPrometheus {
				if pm == metricName {
					//log.Println("SLALite > prometheusadapter [GetValues] Getting metric '" + pm + "' values from PROMETHEUS ...")
					Up := ""
					var TimeFloat float64
					var TimeInt int
					c := exec.Command("curl", "-g", endpoint+"/api/v1/query?query="+pm)
					var stdout, stderr bytes.Buffer

					c.Stdout = &stdout
					c.Stderr = &stderr
					err := c.Run()

					if err != nil {
						log.Println("SLALite > prometheusadapter [GetValues] Metric: '"+pm+"', ERROR executing query [GET "+URLPrometheus+"/api/v1/query?query="+pm+"]", err)
					} else {
						str := string(stdout.Bytes())
						log.Println("SLALite > prometheusadapter [GetValues] Metric: '" + pm + "', Results: " + str)

						res := Query{}
						json.Unmarshal([]byte(str), &res)

						for i := range res.Data_Array.Results {
							job := res.Data_Array.Results[i].Metrics.Job

							if job == "kubernetes-nodes-cadvisor" || job == "compss" || job == "deadlines_missed" || strings.Index(job, "deadlines_missed") != -1 {
								// Metric and Time:
								Up = res.Data_Array.Results[i].Metrics.Name
								TimeFloat = res.Data_Array.Results[i].Values[0].(float64)

								// Value: "kubernetes-nodes-cadvisor": get stacked values
								strVal := res.Data_Array.Results[i].Values[1].(string)

								n, err := strconv.ParseInt(strVal, 10, 64)
								if err == nil {
									if Up == "go_memstats_frees_total" {
										TimeInt = TimeInt + int(n/1024/1024)
									} else {
										// "kubelet_running_pod_count"
										TimeInt = TimeInt + int(n)
									}
								} else {
									log.Println("SLALite > prometheusadapter [GetValues] ERROR: Returning [TimeInt = int(TimeFloat)] ... ", err)
									TimeInt = int(TimeFloat)
								}
							}
						}

						//log.Printf("SLALite > prometheusadapter [GetValues] Result: Key: %s, Time: %f, Value: %d", Up, TimeFloat, TimeInt)
						for i := 0; i < ma.size; i++ {
							val := make(assessment_model.ExpressionData)

							val[metricName] = model.MetricValue{
								DateTime: time.Now(),
								Key:      metricName,
								Value:    TimeInt,
							}

							result = append(result, val)
						}
					}
				}
			}
		}
	}

	return result
}
