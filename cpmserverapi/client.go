package cpmserverapi

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
)

//
// example of calling one of these
//
// request := &cpmserverapi.MetricIostatRequest{"something", "yes"}
// response, err := cpmserverapi.MetricIostatClient("http://localhost:10001", request)
//
func MetricIostatClient(url string, req *MetricIostatRequest) (MetricIostatResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/iostat", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricIostatResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func MetricDfClient(url string, req *MetricDfRequest) (MetricDfResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/df", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricDfResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DiskProvisionClient(url string, req *DiskProvisionRequest) (DiskProvisionResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/disk/provision", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DiskProvisionResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DiskDeleteClient(url string, req *DiskDeleteRequest) (DiskDeleteResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/disk/delete", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DiskDeleteResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DockerInspectClient(url string, req *DockerInspectRequest) (DockerInspectResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/inspect", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerInspectResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DockerRemoveClient(url string, req *DockerRemoveRequest) (DockerRemoveResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/remove", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerRemoveResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DockerStartClient(url string, req *DockerStartRequest) (DockerStartResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/start", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerStartResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func DockerStopClient(url string, req *DockerStopRequest) (DockerStopResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/stop", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerStopResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

//
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
func DockerRunClient(url string, req *DockerRunRequest) (DockerRunResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/docker/run", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := DockerRunResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func MetricCPUClient(url string, req *MetricCPURequest) (MetricCPUResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/cpu", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricCPUResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func MetricMEMClient(url string, req *MetricMEMRequest) (MetricMEMResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url+"/metrics/mem", "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := MetricMEMResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
