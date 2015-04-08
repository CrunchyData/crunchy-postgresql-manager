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

package dummy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"os/exec"
)

type Args struct {
	A, B, C, D, E, CommandPath string
}

type Command struct {
	Output string
}

type InspectCommandOutput struct {
	IPAddress    string
	RunningState string
}

/*
 args.A contains the name of what we are going to execute
 such as 'iostat.sh' or 'df.sh'
*/
func (t *Command) Get(args *Args, reply *Command) error {

	logit.Info.Println("on server, Command Get called A=" + args.A + " B=" + args.B)
	if args.A == "" {
		logit.Error.Println("A was nil")
		return errors.New("Arg A was nil")
	}
	if args.B == "" {
		logit.Info.Println("B was nil")
	} else {
		logit.Info.Println("B was " + args.B)
	}

	var cmd *exec.Cmd

	if args.B == "" {
		cmd = exec.Command(args.A)
	} else {
		cmd = exec.Command(args.A, args.B)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	logit.Info.Println("command output was " + out.String())
	reply.Output = out.String()

	return nil
}
