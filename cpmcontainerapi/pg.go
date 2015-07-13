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
	"io/ioutil"
	"net/http"
	"os/exec"
)

type RemoteWritefileRequest struct {
	Path         string
	Filecontents string
}
type RemoteWritefileResponse struct {
	Status string
}

type InitdbRequest struct {
	ContainerName string
}
type InitdbResponse struct {
	Output string
	Status string
}
type StartPGRequest struct {
	ContainerName string
}
type StartPGResponse struct {
	Output string
	Status string
}
type StartPGOnStandbyRequest struct {
	ContainerName string
}
type StartPGOnStandbyResponse struct {
	Output string
	Status string
}
type StopPGRequest struct {
	ContainerName string
}
type StopPGResponse struct {
	Output string
	Status string
}
type BasebackupRequest struct {
	MasterHostName  string
	StandbyHostName string
}
type BasebackupResponse struct {
	Output string
	Status string
}
type FailoverRequest struct {
	ContainerName string
}
type FailoverResponse struct {
	Output string
	Status string
}
type SeedRequest struct {
	ContainerName string
}
type SeedResponse struct {
	Output string
	Status string
}
type ControldataRequest struct {
	Path string
}
type ControldataResponse struct {
	Output string
	Status string
}

func RemoteWritefile(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("RemoteWritefile called")
	req := RemoteWritefileRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Path == "" {
		rest.Error(w, "Path not supplied in request", http.StatusInternalServerError)
		return
	}

	if req.Filecontents == "" {
		rest.Error(w, "Filecontents not supplied in request", http.StatusInternalServerError)
		return
	}

	d1 := []byte(req.Filecontents)
	err = ioutil.WriteFile(req.Path, d1, 0644)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response RemoteWritefileResponse
	response.Status = "OK"
	w.WriteJson(&response)
}

func Initdb(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("Initdb called")
	req := InitdbRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("initdb.sh", req.ContainerName)
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

	var response InitdbResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func StartPG(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("StartPG called")
	req := StartPGRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("startpg.sh", req.ContainerName)
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

	var response StartPGResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func StopPG(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("StopPG called")
	req := StopPGRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("stoppg.sh", req.ContainerName)
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

	var response StopPGResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func Basebackup(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("Basebackup called")
	req := BasebackupRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.MasterHostName == "" {
		rest.Error(w, "MasterHostName not supplied in request", http.StatusInternalServerError)
		return
	}
	if req.StandbyHostName == "" {
		rest.Error(w, "StandbyHostName not supplied in request", http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("basebackup.sh", req.MasterHostName, req.StandbyHostName)
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

	var response BasebackupResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func Failover(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("Failover called")
	req := FailoverRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.ContainerName == "" {
		logit.Error.Println("ContainerName not supplied in request")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cmd *exec.Cmd
	cmd = exec.Command("fail-over.sh", req.ContainerName)
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

	var response FailoverResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func Seed(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("Seed called")
	req := SeedRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("seed.sh", req.ContainerName)
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

	var response SeedResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func StartPGOnStandby(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("StartPGOnStandby called")
	req := StartPGOnStandbyRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("startpgonstandby.sh", req.ContainerName)
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

	var response StartPGOnStandbyResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}

func Controldata(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("Controldata called")
	req := ControldataRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("pg_controldata", "/pgdata")
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

	var response ControldataResponse
	response.Output = out.String()
	response.Status = "OK"
	w.WriteJson(&response)
}
