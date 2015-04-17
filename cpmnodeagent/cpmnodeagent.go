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

package cpmnodeagent

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	//"github.com/crunchydata/crunchy-postgresql-manager/util"
	"io/ioutil"
	"net/rpc"
	"os/exec"
)

type Args struct {
	A, B, C, D, E, CPU, MEM, CommandPath string
}

type Command struct {
	Output string
}

func (t *Command) ConfigureNode(args *Args, reply *Command) error {

	logit.Info.Println("on server, Command ConfigureMaster called CommandPath=" + args.CommandPath + " A=" + args.A + " B=" + args.B + " C=" + args.C + " D=" + args.D + " cpm=" + args.CPU + " mem=" + args.MEM + " E=" + args.E)

	var cmd *exec.Cmd

	cmd = exec.Command(args.CommandPath, args.A, args.B, args.C, args.D, args.CPU, args.MEM, args.E)
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

//
// Initdb is run on the node during provisioning
//
func (t *Command) PGCommand(args *Args, reply *Command) error {

	logit.Info.Println(args.A)

	var cmd *exec.Cmd

	cmd = exec.Command(args.A, "su - postgres -c 'initdb -D /pgdata'")

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
	logit.Info.Println(args.A + " command output was " + out.String())
	reply.Output = out.String()
	reply.Output = "success"

	return nil
}

//
// Writefile is run on the node during provisioning
// args.A contains the file contents
// args.B contains the path to write the contents to
//
func (t *Command) Writefile(args *Args, reply *Command) error {
	d1 := []byte(args.A)
	err := ioutil.WriteFile(args.B, d1, 0644)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return nil
}

//
// generic command with 1 parameter
//
func (t *Command) Command1(args *Args, reply *Command) error {

	logit.Info.Println(args.A + " " + args.B)

	var cmd *exec.Cmd

	cmd = exec.Command(args.A, args.B)

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
	logit.Info.Println(args.A + " " + args.B + " command output was " + out.String())
	reply.Output = out.String()
	reply.Output = "success"

	return nil
}

func AgentCommand(arga string, argb string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("AgentCommand: dialing:" + err.Error())
	}
	if client == nil {
		logit.Error.Println("AgentCommand: dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	args := &Args{}
	//args.A = util.GetBase() + "/" + arga
	args.A = arga
	args.B = argb
	err = client.Call("Command.Get", args, &command)
	if err != nil {
		logit.Error.Println("AgentCommand:" + arga + " Command Get error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

//used to execute pg_controldata PG binary
func PostgresCommand(arga string, argb string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("PostgresCommand: dialing:" + err.Error())
	}
	if client == nil {
		logit.Error.Println("PostgresCommand: dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	args := &Args{}
	//args.A = util.GetPostgresBase() + "/bin/" + arga
	args.A = arga
	args.B = argb
	err = client.Call("Command.Get", args, &command)
	if err != nil {
		logit.Error.Println("PostgresCommand:" + arga + " Command Get error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func AgentCommandConfigureNode(commandpath string, arga string, argb string, argc string, argd string, arge string, cpu string, mem string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("AgentCommandConfigureNode dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("AgentCommandConfigureNode dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	args := &Args{}
	args.CommandPath = commandpath
	args.A = arga
	args.B = argb
	args.C = argc
	args.D = argd
	args.CPU = cpu
	args.MEM = mem
	args.E = arge
	err = client.Call("Command.ConfigureNode", args, &command)
	if err != nil {
		logit.Error.Println("AgentCommandConfigureNode error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

//used by adminapi clustermgmt
func Command2(pgcommand string, parameter string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("Command2: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("Command2: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	args := &Args{}
	args.A = pgcommand
	args.B = parameter
	err = client.Call("Command.Command2", args, &command)
	if err != nil {
		logit.Error.Println("Command2: " + args.A + " " + args.B + " error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func (t *Command) Command2(args *Args, reply *Command) error {

	logit.Info.Println(args.A + " " + args.B)

	var cmd *exec.Cmd

	cmd = exec.Command(args.A, args.B)

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
	logit.Info.Println(args.A + " " + args.B + " command output was " + out.String())
	reply.Output = out.String()
	reply.Output = "success"

	return nil
}

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
