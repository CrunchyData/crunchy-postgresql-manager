//
// example of calling one of these
//
// request := &backup.Reload{"something", "yes"}
// response, err := backup.ReloadClient("http://cpm-task:13001", request)
//
package task

import (
	"bytes"
	"encoding/json"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"io/ioutil"
	"net/http"
)

const URL = "http://cpm-task:13001"

// ReloadClient is the client to the reload function which causes the scheduled tasks
// to be reread from the database
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

// ExecuteNowClient client for the execute now function which runs a task immediately
func ExecuteNowClient(req *TaskRequest) (ExecuteNowResponse, error) {
	response := ExecuteNowResponse{}
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/executenow"
	logit.Info.Println("executenow about to post to " + url)
	r, err := http.Post(url, "application/json", body)
	if err != nil {
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// StatusUpdateClient is the client to the StatusUpdate function
func StatusUpdateClient(req *TaskStatus) (StatusUpdateResponse, error) {
	response := StatusUpdateResponse{}
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/status/update"
	logit.Info.Println("statusupdate client about to post to " + url)
	r, err := http.Post(url, "application/json", body)
	if err != nil {
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}

// StatusAddClient is the client to the StatusAdd function which allows
// for the addition of new status info
func StatusAddClient(req *TaskStatus) (StatusAddResponse, error) {
	response := StatusAddResponse{}
	buf, _ := json.Marshal(req)
	body := bytes.NewBuffer(buf)
	url := URL + "/api/status/add"
	logit.Info.Println("statusadd client about to post to " + url)
	r, err := http.Post(url, "application/json", body)
	if err != nil {
		return response, err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	//fmt.Println(string(rawresponse))
	return response, err
}
