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
	"github.com/crunchydata/crunchy-postgresql-manager/dummy"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"net"
	"net/http"
	"net/rpc"
)

func main() {

	logit.Info.Println("dummy starting\n")
	command := new(dummy.Command)
	rpc.Register(command)
	logit.Info.Println("Command registered\n")
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":13013")
	logit.Info.Println("listening\n")
	if e != nil {
		logit.Error.Println(e.Error())
		panic("could not listen on rpc socker")
	}
	logit.Info.Println("about to serve\n")
	http.Serve(l, nil)
	logit.Info.Println("after serve\n")
}
