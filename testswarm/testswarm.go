/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
)

func main() {

	fmt.Println("at top of testswarm main")

	var err error
	inspectReq := swarmapi.DockerInspectRequest{}
	inspectReq.ContainerName = "cpm"
	var inspectResp swarmapi.DockerInspectResponse
	inspectResp, err = swarmapi.DockerInspect(&inspectReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(inspectResp.IPAddress)

	runReq := swarmapi.DockerRunRequest{}
	runReq.PGDataPath = "/var/cpm/data/pgsql/swarmtest"
	runReq.ContainerType = "cpm-node"
	runReq.ContainerName = "swarmtest"
	runReq.EnvVars = make(map[string]string)
	runReq.EnvVars["one"] = "value of one"
	runReq.EnvVars["two"] = "value of two"
	runReq.CPU = "0"
	runReq.MEM = "0"
	var runResp swarmapi.DockerRunResponse
	runResp, err = swarmapi.DockerRun(&runReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(runResp.ID)
}
