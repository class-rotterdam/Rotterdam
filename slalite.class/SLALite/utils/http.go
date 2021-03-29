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

package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// httpRequest prepares and executes the HTTP request //// bodyJSON interface{}
func httpRequest(httpMethod string, url string, auth bool, authToken string, body io.Reader) (int, []byte, error) {
	log.Println("SLALite > Utils > HTTP [httpRequest] " + httpMethod + " request [" + url + "], auth [" + strconv.FormatBool(auth) + "] ...")

	// create request with headers and body
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		log.Println("SLALite > Utils > HTTP [httpRequest] ERROR (1)", err)
		return 0, nil, err
	}

	// Content-Type: json / json-patch+json
	if httpMethod == "PATCH" {
		log.Println("SLALite > Utils > HTTP [httpRequest] Content-Type = application/json-patch+json")
		req.Header.Set("Content-Type", "application/json-patch+json")
	} else {
		log.Println("SLALite > Utils > HTTP [httpRequest] Content-Type = application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// Authorization header
	if auth == true {
		log.Println("SLALite > Utils > HTTP [httpRequest] Using Authorization Bearer ...")
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("SLALite > Utils > HTTP [httpRequest] ERROR (2)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("SLALite > Utils > HTTP [httpRequest] ERROR (3)", err)
		return resp.StatusCode, nil, err
	} else if resp.StatusCode >= 400 { // check errors => StatusCode
		log.Println("SLALite > Utils > HTTP [httpRequest] ERROR (4) StatusCode >= 400")
		return resp.StatusCode, nil, errors.New("SLALite > Utils > HTTP [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode) + "")
	}

	log.Println("SLALite > Utils > HTTP [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))

	return resp.StatusCode, data, nil
}

///////////////////////////////////////////////////////////////////////////////
// GET

/*
HTTPGET generic GET request
*/
func HTTPGET(url string, auth bool, authToken string) (int, []byte, error) {
	return httpRequest("GET", url, auth, authToken, nil)
}

/*
HTTPGETStruct GET request that returns a struct of type 'map[string]interface{}'
*/
func HTTPGETStruct(url string, auth bool, authToken string) (int, map[string]interface{}, error) {
	log.Println("SLALite > Utils > HTTP [HTTPGETStruct] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth, authToken)
	if err != nil {
		log.Println("SLALite > Utils > HTTP [HTTPGETStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("SLALite > Utils > HTTP [HTTPGETStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
HTTPGETString GET request that returns a string (response)
*/
func HTTPGETString(url string, auth bool, authToken string) (int, string, error) {
	log.Println("SLALite > Utils > HTTP [HTTPGETString] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth, authToken)
	if err != nil {
		log.Println("SLALite > Utils > HTTP [HTTPGETString] ERROR (1)", err)
		return status, "", err
	}

	return status, string(data), nil
}
