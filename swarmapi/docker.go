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

package swarmapi

import (
	"errors"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/samalba/dockerclient"
	"os"
	"strings"
)

type DockerInspectRequest struct {
	ContainerName string
}
type DockerInspectResponse struct {
	IPAddress    string
	RunningState string
	ServerID     string
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

type DockerInfoResponse struct {
	Output []string
}

type DockerPsInfo struct {
	Name   string
	Status string
	Image  string
}
type DockerPsResponse struct {
	Output []DockerPsInfo
}

type DockerRunRequest struct {
	CPU               string
	MEM               string
	ClusterID         string
	ServerID          string
	ProjectID         string
	Image             string
	IPAddress         string
	Standalone        string
	PGDataPath        string
	ContainerName     string
	ContainerType     string
	CommandOutput     string
	CommandPath       string
	Profile           string
	Token             string
	EnvVars           map[string]string
	RestoreJob        string
	RestoreRemotePath string
	RestoreRemoteHost string
	RestoreRemoteUser string
	RestoreDbUser     string
	RestoreDbPass     string
	RestoreSet        string
}
type DockerRunResponse struct {
	ID string
}

// DockerInspect perform a Docker inspect
func DockerInspect(req *DockerInspectRequest) (DockerInspectResponse, error) {
	response := DockerInspectResponse{}
	var err error
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	//logit.Info.Println("DockerInspect called")

	if req.ContainerName == "" {
		err = errors.New("ContainerName required in request")
		logit.Error.Println(err.Error())
		return response, err
	}

	docker, err := dockerapi.NewClient(swarmURL)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	//if we can't inspect the container, then we give up
	//on trying to stop it
	response.RunningState = "not-found"
	response.IPAddress = ""

	container, err3 := docker.InspectContainer(req.ContainerName)
	if err3 != nil {
		logit.Info.Println("can't inspect container " + req.ContainerName)
		return response, nil
	}

	if container != nil {
		//logit.Info.Println("container found during inspect")
		if container.State.Running {
			response.RunningState = "up"
			//logit.Info.Println("container status is up")
			//logit.Info.Println("container ipaddress is " + container.NetworkSettings.IPAddress)
			response.IPAddress = container.NetworkSettings.IPAddress
		} else {
			response.RunningState = "down"
			//logit.Info.Println("container status is down")
		}
	}

	response.ServerID = container.Node.Addr

	return response, nil
}

// DockerRemove perform a Docker remove
func DockerRemove(req *DockerRemoveRequest) (DockerRemoveResponse, error) {
	response := DockerRemoveResponse{}
	var err error
	logit.Info.Println("DockerRemove called")
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	//if a container exists with that name, then we need
	//to stop it first and then remove it
	docker, err := dockerapi.NewClient(swarmURL)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	//if we can't inspect the container, then we give up
	//on trying to remove it, this is ok to pass since
	//a user can remove the container manually
	container, err3 := docker.InspectContainer(req.ContainerName)
	if err3 != nil {
		logit.Error.Println("during remove....can't inspect container " + req.ContainerName)
		logit.Error.Println("inspect container error was " + err3.Error())
		response.Output = "success"
		return response, nil
	}

	if container != nil {
		//logit.Info.Println("during remove...container found")
		err3 = docker.StopContainer(req.ContainerName, 10)
		if err3 != nil {
			logit.Error.Println("can't stop container " + req.ContainerName)
			logit.Error.Println(err3.Error())
		}
		//logit.Info.Println("during remove....container stopped ")
		opts := dockerapi.RemoveContainerOptions{ID: req.ContainerName}
		err := docker.RemoveContainer(opts)
		if err != nil {
			logit.Error.Println("can't remove container " + req.ContainerName)
			return response, err
		}
		logit.Info.Println("container " + req.ContainerName + " removed ")
	}

	response.Output = "success"

	return response, nil
}

// DockerStart perform a docker start
func DockerStart(req *DockerStartRequest) (DockerStartResponse, error) {
	var response DockerStartResponse

	//logit.Info.Println("DockerStart called")
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	docker, err := dockerapi.NewClient(swarmURL)
	if err != nil {
		logit.Error.Println("can't get connection to docker socket")
		return response, err
	}
	err = docker.StartContainer(req.ContainerName, nil)
	if err != nil {
		logit.Error.Println("can't start container " + req.ContainerName)
		return response, err
	}

	response.Output = "success"
	return response, nil
}

// DockerStop perform a docker stop
func DockerStop(req *DockerStopRequest) (DockerStopResponse, error) {
	var response DockerStopResponse
	//logit.Info.Println("DockerStop called")
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	docker, err := dockerapi.NewClient(swarmURL)
	if err != nil {
		logit.Error.Println("can't get connection to docker socket")
		return response, err
	}
	err = docker.StopContainer(req.ContainerName, 10)
	if err != nil {
		logit.Error.Println("can't stop container " + req.ContainerName)
		return response, err
	}

	response.Output = "success"
	return response, nil
}

// DockerRun perform a docker run
func DockerRun(req *DockerRunRequest) (DockerRunResponse, error) {
	response := DockerRunResponse{}
	//logit.Info.Println("DockerRun called")
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	var envvars []string
	var i = 0
	if req.EnvVars != nil {
		envvars = make([]string, len(req.EnvVars)+1)
		for k, v := range req.EnvVars {
			envvars[i] = k + "=" + v
			i++
		}
	} else {
		envvars = make([]string, 1)
	}

	if req.Profile == "" {
		return response, errors.New("Profile was empty and should not be")
	}

	//typical case is to always add the profile constraint env var
	//like SM, MED, LG, however in the case of a restore job, we
	//use a hard constraint of the host ipaddress to pin
	//the restored container to the same host as where the backup
	//is stored
	if req.IPAddress != "" {
		envvars[i] = "constraint:host==~" + req.IPAddress
	} else {
		envvars[i] = "constraint:profile==~" + req.Profile
	}

	docker, err := dockerapi.NewClient(swarmURL)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	options := dockerapi.CreateContainerOptions{}
	config := dockerapi.Config{}
	config.Hostname = req.ContainerName
	options.Config = &config
	hostConfig := dockerapi.HostConfig{}
	options.HostConfig = &hostConfig
	options.Name = req.ContainerName
	options.Config.Env = envvars
	options.Config.Image = "crunchydata/" + req.Image
	//logit.Info.Println("swarmapi using " + options.Config.Image + " as the image name")
	options.Config.Volumes = make(map[string]struct{})

	//TODO figure out cpu shares and memory settings, these are different
	//than what I was using before due to me using the docker api directly
	//with this swarm implementation...use the defaults for now

	//options.HostConfig.CPUShares, err = strconv.ParseInt(req.CPU, 0, 64)
	//if err != nil {
	//logit.Error.Println(err.Error())
	//return response, err
	//}
	//options.HostConfig.Memory = req.MEM

	options.HostConfig.Binds = make([]string, 3)
	options.HostConfig.Binds[0] = req.PGDataPath + ":/pgdata"
	options.HostConfig.Binds[1] = "/var/cpm/data/keys:/keys"
	options.HostConfig.Binds[2] = "/var/cpm/config:/syslogconfig"

	container, err3 := docker.CreateContainer(options)
	if err3 != nil {
		logit.Error.Println(err3.Error())
		return response, err3
	}

	var startResponse DockerStartResponse
	startRequest := DockerStartRequest{}
	startRequest.ContainerName = req.ContainerName
	startResponse, err = DockerStart(&startRequest)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	logit.Info.Println(startResponse.Output)
	//cmd := exec.Command(req.CommandPath, req.PGDataPath, req.ContainerName,
	//req.Image, req.CPU, req.MEM, allEnvVars)

	response.ID = container.ID
	return response, nil
}

func DockerInfo() (DockerInfoResponse, error) {
	response := DockerInfoResponse{}
	var err error
	swarmURL := os.Getenv("SWARM_MANAGER_URL")
	if swarmURL == "" {
		logit.Error.Println("SWARM_MANAGER_URL not set")
		return response, errors.New("SWARM_MANAGER_URL not set")
	}

	//logit.Info.Println("DockerInfo called")

	docker, err := dockerclient.NewDockerClient(swarmURL, nil)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	var info *dockerclient.Info
	info, err = docker.Info()
	if err != nil {
		logit.Error.Println(err.Error())
		return response, nil
	}

	var colonLoc int
	response.Output = make([]string, 0)
	for x := 0; x < len(info.DriverStatus); {
		//fmt.Printf("index %d  [%s]\n", x, info.DriverStatus[x])
		trimmedStr := strings.TrimSpace(info.DriverStatus[x][1])
		//fmt.Printf("trimmed [%s]\n", strings.TrimSpace(info.DriverStatus[x][1]))

		colonLoc = strings.Index(trimmedStr, ":")
		//fmt.Printf("colonLoc=%d\n", colonLoc)
		if colonLoc > 0 {
			//logit.Info.Println("found " + trimmedStr)
			//parts := strings.Split(trimmedStr, ":")
			//response.Output = append(response.Output, parts[0])
			response.Output = append(response.Output, trimmedStr)
		}
		x++
	}

	return response, nil
}

func DockerPs(serverid string) (DockerPsResponse, error) {
	response := DockerPsResponse{}
	var err error

	//logit.Info.Println("DockerPs called on " + serverid)

	docker, err := dockerapi.NewClient("tcp://" + serverid)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	var info []dockerapi.APIContainers
	options := dockerapi.ListContainersOptions{}
	options.All = true

	info, err = docker.ListContainers(options)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	response.Output = make([]DockerPsInfo, 0)

	var apicontainer dockerapi.APIContainers
	for x := range info {
		apicontainer = info[x]
		//logit.Info.Printf("x=%d\n", x)
		if strings.Index(apicontainer.Image, "cpm-node") > 0 ||
			strings.Index(apicontainer.Image, "cpm-pgpool") > 0 {
			cinfo := DockerPsInfo{}
			if len(apicontainer.Names) > 0 {
				cinfo.Name = strings.Trim(apicontainer.Names[0], "/")
				//logit.Info.Println("name=" + cinfo.Name + " status=" + apicontainer.Status)
			}
			cinfo.Status = apicontainer.Status
			cinfo.Image = apicontainer.Image
			response.Output = append(response.Output, cinfo)
		}
	}

	//logit.Info.Println("dockerps returning results")
	return response, nil
}
