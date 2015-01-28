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

package mon

import (
	"crunchy.com/logutil"
	"errors"
	"net/rpc"
)

//just a placeholder for any client calls to the monitor server
//in the future
func PlaceholderClient(ipaddress string, status string) error {

	logutil.Log("Placeholder called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("Placeholder: dialing:" + err.Error())
		return err
	}
	if client == nil {
		logutil.Log("Placeholder: client was nil")
		return errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.Placeholder", &status, &command)
	if err != nil {
		logutil.Log("Placeholder: error " + err.Error())
		return err
	}
	logutil.Log("status=" + status)

	return nil
}
