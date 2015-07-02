package main

import (
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	testMetricIostat()
	//testMetricDf()
	//testDiskProvision()
	//testDiskDelete()
	//testDockerRun()
	//testDockerInspect()
	//testDockerRemove()
	//testDockerStop()
	//testDockerStart()
	//testMetricCPU()
	//testMetricMEM()
}

func testMetricDf() {
	var response cpmserverapi.MetricDfResponse
	var err error
	request := &cpmserverapi.MetricDfRequest{"something", "yes"}
	response, err = cpmserverapi.MetricDfClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}

func testMetricIostat() {
	var response cpmserverapi.MetricIostatResponse
	var err error
	request := &cpmserverapi.MetricIostatRequest{"something", "yes"}
	response, err = cpmserverapi.MetricIostatClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDiskProvision() {
	var response cpmserverapi.DiskProvisionResponse
	var err error
	request := &cpmserverapi.DiskProvisionRequest{"/tmp/foo"}
	response, err = cpmserverapi.DiskProvisionClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDiskDelete() {
	var response cpmserverapi.DiskDeleteResponse
	var err error
	request := &cpmserverapi.DiskDeleteRequest{"/tmp/foo"}
	response, err = cpmserverapi.DiskDeleteClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDockerRun() {
	var response cpmserverapi.DockerRunResponse
	var err error
	request := &cpmserverapi.DockerRunRequest{}
	request.CommandPath = "docker-run.sh"
	request.CPU = "0"
	request.MEM = "0"
	request.Image = "cpm-node"
	request.PGDataPath = "/tmp/foo"
	request.ContainerName = "testpoo"
	envvars := make(map[string]string)
	envvars["one"] = "uno"
	request.EnvVars = envvars
	response, err = cpmserverapi.DockerRunClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDockerInspect() {
	var response cpmserverapi.DockerInspectResponse
	var err error
	request := &cpmserverapi.DockerInspectRequest{}
	request.ContainerName = "testpoo"
	response, err = cpmserverapi.DockerInspectClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.IPAddress + " " + response.RunningState)
}
func testDockerRemove() {
	var response cpmserverapi.DockerRemoveResponse
	var err error
	request := &cpmserverapi.DockerRemoveRequest{}
	request.ContainerName = "testpoo"
	response, err = cpmserverapi.DockerRemoveClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDockerStop() {
	var response cpmserverapi.DockerStopResponse
	var err error
	request := &cpmserverapi.DockerStopRequest{}
	request.ContainerName = "testpoo"
	response, err = cpmserverapi.DockerStopClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testDockerStart() {
	var response cpmserverapi.DockerStartResponse
	var err error
	request := &cpmserverapi.DockerStartRequest{}
	request.ContainerName = "testpoo"
	response, err = cpmserverapi.DockerStartClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}

func testMetricMEM() {
	var response cpmserverapi.MetricMEMResponse
	var err error
	request := &cpmserverapi.MetricMEMRequest{}
	response, err = cpmserverapi.MetricMEMClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
func testMetricCPU() {
	var response cpmserverapi.MetricCPUResponse
	var err error
	request := &cpmserverapi.MetricCPURequest{}
	response, err = cpmserverapi.MetricCPUClient("http://192.168.0.106:10001", request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("returned " + response.Output)
}
