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

package kubeclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"io/ioutil"
	"net/http"
)

type MyPod struct {
	CurrentState struct {
		Status string
	}
}

func getHttpClient() (*http.Client, error) {
	var caFile = "/kubekeys/root.crt"
	var certFile = "/kubekeys/cert.crt"
	var keyFile = "/kubekeys/key.key"

	var client *http.Client

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logit.Error.Println(err.Error())
		return client, err
	}
	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		logit.Error.Println(err.Error())
		return client, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client = &http.Client{Transport: transport}

	return client, nil
}

// DeleteService deletes a kube service
// kubeURL  - the URL to kube
// ID - the ID of the service we want to delete
// it returns an error is there was a problem
func DeleteService(kubeURL string, ID string) error {
	logit.Info.Println("deleting service " + ID)

	client, err := getHttpClient()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	// DELETE service
	var url = kubeURL + "/api/v1beta1/services/" + ID
	logit.Info.Println("url is " + url)
	request, err2 := http.NewRequest("DELETE", url, nil)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		return err2
	}

	resp, err := client.Do(request)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		logit.Error.Println(err3.Error())
		return err3
	}
	logit.Info.Println(string(data))

	return nil
}

// DeletePod deletes a kube pod that should already exist
// kubeURL  - the URL to kube
// ID - the ID of the Pod we want to delete
// it returns an error is there was a problem
func DeletePod(kubeURL string, ID string) error {
	logit.Info.Println("deleting pod " + ID)

	client, err4 := getHttpClient()
	if err4 != nil {
		logit.Error.Println(err4.Error())
		return err4
	}

	// DELETE pod
	var url = kubeURL + "/api/v1beta1/pods/" + ID
	logit.Info.Println("url is " + url)
	request, err2 := http.NewRequest("DELETE", url, nil)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		return err2
	}

	resp, err := client.Do(request)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		return err2
	}
	logit.Info.Println(string(data))

	return nil
}

// CreatePod creates a new pod and service using passed in values
// kubeURL - the URL to the kube
// podInfo - the params used to configure the pod
// return an error if anything goes wrong
func CreatePod(kubeURL string, podInfo template.KubePodParams) error {
	client, err := getHttpClient()

	logit.Info.Println("creating pod " + podInfo.ID)

	//use a pod template to build the pod definition
	data, err := template.KubeNodePod(podInfo)
	if err != nil {
		logit.Error.Println("CreatePod:" + err.Error())
		return err
	}
	logit.Info.Println(string(data[:]))

	var bodyType = "application/json"
	var url = kubeURL + "/api/v1beta1/pods"
	logit.Info.Println("url is " + url)

	// POST POD
	resp, err := client.Post(url, bodyType, bytes.NewReader(data))
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println(string(data))

	return nil

}

// GetPods gets all the pods
// kubeURL - the URL to the kube
// podInfo - the params used to configure the pod
// return an error if anything goes wrong
func GetPods(kubeURL string, podInfo template.KubePodParams) error {
	logit.Info.Println("creating pod " + podInfo.ID)

	client, err := getHttpClient()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//use a pod template to build the pod definition
	data, err := template.KubeNodePod(podInfo)
	if err != nil {
		logit.Error.Println("CreatePod:" + err.Error())
		return err
	}

	logit.Info.Println(string(data[:]))

	// Do GET something
	resp, err2 := client.Get(kubeURL + "/api/v1beta1/pods")
	if err2 != nil {
		logit.Error.Println(err2.Error())
		return err2
	}
	defer resp.Body.Close()

	// Dump response
	data, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		logit.Error.Println(err3.Error())
		return err3
	}
	logit.Info.Println(string(data))

	return nil
}

// GetPod gets information about a single pod from kube
// kubeURL - the URL to the kube
// podName - the pod name
// return an error if anything goes wrong
func GetPod(kubeURL string, podName string) (MyPod, error) {
	var podInfo MyPod

	logit.Info.Println("getting pod info " + podName)

	client, err := getHttpClient()
	if err != nil {
		logit.Error.Println(err.Error())
		return podInfo, err
	}

	// Do GET something
	resp, err2 := client.Get(kubeURL + "/api/v1beta1/pods/" + podName)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		return podInfo, err2
	}
	defer resp.Body.Close()

	// Dump response
	data, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		logit.Error.Println(err3.Error())
		return podInfo, err3
	}
	logit.Info.Println(string(data))
	err2 = json.Unmarshal(data, &podInfo)
	if err2 != nil {
		logit.Error.Println("error in unmarshalling pod " + err2.Error())
		return podInfo, err2
	}

	return podInfo, nil
}

// CreateService creates a service
// kubeURL - the URL to the kube
// podName - the pod name
// return an error if anything goes wrong
func CreateService(kubeURL string, podInfo template.KubePodParams) error {
	var s1data []byte
	var err error
	var serviceurl = kubeURL + "/api/v1beta1/services"
	var bodyType = "application/json"

	logit.Info.Println("create service called")

	client, err := getHttpClient()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	s1data, err = template.KubeNodeService(podInfo)
	if err != nil {
		logit.Error.Println("CreateService:" + err.Error())
		return err
	}
	logit.Info.Println("create service request...")
	logit.Info.Println(string(s1data[:]))

	// POST admin SERVICE at port 13000
	resp1, err1 := client.Post(serviceurl, bodyType, bytes.NewReader(s1data))
	if err1 != nil {
		logit.Error.Println(err1.Error())
		return err1
	}
	defer resp1.Body.Close()

	// Dump response
	data, err4 := ioutil.ReadAll(resp1.Body)
	if err4 != nil {
		logit.Error.Println(err4.Error())
		return err4
	}
	logit.Info.Println("create service response..." + string(data))

	return nil
}
