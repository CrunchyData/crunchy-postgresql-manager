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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"io/ioutil"
	"net/http"
	"strings"
)

const URL = "http://"
const PORT = ":10001"

// MetricIostatClient perform an iostat on the host and return the results of it
func MetricIostatClient(serverID string, req *MetricIostatRequest) (MetricIostatResponse, error) {
	var err error
	response := MetricIostatResponse{}
	serverParts := strings.Split(serverID, ":")
	var url = URL + serverParts[0] + PORT
	logit.Info.Println("MetricIostatClient url is " + url)

	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(url+"/metrics/iostat", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}

	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// MetricDfClient perform a df command on the host and return the results of it
func MetricDfClient(serverID string, req *MetricDfRequest) (MetricDfResponse, error) {
	var err error
	response := MetricDfResponse{}
	serverParts := strings.Split(serverID, ":")
	var url = URL + serverParts[0] + PORT
	logit.Info.Println("MetricDfClient url is " + url)
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(url+"/metrics/df", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DiskProvisionClient client for provisioning a new directory on the host
func DiskProvisionClient(serverName string, req *DiskProvisionRequest) (DiskProvisionResponse, error) {
	var err error
	response := DiskProvisionResponse{}
	serverParts := strings.Split(serverName, ":")
	var url = URL + serverParts[0] + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(url+"/disk/provision", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// DiskDeleteClient client for deleting a directory on the  host
func DiskDeleteClient(serverName string, req *DiskDeleteRequest) (DiskDeleteResponse, error) {
	response := DiskDeleteResponse{}
	serverParts := strings.Split(serverName, ":")
	var url = URL + serverParts[0] + PORT
	logit.Info.Println("deleting disk client with url=" + url)
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(url+"/disk/delete", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// MetricCPUClient client for getting the cpu metrics
func MetricCPUClient(serverID string, req *MetricCPURequest) (MetricCPUResponse, error) {
	response := MetricCPUResponse{}
	serverParts := strings.Split(serverID, ":")
	var url = URL + serverParts[0] + PORT
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	logit.Info.Println("calling cpu with url " + url)
	r, err := http.Post(url+"/metrics/cpu", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// MetricMEMClient client for getting the memory metrics
func MetricMEMClient(serverID string, req *MetricMEMRequest) (MetricMEMResponse, error) {
	response := MetricMEMResponse{}
	serverParts := strings.Split(serverID, ":")
	var url = URL + serverParts[0] + PORT
	logit.Info.Println("calling mem with url " + url)
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(url+"/metrics/mem", "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
