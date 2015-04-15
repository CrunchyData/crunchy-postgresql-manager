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
	"flag"
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"strings"
	"time"
)

var ACTION string
var CONTAINERS string
var DOCKER_HOST = "unix:///var/run/docker.sock"

func init() {
	flag.StringVar(&ACTION, "a", "", "action to perform, valid values are start and stop")
	flag.StringVar(&CONTAINERS, "c", "", "container names to act upon")

	flag.Parse()
}

func main() {

	fmt.Println("dockerapi started...")
	time.Sleep(time.Millisecond * 5000)
	fmt.Println("action: " + ACTION)
	fmt.Println("containers: " + CONTAINERS)
	var containers = strings.Split(CONTAINERS, ",")
	var docker *dockerapi.Client
	var err error
	docker, err = dockerapi.NewClient(DOCKER_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = docker.Ping()

	for i := range containers {
		fmt.Println(ACTION + " issued for " + containers[i])
		switch ACTION {
		case "stop":
			err = docker.StopContainer(containers[i], 5)
		case "start":
			err = docker.StartContainer(containers[i], nil)
		default:
			fmt.Println(ACTION + " unsupported action")
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Millisecond * 2000)
	}
}
