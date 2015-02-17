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

package cpmagent

import (
	"bytes"
	"errors"
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"io/ioutil"
	"net/rpc"
	"os/exec"
	"strings"
)

type InspectOutput struct {
	IPAddress    string
	RunningState string
}

type Args struct {
	A, B, C, D, E, CPU, MEM, CommandPath string
}
type DockerRunArgs struct {
	CPU           string
	MEM           string
	ClusterID     string
	ServerID      string
	Image         string
	IPAddress     string
	Standalone    string
	PGDataPath    string
	ContainerName string
	ContainerType string
	CommandOutput string
	CommandPath   string
	EnvVars       map[string]string
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

	glog.Infoln("on server, Command Get called A=" + args.A + " B=" + args.B)
	if args.A == "" {
		glog.Errorln("A was nil")
		return errors.New("Arg A was nil")
	}
	if args.B == "" {
		glog.Infoln("B was nil")
	} else {
		glog.Infoln("B was " + args.B)
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
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	glog.Infoln("command output was " + out.String())
	reply.Output = out.String()

	return nil
}

func (t *Command) ConfigureNode(args *Args, reply *Command) error {

	glog.Infoln("on server, Command ConfigureMaster called CommandPath=" + args.CommandPath + " A=" + args.A + " B=" + args.B + " C=" + args.C + " D=" + args.D + " cpm=" + args.CPU + " mem=" + args.MEM + " E=" + args.E)

	var cmd *exec.Cmd

	cmd = exec.Command(args.CommandPath, args.A, args.B, args.C, args.D, args.CPU, args.MEM, args.E)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	glog.Infoln("command output was " + out.String())
	reply.Output = out.String()

	return nil
}

//
// DockerInspectCommand is to be run on the server that is running
// docker, it will connnect to docker via the unix domain socket
//
func (t *Command) DockerInspectCommand(args *Args, reply *Command) error {

	glog.Infoln("DockerInspectCommand called A=" + args.A)
	if args.A == "" {
		glog.Errorln("A was nil")
		return errors.New("Arg A was nil")
	}

	var cmd *exec.Cmd

	cmd = exec.Command("docker",
		"inspect",
		"--format",
		"{{ .NetworkSettings.IPAddress }}",
		args.A)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}

	ipaddr := strings.Trim(out.String(), "\n")
	glog.Infoln("docker inspect command output was " + ipaddr)
	reply.Output = ipaddr

	return nil
}

//
// DockerRemoveCommand is to be run on the server that is running
// docker, it will connnect to docker via the unix domain socket
// and remove an existing container
//
func (t *Command) DockerRemoveCommand(args *Args, reply *Command) error {

	glog.Infoln("DockerRemoveCommand called A=" + args.A)
	if args.A == "" {
		glog.Errorln("A was nil")
		return errors.New("Arg A was nil")
	}

	var containerName = args.A

	//if a container exists with that name, then we need
	//to stop it first and then remove it
	//inspect
	//stop
	//remove
	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		glog.Errorln("can't get connection to docker socket")
		return err
	}

	//if we can't inspect the container, then we give up
	//on trying to remove it, this is ok to pass since
	//a user can remove the container manually
	container, err3 := docker.InspectContainer(containerName)
	if err3 != nil {
		glog.Errorln("can't inspect container " + containerName)
		glog.Errorln("inspect container error was " + err3.Error())
		return nil
	}

	if container != nil {
		glog.Infoln("container found during inspect")
		err3 = docker.StopContainer(containerName, 10)
		if err3 != nil {
			glog.Infoln("can't stop container " + containerName)
		}
		glog.Infoln("container stopped ")
		opts := dockerapi.RemoveContainerOptions{ID: containerName}
		err := docker.RemoveContainer(opts)
		if err != nil {
			glog.Infoln("can't remove container " + containerName)
		}
		glog.Infoln("container removed ")
	}

	reply.Output = "success"

	return nil
}

//
// Initdb is run on the node during provisioning
//
func (t *Command) PGCommand(args *Args, reply *Command) error {

	glog.Infoln(args.A)

	var cmd *exec.Cmd

	cmd = exec.Command(args.A, "su - postgres -c 'initdb -D /pgdata'")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	glog.Infoln(args.A + " command output was " + out.String())
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
		glog.Errorln(err.Error())
		return err
	}

	return nil
}

