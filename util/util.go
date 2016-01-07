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

package util

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var CPMBASE string

func GetBase() string {
	if CPMBASE == "" {
		CPMBASE = os.Getenv("CPMBASE")
		if CPMBASE == "" {
			CPMBASE = "/var/cpm"
		}
	}
	return CPMBASE
}

func CleanName(name string) string {
	return strings.Replace(strings.ToLower(name), " ", "", -1)
}

// FastPing returns either OFFLINE or RUNNING for a host and port combination
func FastPing(port string, host string) (string, error) {

	var err error

	var cmd *exec.Cmd
	cmd = exec.Command("ping-wrapper.sh", host, port)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "OFFLINE", err
	}
	var rc = out.String()
	if rc != "connected" {
		return "OFFLINE", nil
	}

	return "RUNNING", nil
}
