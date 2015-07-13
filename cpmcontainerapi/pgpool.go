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

package cpmcontainerapi

import (
	"bytes"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"net/http"
	"os/exec"
)

type StartPgpoolRequest struct {
	ContainerName string
	Path          string
}
type StartPgpoolResponse struct {
	Output string
	Status string
}
type StopPgpoolRequest struct {
	ContainerName string
	Path          string
}
type StopPgpoolResponse struct {
	Output string
	Status string
}

func StartPgpool(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("StartPgpool called")
	response := StartPgpoolResponse{}
	req := StartPgpoolRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.ContainerName == "" {
		logit.Error.Println("ContainerName required")
		rest.Error(w, err.Error(), 400)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("startpgpool.sh", req.Path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	response.Output = out.String()
	response.Status = "OK"

	w.WriteJson(&response)
}

func StopPgpool(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("StopPgpool called")
	response := StopPgpoolResponse{}
	req := StopPgpoolRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("stop-pgpool.sh", req.Path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	response.Output = out.String()
	response.Status = "OK"

	w.WriteJson(&response)
}