//
// generic command with 1 parameter
//
func (t *Command) Command1(args *Args, reply *Command) error {

	glog.Infoln(args.A + " " + args.B)

	var cmd *exec.Cmd

	cmd = exec.Command(args.A, args.B)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	glog.Infoln(args.A + " " + args.B + " command output was " + out.String())
	reply.Output = out.String()
	reply.Output = "success"

	return nil
}

//
// DockerStartCommand is to be run on the server that is running
// docker, it will connnect to docker via the unix domain socket
// and start an existing container
//
func (t *Command) DockerStartCommand(args *Args, reply *Command) error {

	glog.Infoln("DockerStartCommand called A=" + args.A)
	if args.A == "" {
		glog.Errorln("A was nil")
		return errors.New("Arg A was nil")
	}

	var containerName = args.A

	//if a container exists with that name, then we need
	//to stop it first and then remove it
	//inspect
	//start
	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		glog.Errorln("can't get connection to docker socket")
		return err
	}

	//if we can't inspect the container, then we give up
	//on trying to start it
	container, err3 := docker.InspectContainer(containerName)
	if err3 != nil {
		glog.Errorln("can't inspect container " + containerName)
		glog.Errorln("inspect container error was " + err3.Error())
		return errors.New("container " + containerName + " not found")
	}

	if container != nil {
		glog.Infoln("container found during inspect")

		if container.State.Running {
			glog.Infoln("container " + containerName + " was already running, no need to start it")
			return nil
		}

		err3 = docker.StartContainer(containerName, nil)
		if err3 != nil {
			glog.Errorln("can't start container " + containerName)
			return errors.New("can not start container " + containerName)
		}
		glog.Infoln("container started ")
	}

	reply.Output = "success"

	return nil
}

//
// DockerStopCommand is to be run on the server that is running
// docker, it will connnect to docker via the unix domain socket
// and stop an existing container in a running state
//
func (t *Command) DockerStopCommand(args *Args, reply *Command) error {

	glog.Infoln("DockerStopCommand called A=" + args.A)
	if args.A == "" {
		glog.Errorln("A was nil")
		return errors.New("Arg A was nil")
	}

	var containerName = args.A

	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		glog.Errorln("can't get connection to docker socket")
		return err
	}

	//if we can't inspect the container, then we give up
	//on trying to stop it
	container, err3 := docker.InspectContainer(containerName)
	if err3 != nil {
		glog.Infoln("during stop, can't inspect container " + containerName)
		return nil
	}

	if container != nil {
		glog.Infoln("container found during inspect")
		var timeout uint
		timeout = 10
		err3 = docker.StopContainer(containerName, timeout)
		if err3 != nil {
			glog.Infoln("can't stop container " + containerName)
		}
		glog.Infoln("container " + containerName + " stopped ")
	}

	reply.Output = "success"

	return nil
}

//
// DockerInspectFullCommand is to be run on the server that is running
// docker, it will connnect to docker via the unix domain socket
//
func (t *Command) DockerInspect2Command(args *Args, reply *InspectCommandOutput) error {

	glog.Infoln("DockerInspect2Command called containerName=" + args.A)
	if args.A == "" {
		glog.Errorln("containerName was nil")
		return errors.New("containerName was nil")
	}

	var containerName = args.A

	docker, err := dockerapi.NewClient("unix://var/run/docker.sock")
	if err != nil {
		glog.Errorln("can't get connection to docker socket")
		return err
	}

	//if we can't inspect the container, then we give up
	//on trying to stop it
	reply.RunningState = "down"
	reply.IPAddress = ""

	container, err3 := docker.InspectContainer(containerName)
	if err3 != nil {
		glog.Infoln("can't inspect container " + containerName)
		return err3
	}

	if container != nil {
		glog.Infoln("container found during inspect")
		if container.State.Running {
			reply.RunningState = "up"
			glog.Infoln("container status is up")
			glog.Infoln("container ipaddress is " + container.NetworkSettings.IPAddress)
			reply.IPAddress = container.NetworkSettings.IPAddress
		} else {
			reply.RunningState = "down"
			glog.Infoln("container status is down")
		}
	}

	return nil
}

