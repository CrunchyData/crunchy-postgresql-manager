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
	"crunchy.com/cpmagent"
	"crunchy.com/logit"
	"flag"
	"log"
	"net/rpc"
)

/*
  command line args:
 -h ipaddress to connect to, localhost if none supplied
 -m metric to execute (required)
*/
func main() {
	hostFlag := flag.String("h", "127.0.0.1", "ipaddress of cpm server")
	metricFlag := flag.String("m", "none", "metric to execute")
	flag.Parse()

	logit.Info.Println("main: connecting to " + *hostFlag)
	if *metricFlag == "none" {
		log.Fatal("no metric specified on command line -m")
	}

	client, err := rpc.DialHTTP("tcp", *hostFlag+":13000")
	if err != nil {
		log.Fatal("main:dialing:" + err.Error())
	}
	if client == nil {
		log.Fatal("main:dialing:" + err.Error())
	}

	var command cpmagent.Command

	logit.Info.Println("main: executing metric: " + *metricFlag)

	args := &cpmagent.Args{}
	args.A = *metricFlag
	err = client.Call("Command.Get", args, &command)
	if err != nil {
		logit.Error.Println("main:Get error:" + err.Error())
	} else {
		logit.Info.Println("main: output ......")
		logit.Info.Println("main:" + command.Output)
		logit.Info.Println("main: ")
	}

}
