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

package common

import (
	cfg "atos/rotterdam/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////
// GET

/*
 * HttpGET: generic GET request
 */
func HttpGET(url string) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [HttpGET] GET request [" + url + "] ...")

	// create request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET] ERROR (1)", err)
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute GET request
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET] ERROR (2)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET] ERROR (3)", err)
		return resp.StatusCode, nil, err
	}

	log.Println("Rotterdam > CAAS > http [HttpGET] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	log.Println("Rotterdam > CAAS > http [HttpGET] RESPONSE: " + string(data))

	return resp.StatusCode, data, nil
}

/*
 * HttpGET_GenericStruct: GET request that returns a struct of type 'map[string]interface{}'
 */
func HttpGET_GenericStruct(url string) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HttpGET_GenericStruct] GET request [" + url + "] ...")

	status, data, err := HttpGET(url)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET_GenericStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET_GenericStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
 * HttpGET_String: GET request that returns a string (response)
 */
func HttpGET_String(url string) (int, string, error) {
	log.Println("Rotterdam > CAAS > http [HttpGET_String] GET request [" + url + "] ...")

	status, data, err := HttpGET(url)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGET_String] ERROR (1)", err)
		return status, "", err
	}

	return status, string(data), nil
}

/*
 * HttpGETtest1: Test method
 */
func HttpGETtest1(url string) (string, error) {
	log.Println("Rotterdam > CAAS > http [HttpGETtest1] GET request ...")

	_, objmap, err := HttpGET_GenericStruct(url)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGETtest1] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > http [HttpGETtest1] url: " + objmap["url"].(string))
	log.Println("Rotterdam > CAAS > http [HttpGETtest1] args/foo1: " + objmap["args"].(map[string]interface{})["foo1"].(string))

	return objmap["url"].(string), nil
}

/*
 * HttpGETtest2: Test method
 */
func HttpGETtest2(url string) (string, error) {
	log.Println("Rotterdam > CAAS > http [HttpGETtest2] GET request ...")

	_, data, err := HttpGET_String(url)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpGETtest2] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > http [HttpGETtest2] url: " + data)

	return data, nil
}

///////////////////////////////////////////////////////////////////////////////
// POST

/*
 * HttpPOST: Generic POST request
 */
func HttpPOST(url string, body_json interface{}) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [HttpPOST] POST request [" + url + "] ...")

	// create request's body'
	bodyBytes, err := json.Marshal(body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST] ERROR (1)", err)
		return 0, nil, err
	}
	body := bytes.NewReader(bodyBytes)

	// create request with headers and body
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST] ERROR (2)", err)
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute POST request
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST] ERROR (4)", err)
		return resp.StatusCode, nil, err
	}

	log.Println("Rotterdam > CAAS > http [HttpPOST] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	log.Println("Rotterdam > CAAS > http [HttpPOST] RESPONSE: " + string(data))

	// return []byte
	return resp.StatusCode, data, nil
}

/*
 * HttpPOST_GenericStruct: POST request that returns a struct of type 'map[string]interface{}'
 */
func HttpPOST_GenericStruct(url string, body_json interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HttpPOST_GenericStruct] POST request [" + url + "] ...")

	status, data, err := HttpPOST(url, body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST_GenericStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOST_GenericStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
 * MOCKUP_HttpPOST_GenericStruct: POST request that returns a struct of type 'map[string]interface{}'
 */
func MOCKUP_HttpPOST_GenericStruct(url string, body_json interface{}) (map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [MOCKUP_HttpPOST_GenericStruct] POST request [" + url + "] ...")

	res := map[string]interface{}{
		"status": "200"}

	return res, nil
}

/*
 *
 */
func HttpPOSTTest1(url string) (string, error) {
	log.Println("Rotterdam > CAAS > http [HttpPOSTTest1] POST request ...")

	type Payload struct {
		Type     string      `json:"type"`
		Name     string      `json:"name"`
		Data     string      `json:"data"`
		Priority interface{} `json:"priority"`
		Port     interface{} `json:"port"`
		Weight   interface{} `json:"weight"`
	}

	data := Payload{
		Type:     "tipo1",
		Name:     "name1",
		Data:     "datadata",
		Priority: 1,
		Port:     8080,
		Weight:   300}

	_, resp_bytes, err := HttpPOST(url, data)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPOSTTest1] ERROR (1)", err)
		return "", err
	}

	return string(resp_bytes), nil
}

