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
	"crunchy.com/template"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	//client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/golang/glog"
)

// DeletePod deletes a kube pod that should already exist
// kubeURL  - the URL to kube
// ID - the ID of the Pod we want to delete
// it returns an error is there was a problem
func DeletePod(kubeURL string, ID string) error {
	glog.Infoln("deleting pod " + ID)
	/**
	var c *client.Client
	c = client.NewOrDie(&client.Config{
		Host:    kubeURL,
		Version: "v1beta1",
	})
	if c != nil {
		glog.Infoln("connection to kube ok....")
	}
	err := c.Pods(api.NamespaceDefault).Delete(ID)
	if err != nil {
		glog.Errorln("DeletePod:" + err.Error())
		return err
	}
	*/

	return nil
}

// CreatePod creates a new pod using passed in values
// kubeURL - the URL to the kube
// podInfo - the params used to configure the pod
// return an error if anything goes wrong
func CreatePod(kubeURL string, podInfo template.KubePodParams) error {
	glog.Infoln("creating pod " + podInfo.ID)

	//use a pod template to build the pod definition
	data, err := template.KubeNodePod(podInfo)
	if err != nil {
		glog.Errorln("CreatePod:" + err.Error())
		return err
	}

	glog.Infoln(string(data[:]))

	/*
		//use the kube api directly for now, later on probably an openshift wrapping api
		var obj runtime.Object
		wasCreated := true
		var c *client.Client
		c = client.NewOrDie(&client.Config{
			Host:    kubeURL,
			Version: "v1beta1",
		})
		if c != nil {
			glog.Infoln("connection to kube ok....")
		}

		obj, err = c.Verb("POST").Path("pods").Body(data).Do().WasCreated(&wasCreated).Get()
		if err != nil {
			glog.Errorln("CreatePod:" + err.Error())
			return err
		}
		if obj != nil {
			glog.Infoln("got the object from the kube pod create")
		}
	*/
	return nil
}