func (t *Command) DockerRun(args *DockerRunArgs, reply *Command) error {

	glog.Infoln("on server, Command DockerRun called CommandPath=" + args.CommandPath + " PGDataPath=" + args.PGDataPath + " ContainerName=" + args.ContainerName + " Image=" + args.Image + " cpm=" + args.CPU + " mem=" + args.MEM)

	var cmd *exec.Cmd

	var allEnvVars = ""
	if args.EnvVars != nil {
		for k, v := range args.EnvVars {
			allEnvVars = allEnvVars + " -e " + k + "=" + v
		}
	}
	glog.Infoln("env vars " + allEnvVars)

	cmd = exec.Command(args.CommandPath, args.PGDataPath, args.ContainerName, args.Image, args.CPU, args.MEM, allEnvVars)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		glog.Errorln(err.Error())
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		return errors.New(errorString)
	}
	glog.Infoln("command output was " + out.String())
	reply.Output = out.String()

	return nil
}

func AgentCommand(arga string, argb string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("AgentCommand: dialing:" + err.Error())
	}
	if client == nil {
		glog.Errorln("AgentCommand: dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	args := &Args{}
	args.A = arga
	args.B = argb
	err = client.Call("Command.Get", args, &command)
	if err != nil {
		glog.Errorln("AgentCommand:" + arga + " Command Get error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func AgentDockerRun(args DockerRunArgs, ipaddress string) (string, error) {
	glog.Infoln("server.go:AgentDockerRun:cpu=" + args.CPU + " mem=" + args.MEM)
	if args.EnvVars != nil {
		for k, v := range args.EnvVars {
			glog.Infoln("server.go: env var " + k + " " + v)
		}
	}
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("AgentDockerRun dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("AgentDockerRun dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	err = client.Call("Command.DockerRun", args, &command)
	if err != nil {
		glog.Errorln("DockerRun error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func DockerInspectCommand(arga string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("DockerInspectCommand: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("DockerInspectCommand: dialing: client was nil")
		return "", errors.New("client was nil")
	}

	var command Command

	args := &Args{}
	args.A = arga
	err = client.Call("Command.DockerInspectCommand", args, &command)
	if err != nil {
		glog.Errorln("DockerInspectCommand:" + arga + " Command DockerInspectCommand error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func DockerInspect2Command(arga string, ipaddress string) (InspectOutput, error) {
	var command InspectCommandOutput
	var output InspectOutput

	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("DockerInspect2Command: dialing:" + err.Error())
		return output, err
	}
	if client == nil {
		glog.Errorln("DockerInspect2Command dialing: client was nil")
		return output, errors.New("client was nil here")
	}

	args := &Args{}
	args.A = arga
	err = client.Call("Command.DockerInspect2Command", args, &command)
	if err != nil {
		glog.Errorln("DockerInspect2Command:" + arga + " Command DockerInspect2Command error:" + err.Error())
		return output, err
	}
	output.IPAddress = command.IPAddress
	output.RunningState = command.RunningState

	return output, nil
}

func DockerRemoveContainer(arga string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("DockerRemoveContainer: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("DockerRemoveContainer: dialing:" + err.Error())
		return "", errors.New("client was nil here2")
	}

	var command Command

	args := &Args{}
	args.A = arga
	err = client.Call("Command.DockerRemoveCommand", args, &command)
	if err != nil {
		glog.Errorln("DockerRemoveCommand arga error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func AgentCommandConfigureNode(commandpath string, arga string, argb string, argc string, argd string, arge string, cpu string, mem string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("AgentCommandConfigureNode dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("AgentCommandConfigureNode dialing: client was nil")
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
		glog.Errorln("AgentCommandConfigureNode error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func DockerStartContainer(containerName string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("DockerStartContainer dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("DockerStartContainer dialing:" + err.Error())
		return "", errors.New("client was nil 3")
	}

	var command Command

	args := &Args{}
	args.A = containerName
	err = client.Call("Command.DockerStartCommand", args, &command)
	if err != nil {
		glog.Errorln("DockerStartContainer containerName error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func DockerStopContainer(containerName string, ipaddress string) (string, error) {

	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("DockerStopContainer dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("DockerStopContainer dialing:" + err.Error())
		return "", errors.New("client was nil here 4")
	}

	var command Command

	args := &Args{}
	args.A = containerName
	err = client.Call("Command.DockerStopCommand", args, &command)
	if err != nil {
		glog.Infoln("DockerStopContainer containerName error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func Command1(pgcommand string, parameter string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("Command1: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("Command1: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	args := &Args{}
	args.A = pgcommand
	args.B = parameter
	err = client.Call("Command.Command1", args, &command)
	if err != nil {
		glog.Errorln("Command1: " + args.A + " " + args.B + " error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}
