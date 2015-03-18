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
	"crunchy.com/admindb"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
	"net/http"
	"strconv"
)

func GetHC1(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetHC1: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var domain string
	domain, err = admindb.GetDomain()

	c, err := client.NewClient(&client.ClientConfig{
		Host:     "cpm-mon." + domain + ":8086",
		Username: "root",
		Password: "root",
		Database: "cpm",
	})

	if err != nil {
		glog.Errorln("GetHC1: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []*client.Series

	//get the latest HC1 record and it's seconds value
	var query = "select seconds, service, servicetype, status from hc1 limit 1"
	glog.Infoln(query)

	results, err = c.Query(query)
	if err != nil {
		glog.Errorln(err.Error())
		w.WriteJson(&results)
		//rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(results) == 0 {
		glog.Infoln("GetHC1: no results yet")
		w.WriteJson(&results)
		return
	}

	var resultsLen = len(results[0].Points)
	if resultsLen == 0 {
		glog.Infoln("GetHC1: no results yet 2")
		w.WriteJson(&results)
		return
	}

	glog.Infof("results len = %d\n", resultsLen)
	var seconds = results[0].Points[0][2].(float64)
	glog.Infof("results seconds=%f\n", seconds)

	query = "select seconds, service, servicetype, status from hc1 where seconds = " + strconv.FormatFloat(seconds, 'f', 2, 64)
	glog.Infoln(query)

	results, err = c.Query(query)

	resultsLen = len(results[0].Points)
	if resultsLen == 0 {
		glog.Infoln("GetHC1 b: no results")
	}

	w.WriteJson(&results)

}
