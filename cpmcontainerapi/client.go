package cpmcontainerapi

import (
	"bytes"
	"encoding/json"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"io/ioutil"
	"net/http"
)

const PORT = ":10001"

//
// example of calling one of these
//
// request := &cpmcontainerapi.RemovewritefileRequest{"something", "yes"}
// response, err := cpmcontainerapi.RemoteWritefileClient("http://localhost:10001", request)
//
func RemoteWritefileClient(path string, contents string, ipaddress string) (RemoteWritefileResponse, error) {
	var req = RemoteWritefileRequest{}
	response := RemoteWritefileResponse{}
	req.Path = path
	req.Filecontents = contents
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + ipaddress + PORT + "/api/remotewritefile"
	logit.Info.Println("remotewritefile about to post to " + url)
	r, err := http.Post(url, "application/json", body)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	rawresponse, e2 := ioutil.ReadAll(r.Body)
	if e2 != nil {
		logit.Error.Println(e2.Error())
		return response, e2
	}
	err = json.Unmarshal(rawresponse, &response)
	if err != nil {
		logit.Error.Println(err.Error())
		return response, err
	}
	logit.Info.Println(string(rawresponse))
	return response, err
}

func InitdbClient(host string) (InitdbResponse, error) {
	var err error
	req := InitdbRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/initdb"
	logit.Info.Println("initdbclient about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := InitdbResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StartPGClient(host string) (StartPGResponse, error) {
	var err error
	req := StartPGRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/startpg"
	logit.Info.Println("startpg client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StartPGResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StartPGOnStandbyClient(host string) (StartPGOnStandbyResponse, error) {
	var err error
	req := StartPGOnStandbyRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/startpgonstandby"
	logit.Info.Println("startpgonstandby client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StartPGOnStandbyResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StopPGClient(host string) (StopPGResponse, error) {
	var err error
	req := StopPGRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/stoppg"
	logit.Info.Println("stoppg client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StopPGResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StartPgpoolClient(host string) (StartPgpoolResponse, error) {
	var err error
	req := StartPgpoolRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/startpgpool"
	logit.Info.Println("startpgpool client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StartPgpoolResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StopPgpoolClient(host string) (StopPgpoolResponse, error) {
	var err error
	req := StopPgpoolRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/stoppgpool"
	logit.Info.Println("stoppgpool client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StopPgpoolResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func BasebackupClient(master string, standby string) (BasebackupResponse, error) {
	var err error
	req := BasebackupRequest{}
	req.MasterHostName = master
	req.StandbyHostName = standby

	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + standby + PORT + "/api/basebackup"
	logit.Info.Println("stoppgpool client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := BasebackupResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func FailoverClient(host string) (FailoverResponse, error) {
	var err error
	req := FailoverRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/failover"
	logit.Info.Println("failover client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := FailoverResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func SeedClient(host string) (SeedResponse, error) {
	var err error
	req := SeedRequest{}
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/seed"
	logit.Info.Println("seed client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := SeedResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func ControldataClient(host string) (ControldataResponse, error) {
	var err error
	req := ControldataRequest{}
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/controldata"
	logit.Info.Println("controldata client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := ControldataResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func BadgerGenerateClient(host string) (BadgerGenerateResponse, error) {
	var err error
	req := BadgerGenerateRequest{}
	req.ContainerName = host
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := "http://" + host + PORT + "/api/badgergenerate"
	logit.Info.Println("badgergenerate client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := BadgerGenerateResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
