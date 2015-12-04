package adminapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	response := LoginToken{}
	var PORT = ":13001"
	var host = "cpm-admin"
	t.Log("getting a token")
	url := "http://" + host + PORT + "/sec/login/cpm.cpm"
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	fmt.Println(response.Contents)
	if response.Contents == "" {
		t.Error("token was blank")
	} else if len(response.Contents) != 16 {
		t.Error("token was not 16 bytes long")
	}
}
