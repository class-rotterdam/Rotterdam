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
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// httpRequest prepares and executes the HTTP request //// bodyJSON interface{}
func httpRequest(httpMethod string, url string, auth bool, body io.Reader) (int, []byte, error) {
	log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] " + httpMethod + " request [" + url + "], auth [" + strconv.FormatBool(auth) + "] ...")

	// create request with headers and body
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] ERROR (2)", err)
		return 0, nil, err
	}

	// Content-Type: json / json-patch+json
	if httpMethod == "PATCH" {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] Content-Type = application/json-patch+json")
		req.Header.Set("Content-Type", "application/json-patch+json")
	} else {
		log.Println("Rotterdam > Adaptation-Engine > monitoring http [httpRequest] Content-Type = application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute HTTP request
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] ERROR (4)", err)
		return resp.StatusCode, nil, err
	} else if resp.StatusCode >= 400 { // check errors => StatusCode
		log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] ERROR (5) StatusCode >= 400")
		return resp.StatusCode, nil, errors.New("Rotterdam > CAAS > http [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode) + "")
	}

	log.Println("Rotterdam > Adaptation-Engine > monitoring [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	//log.Println("Rotterdam > CAAS > http [HTTPRequest] RESPONSE: " + string(data))

	return resp.StatusCode, data, nil
}

///////////////////////////////////////////////////////////////////////////////
// GET

/*
HTTPGET generic GET request
*/
func HTTPGET(url string, auth bool) (int, []byte, error) {
	return httpRequest("GET", url, auth, nil)
}

/*
GETHttpStruct GET request that returns a struct of type 'map[string]interface{}'
*/
func GETHttpStruct(url string, auth bool) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > Adaptation-Engine > monitoring [HTTPGETStruct] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [HTTPGETStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [HTTPGETStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
GETHttpString GET request that returns a string (response)
*/
func GETHttpString(url string, auth bool) (int, string, error) {
	log.Println("Rotterdam > Adaptation-Engine > monitoring [HTTPGETString] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth)
	if err != nil {
		log.Println("Rotterdam > Adaptation-Engine > monitoring [HTTPGETString] ERROR (1)", err)
		return status, "", err
	}

	return status, string(data), nil
}
