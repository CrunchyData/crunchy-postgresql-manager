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
	"bytes"
	"crunchy.com/template"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"io/ioutil"
	"log"
	"net/http"
)

type MyPod struct {
	CurrentState struct {
		Status string
	}
}

func TestCreate(w rest.ResponseWriter, r *rest.Request) {

	glog.Infoln("here in Test Create")

	podInfo := template.KubePodParams{
		"testnode",
		"0", "0",
		"crunchydata/cpm-node",
		"/opt/cpm/data/pgsql/testnode", "13000"}
	err := CreatePod(kubeURL, podInfo)
	if err != nil {
		glog.Infoln(err.Error())
		glog.Errorln("TestCreate:" + err.Error())
	}
	glog.Infoln("no error on create pod")

	response := KubeResponse{}
	response.URL = "here in TestCreate"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&response)
}

func TestDelete(w rest.ResponseWriter, r *rest.Request) {

	glog.Infoln("here in Test Delete")
	err := DeletePod(kubeURL, "testnode")
	if err != nil {
		glog.Infoln(err.Error())
		glog.Errorln("TestCreate:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	glog.Infoln("no error on delete pod")
	response := KubeResponse{}
	response.URL = "here in TestDelete"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&response)
}

// DeletePod deletes a kube pod that should already exist
// kubeURL  - the URL to kube
// ID - the ID of the Pod we want to delete
// it returns an error is there was a problem
func DeletePod(kubeURL string, ID string) error {
	glog.Infoln("deleting pod " + ID)

	var caFile = "/kubekeys/root.crt"
	var certFile = "/kubekeys/cert.crt"
	var keyFile = "/kubekeys/key.key"

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
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
	client := &http.Client{Transport: transport}

	// DELETE pod
	var url = kubeURL + "/api/v1beta1/pods/" + ID
	glog.Infoln("url is " + url)
	request, err2 := http.NewRequest("DELETE", url, nil)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}

	resp, err := client.Do(request)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}
	glog.Infoln(string(data))

	// DELETE service 1 (port 13000)
	url = kubeURL + "/api/v1beta1/services/" + ID
	glog.Infoln("url is " + url)
	request, err2 = http.NewRequest("DELETE", url, nil)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}

	resp, err = client.Do(request)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 = ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}
	glog.Infoln(string(data))

	// DELETE service 2 (port 5432)
	url = kubeURL + "/api/v1beta1/services/" + ID + "-db"
	glog.Infoln("url is " + url)
	request, err2 = http.NewRequest("DELETE", url, nil)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}

	resp, err = client.Do(request)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 = ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return err2
	}
	glog.Infoln(string(data))

	return nil
}

// CreatePod creates a new pod and service using passed in values
// kubeURL - the URL to the kube
// podInfo - the params used to configure the pod
// return an error if anything goes wrong
func CreatePod(kubeURL string, podInfo template.KubePodParams) error {
	var caFile = "/kubekeys/root.crt"
	var certFile = "/kubekeys/cert.crt"
	var keyFile = "/kubekeys/key.key"

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
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
	client := &http.Client{Transport: transport}

	glog.Infoln("creating pod " + podInfo.ID)

	//use a pod template to build the pod definition
	data, err := template.KubeNodePod(podInfo)
	if err != nil {
		glog.Errorln("CreatePod:" + err.Error())
		return err
	}
	glog.Infoln(string(data[:]))

	var bodyType = "application/json"
	var url = kubeURL + "/api/v1beta1/pods"
	var serviceurl = kubeURL + "/api/v1beta1/services"
	glog.Infoln("url is " + url)
	glog.Infoln("serviceurl is " + serviceurl)

	// POST POD
	resp, err := client.Post(url, bodyType, bytes.NewReader(data))
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}
	defer resp.Body.Close()

	// Dump response
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}
	log.Println(string(data))

	//use a service template to build the service definition
	//service 1 is for the admin port 13000
	podInfo.PORT = "13000"
	var s1data []byte
	s1data, err = template.KubeNodeService(podInfo)
	if err != nil {
		glog.Errorln("CreatePod:" + err.Error())
		return err
	}
	glog.Infoln("service 1 request...")
	glog.Infoln(string(s1data[:]))

	// POST admin SERVICE at port 13000
	resp1, err1 := client.Post(serviceurl, bodyType, bytes.NewReader(s1data))
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}
	defer resp1.Body.Close()

	// Dump response
	data, err = ioutil.ReadAll(resp1.Body)
	if err != nil {
		glog.Errorln(err1.Error())
		return nil
	}
	log.Println("service 1 response..." + string(data))

	// POST pg SERVICE at port 5432 adding "-db" as suffix to name
	podInfo.PORT = "5432"
	podInfo.ID = podInfo.ID + "-db"
	var s2data []byte
	s2data, err = template.KubeNodeService(podInfo)
	if err != nil {
		glog.Errorln("CreatePod:" + err.Error())
		return err
	}

	glog.Infoln("service 2 request...")
	glog.Infoln(string(s2data[:]))
	resp2, err2 := client.Post(serviceurl, bodyType, bytes.NewReader(s2data))
	if err2 != nil {
		glog.Errorln(err2.Error())
		return nil
	}
	defer resp2.Body.Close()

	// Dump response
	data, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}
	log.Println("service 2 response ..." + string(data))

	return nil

}

// GetPods gets all the pods
// kubeURL - the URL to the kube
// podInfo - the params used to configure the pod
// return an error if anything goes wrong
func GetPods(kubeURL string, podInfo template.KubePodParams) error {
	var caFile = "/kubekeys/root.crt"
	var certFile = "/kubekeys/cert.crt"
	var keyFile = "/kubekeys/key.key"

	glog.Infoln("creating pod " + podInfo.ID)

	//use a pod template to build the pod definition
	data, err := template.KubeNodePod(podInfo)
	if err != nil {
		glog.Errorln("CreatePod:" + err.Error())
		return err
	}

	glog.Infoln(string(data[:]))

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorln(err.Error())
		return nil
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
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get(kubeURL + "/api/v1beta1/pods")
	if err != nil {
		glog.Errorln(err.Error())
		return nil
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return nil
	}
	log.Println(string(data))

	return nil
}

// GetPod gets information about a single pod from kube
// kubeURL - the URL to the kube
// podName - the pod name
// return an error if anything goes wrong
func GetPod(kubeURL string, podName string) (MyPod, error) {
	var podInfo MyPod
	var caFile = "/kubekeys/root.crt"
	var certFile = "/kubekeys/cert.crt"
	var keyFile = "/kubekeys/key.key"

	glog.Infoln("getting pod info " + podName)

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Errorln(err.Error())
		return podInfo, nil
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorln(err.Error())
		return podInfo, nil
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
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get(kubeURL + "/api/v1beta1/pods/" + podName)
	if err != nil {
		glog.Errorln(err.Error())
		return podInfo, nil
	}
	defer resp.Body.Close()

	// Dump response
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln(err2.Error())
		return podInfo, nil
	}
	log.Println(string(data))
	err2 = json.Unmarshal(data, &podInfo)
	if err2 != nil {
		glog.Errorln("error in unmarshalling pod " + err2.Error())
		return podInfo, err2
	}

	return podInfo, nil
}
