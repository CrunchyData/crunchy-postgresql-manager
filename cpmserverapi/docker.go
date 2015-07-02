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
	dockerapi "github.com/fsouza/go-dockerclient"
	"net/http"
	"os/exec"
)

type DockerInspectRequest struct {
	ContainerName string
}
type DockerInspectResponse struct {
	IPAddress    string
	RunningState string
}
type DockerRemoveRequest struct {
	ContainerName string
}
type DockerRemoveResponse struct {
	Output string
}
type DockerStartRequest struct {
	ContainerName string
}
type DockerStartResponse struct {
	Output string
}
type DockerStopRequest struct {
	ContainerName string
}
type DockerStopResponse struct {
	Output string
}

type DockerRunRequest struct {
	CPU           string
	MEM           string
	ClusterID     string
	ServerID      string
	ProjectID     string
	Image         string
	IPAddress     string
	Standalone    string
	PGDataPath    string
	ContainerName string
	ContainerType string
	CommandOutput string
	CommandPath   string
	EnvVars       map[string]string
}
type DockerRunResponse struct {
	Output string
}

func DockerInspect(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DockerInspect called")
	req := DockerInspectRequest{}
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

	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		logit.Error.Println("can't get connection to docker socket")
		rest.Error(w, err.Error(), 400)
		return
	}

	//if we can't inspect the container, then we give up
	//on trying to stop it
	response := DockerInspectResponse{}
	response.RunningState = "down"
	response.IPAddress = ""

	container, err3 := docker.InspectContainer(req.ContainerName)
	if err3 != nil {
		logit.Info.Println("can't inspect container " + req.ContainerName)
		w.WriteJson(&response)
		return
	}

	if container != nil {
		logit.Info.Println("container found during inspect")
		if container.State.Running {
			response.RunningState = "up"
			logit.Info.Println("container status is up")
			logit.Info.Println("container ipaddress is " + container.NetworkSettings.IPAddress)
			response.IPAddress = container.NetworkSettings.IPAddress
		} else {
			response.RunningState = "down"
			logit.Info.Println("container status is down")
		}
	}

	w.WriteJson(&response)
}

func DockerRemove(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DockerRemove called")
	req := DockerRemoveRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//if a container exists with that name, then we need
	//to stop it first and then remove it
	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		logit.Error.Println("can't get connection to docker socket")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response DockerRemoveResponse

	//if we can't inspect the container, then we give up
	//on trying to remove it, this is ok to pass since
	//a user can remove the container manually
	container, err3 := docker.InspectContainer(req.ContainerName)
	if err3 != nil {
		logit.Error.Println("during remove....can't inspect container " + req.ContainerName)
		logit.Error.Println("inspect container error was " + err3.Error())
		response.Output = "success"
		w.WriteJson(&response)
		return
	}

	if container != nil {
		logit.Info.Println("during remove...container found")
		err3 = docker.StopContainer(req.ContainerName, 10)
		if err3 != nil {
			logit.Error.Println("can't stop container " + req.ContainerName)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logit.Info.Println("during remove....container stopped ")
		opts := dockerapi.RemoveContainerOptions{ID: req.ContainerName}
		err := docker.RemoveContainer(opts)
		if err != nil {
			logit.Error.Println("can't remove container " + req.ContainerName)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logit.Info.Println("container " + req.ContainerName + " removed ")
	}

	response.Output = "success"

	w.WriteJson(&response)
}

func DockerStart(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DockerStart called")
	req := DockerStartRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docker, err3 := dockerapi.NewClient("unix://var/run/docker.sock")
	if err3 != nil {
		logit.Error.Println("can't get connection to docker socket")
		rest.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	err3 = docker.StartContainer(req.ContainerName, nil)
	if err3 != nil {
		logit.Error.Println("can't start container " + req.ContainerName)
		rest.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	var response DockerStartResponse
	response.Output = "success"
	w.WriteJson(&response)
}

func DockerStop(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DockerStop called")
	req := DockerStopRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	docker, err3 := dockerapi.NewClient("unix://var/run/docker.sock")
	if err3 != nil {
		logit.Error.Println("can't get connection to docker socket")
		rest.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	err3 = docker.StopContainer(req.ContainerName, 10)
	if err3 != nil {
		logit.Error.Println("can't stop container " + req.ContainerName)
		rest.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	var response DockerStopResponse
	response.Output = "success"
	w.WriteJson(&response)
}

func DockerRun(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("DockerRun called")
	req := DockerRunRequest{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allEnvVars = ""
	if req.EnvVars != nil {
		for k, v := range req.EnvVars {
			allEnvVars = allEnvVars + " -e " + k + "=" + v
		}
	}
	logit.Info.Println("env vars " + allEnvVars)

	cmd := exec.Command(req.CommandPath, req.PGDataPath, req.ContainerName,
		req.Image, req.CPU, req.MEM, allEnvVars)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		logit.Error.Println(stderr.String())
		logit.Error.Println(out.String())
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	response := DockerRunResponse{}
	response.Output = out.String()
	w.WriteJson(&response)
}
