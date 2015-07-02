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

type DiskProvisionRequest struct {
	Path string
}
type DiskProvisionResponse struct {
	Output string
	Status string
}

type DiskDeleteRequest struct {
	Path string
}
type DiskDeleteResponse struct {
	Output string
	Status string
}

func DiskProvision(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DiskProvision called")
	req := DiskProvisionRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("provisionvolume.sh", req.Path)
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

	var response DiskProvisionResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func DiskDelete(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DiskDelete called")
	req := DiskDeleteRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("deletevolume", req.Path)
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

	var response DiskDeleteResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}
