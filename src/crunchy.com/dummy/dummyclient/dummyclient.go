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
	"crunchy.com/dummy"
	"crunchy.com/logutil"
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

	logutil.Log("connecting to " + *hostFlag)
	if *metricFlag == "none" {
		log.Fatal("no metric specified on command line -m")
	}

	client, err := rpc.DialHTTP("tcp", *hostFlag+":13013")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	if client == nil {
		log.Fatal("dialing:", err)
	}

	var command dummy.Command

	logutil.Log("executing metric: " + *metricFlag)

	args := &dummy.Args{}
	args.A = *metricFlag
	err = client.Call("Command.Get", args, &command)
	if err != nil {
		logutil.Log("Get error:" + err.Error())
	} else {
		logutil.Log(" output ......")
		logutil.Log(command.Output)
		logutil.Log(" ")
	}

}
