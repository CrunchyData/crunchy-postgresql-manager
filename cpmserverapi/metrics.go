/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cpmserverapi

import (
	"bytes"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"net/http"
	"os/exec"
)

type MetricCPURequest struct {
	Something string
	Other     string
}
type MetricCPUResponse struct {
	Output string
}
type MetricMEMRequest struct {
	Something string
	Other     string
}
type MetricMEMResponse struct {
	Output string
}

type MetricIostatRequest struct {
	Something string
	Other     string
}
type MetricIostatResponse struct {
	Output string
}
type MetricDfRequest struct {
	Something string
	Other     string
}
type MetricDfResponse struct {
	Output string
}

func MetricCPU(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("MetricCPU called")
	req := MetricCPURequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("monitor-load")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
	}

	var response MetricCPUResponse
	response.Output = out.String()
	w.WriteJson(&response)
}

func MetricMEM(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("MetricMEM called")
	req := MetricMEMRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("monitor-mem")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
	}

	var response MetricMEMResponse
	response.Output = out.String()
	w.WriteJson(&response)
}

func MetricIostat(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("MetricIostat called")
	post := MetricIostatRequest{}
	err := r.DecodeJsonPayload(&post)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("cpmiostat")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
	}

	var response MetricIostatResponse
	response.Output = out.String()
	w.WriteJson(&response)
}

func MetricDf(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("MetricDf called")
	post := MetricDfRequest{}
	err := r.DecodeJsonPayload(&post)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("cpmdf")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
	}
	var response MetricDfResponse
	response.Output = out.String()
	w.WriteJson(&response)
}
