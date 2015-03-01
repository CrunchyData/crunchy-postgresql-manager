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
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
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

	glog.Infoln("dockerapi started...")
	time.Sleep(time.Millisecond * 5000)
	glog.Infoln("action: " + ACTION)
	glog.Infoln("containers: " + CONTAINERS)
	var containers = strings.Split(CONTAINERS, ",")
	var docker *dockerapi.Client
	var err error
	docker, err = dockerapi.NewClient(DOCKER_HOST)
	if err != nil {
		glog.Errorln(err.Error())
	}
	err = docker.Ping()
	glog.Flush()

	for i := range containers {
		glog.Infoln(ACTION + " issued for " + containers[i])
		switch ACTION {
		case "stop":
			err = docker.StopContainer(containers[i], 5)
		case "start":
			err = docker.StartContainer(containers[i], nil)
		default:
			glog.Infoln(ACTION + " unsupported action")
		}
		if err != nil {
			glog.Errorln(err.Error())
		}
		time.Sleep(time.Millisecond * 2000)
	}
	glog.Flush()
}
