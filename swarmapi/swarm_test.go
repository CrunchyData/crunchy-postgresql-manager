package swarmapi

import (
	"fmt"
	"testing"
)

func TestInspect(t *testing.T) {
	request := new(DockerInspectRequest)
	request.ContainerName = "cpm-admin"

	var response DockerInspectResponse
	var err error

	response, err = DockerInspect(request)
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(response.IPAddress)
}

func TestInfo(t *testing.T) {

	var response DockerInfoResponse
	var err error

	response, err = DockerInfo()
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}

	for i := 0; i < len(response.Output); i++ {
		fmt.Println(response.Output[i])
	}
}

func TestPs(t *testing.T) {

	var response DockerInfoResponse
	var psResponse DockerPsResponse
	var err error

	response, err = DockerInfo()
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}

	for i := 0; i < len(response.Output); i++ {
		fmt.Println(response.Output[i])
		psResponse, err = DockerPs(response.Output[i])
		if err != nil {
			//t.Error(err.Error())
			fmt.Println(err.Error())
			t.Fail()
		}
		for j := 0; j < len(psResponse.Output); j++ {
			fmt.Println(psResponse.Output[j].Name)
		}
	}
}
