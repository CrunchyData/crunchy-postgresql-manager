//
// example of calling one of these
//
// request := &cpmserverapi.MetricIostatRequest{"something", "yes"}
// response, err := cpmserverapi.MetricIostatClient("http://localhost:10001", request)
//
package cpmserverapi

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
)

const URL = "http://cpm-"
const PORT = ":10001"

// MetricIostatClient perform an iostat on the host and return the results of it
func MetricIostatClient(serverName string, req *MetricIostatRequest) (MetricIostatResponse, error) {
	var err error
	var url = URL + serverName + PORT

	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/iostat", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricIostatResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// MetricDfClient perform a df command on the host and return the results of it
func MetricDfClient(serverName string, req *MetricDfRequest) (MetricDfResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/df", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricDfResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DiskProvisionClient client for provisioning a new directory on the host
func DiskProvisionClient(serverName string, req *DiskProvisionRequest) (DiskProvisionResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/disk/provision", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DiskProvisionResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DiskDeleteClient client for deleting a directory on the  host
func DiskDeleteClient(serverName string, req *DiskDeleteRequest) (DiskDeleteResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/disk/delete", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DiskDeleteResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

/**
// DockerInspectClient client for performing a Docker inspect command
func DockerInspectClient(serverName string, req *DockerInspectRequest) (DockerInspectResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/inspect", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerInspectResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DockerRemoveClient client for performing a Docker remove command
func DockerRemoveClient(serverName string, req *DockerRemoveRequest) (DockerRemoveResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/remove", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerRemoveResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DockerStartClient client for performing a Docker start command
func DockerStartClient(serverName string, req *DockerStartRequest) (DockerStartResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/start", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerStartResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DockerStopClient client for performing a Docker stop command
func DockerStopClient(serverName string, req *DockerStopRequest) (DockerStopResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/stop", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerStopResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DockerRunClient client for performing a Docker run command
// example:
//
//request := &cpmserverapi.DockerRunRequest{}
//request.CommandPath = "docker-run.sh"
//request.CPU = "0"
//request.MEM = "0"
//request.Image = "cpm-node"
//request.PGDataPath = "/tmp/foo"
//request.ContainerName = "testpoo"
//envvars := make(map[string]string)
//envvars["one"] = "uno"
//request.EnvVars = envvars
//
func DockerRunClient(serverName string, req *DockerRunRequest) (DockerRunResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/run", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerRunResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
*/

// MetricCPUClient client for getting the cpu metrics
func MetricCPUClient(serverName string, req *MetricCPURequest) (MetricCPUResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/cpu", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricCPUResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// MetricMEMClient client for getting the memory metrics
func MetricMEMClient(serverName string, req *MetricMEMRequest) (MetricMEMResponse, error) {
	var err error
	var url = URL + serverName + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/mem", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricMEMResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