/*
 *
 */
func HttpPOSTTest2(url string) {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	type Payload struct {
		Type     string      `json:"type"`
		Name     string      `json:"name"`
		Data     string      `json:"data"`
		Priority interface{} `json:"priority"`
		Port     interface{} `json:"port"`
		Weight   interface{} `json:"weight"`
	}

	data := Payload{
		// fill struct
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.digitalocean.com/v2/domains/example.com/records", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer b7d03a6947b217efb6f3ec3bd3504582")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}

///////////////////////////////////////////////////////////////////////////////
// DELETE

/*
 * HttpDELETE: Generic DELETE request
 */
func HttpDELETE(url string, body_json interface{}) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [HttpDELETE] DELETE request [" + url + "] ...")

	// create request's body'
	bodyBytes, err := json.Marshal(body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE] ERROR (1)", err)
		return 0, nil, err
	}
	body := bytes.NewReader(bodyBytes)

	// create request with headers and body
	req, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE] ERROR (2)", err)
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute DELETE request
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE] ERROR (4)", err)
		return resp.StatusCode, nil, err
	}

	log.Println("Rotterdam > CAAS > http [HttpDELETE] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	log.Println("Rotterdam > CAAS > http [HttpDELETE] RESPONSE: " + string(data))

	// return []byte
	return resp.StatusCode, data, nil
}

/*
 * HttpDELETE_GenericStruct: DELETE request that returns a struct of type 'map[string]interface{}'
 */
func HttpDELETE_GenericStruct(url string) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HttpDELETE_GenericStruct] DELETE request [" + url + "] ...")

	type Body struct {
		Content interface{}
	}

	status, data, err := HttpDELETE(url, Body{})
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE_GenericStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HttpDELETE_GenericStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUT

/*
 * HttpPUT: Generic PUT request
 */
func HttpPUT(url string, body_json interface{}) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [HttpPUT] PUT request [" + url + "] ...")

	// create request's body'
	bodyBytes, err := json.Marshal(body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT] ERROR (1)", err)
		return 0, nil, err
	}
	body := bytes.NewReader(bodyBytes)

	// create request with headers and body
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT] ERROR (2)", err)
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute POST request
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT] ERROR (4)", err)
		return resp.StatusCode, nil, err
	}

	log.Println("Rotterdam > CAAS > http [HttpPUT] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	log.Println("Rotterdam > CAAS > http [HttpPUT] RESPONSE: " + string(data))

	// return []byte
	return resp.StatusCode, data, nil
}

/*
 * HttpPUT_GenericStruct: PUT request that returns a struct of type 'map[string]interface{}'
 */
func HttpPUT_GenericStruct(url string, body_json interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HttpPUT_GenericStruct] PUT request [" + url + "] ...")

	status, data, err := HttpPUT(url, body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT_GenericStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPUT_GenericStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PATCH

/*
 * HttpPATCH: Generic PATCH request
 */
func HttpPATCH(url string, body_json interface{}) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [HttpPATCH] PATCH request [" + url + "] ...")

	// create request's body'
	bodyBytes, err := json.Marshal(body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH] ERROR (1)", err)
		return 0, nil, err
	}
	body := bytes.NewReader(bodyBytes)

	// create request with headers and body
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH] ERROR (2)", err)
		return 0, nil, err
	}
	// -H "Content-Type:application/json-patch+json"
	req.Header.Set("Content-Type", "application/json-patch+json") // old value: application/json
	// add authorization header to the req
	req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute POST request
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH] ERROR (4)", err)
		return resp.StatusCode, nil, err
	}

	log.Println("Rotterdam > CAAS > http [HttpPATCH] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	log.Println("Rotterdam > CAAS > http [HttpPATCH] RESPONSE: " + string(data))

	// return []byte
	return resp.StatusCode, data, nil
}

/*
 * HttpPATCH_GenericStruct: PATCH request that returns a struct of type 'map[string]interface{}'
 */
func HttpPATCH_GenericStruct(url string, body_json interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HttpPATCH_GenericStruct] PATCH request [" + url + "] ...")

	status, data, err := HttpPATCH(url, body_json)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH_GenericStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HttpPATCH_GenericStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}
