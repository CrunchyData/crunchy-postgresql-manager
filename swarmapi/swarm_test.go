package swarmapi

import (
	"fmt"
	"testing"
)

var TestContainerName = "dfw1"

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

func TestRun(t *testing.T) {

	request := new(DockerRunRequest)
	var response DockerRunResponse
	var err error

	request.ContainerName = TestContainerName
	request.PGDataPath = "/tmp/" + TestContainerName
	request.Profile = "SM"
	request.Image = "cpm-node"
	request.EnvVars = make(map[string]string)
	request.CPU = "0"
	request.MEM = "0"
	response, err = DockerRun(request)
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(response.ID)
}

func TestStop(t *testing.T) {

	var response DockerStopResponse
	var err error

	request := new(DockerStopRequest)
	request.ContainerName = TestContainerName
	response, err = DockerStop(request)
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(response.Output)
}

func TestStart(t *testing.T) {

	request := new(DockerStartRequest)
	var response DockerStartResponse
	var err error

	request.ContainerName = TestContainerName
	response, err = DockerStart(request)
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(response.Output)
}

func TestInspect(t *testing.T) {
	request := new(DockerInspectRequest)
	request.ContainerName = TestContainerName

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

func TestRemove(t *testing.T) {
	TestStop(t)
	request := new(DockerRemoveRequest)
	request.ContainerName = TestContainerName

	var err error

	_, err = DockerRemove(request)
	if err != nil {
		//t.Error(err.Error())
		fmt.Println(err.Error())
		t.Fail()
	}
}
