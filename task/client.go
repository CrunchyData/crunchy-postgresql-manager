package task

import (
	"bytes"
	"encoding/json"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"io/ioutil"
	"net/http"
)

const URL = "http://cpm-task:13001"

//
// example of calling one of these
//
// request := &backup.Reload{"something", "yes"}
// response, err := backup.ReloadClient("http://cpm-task:13001", request)
//
func ReloadClient() (ReloadResponse, error) {
	var req = ReloadRequest{}
	response := ReloadResponse{}

	req.Name = "foo"
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/reload"
	logit.Info.Println("reload about to post to " + url)
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

func ExecuteNowClient(req *TaskRequest) (ExecuteNowResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/executenow"
	logit.Info.Println("executenow about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := ExecuteNowResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StatusUpdateClient(req *TaskStatus) (StatusUpdateResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/status/update"
	logit.Info.Println("statusupdate client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StatusUpdateResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

func StatusAddClient(req *TaskStatus) (StatusAddResponse, error) {
	var err error
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/status/add"
	logit.Info.Println("statusadd client about to post to " + url)
	r, _ := http.Post(url, "application/json", body)
	rawresponse, _ := ioutil.ReadAll(r.Body)
	response := StatusAddResponse{}
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
