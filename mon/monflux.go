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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/myinfluxdb/client"
	"log"
)

//make sure to create the cpm database if it doesn't exist
func Bootdb() {
	logit.Info.Println("monflux:Bootdb: called ")

	//create a connection to influx but not to a database yet
	c, err := client.NewClient(&client.ClientConfig{
		Username: "root",
		Password: "root",
	})

	if err != nil {
		logit.Error.Println(err.Error())
		panic("cant connect to influxdb")
	}

	//verify that the cpm database exists
	dbs, err := c.GetDatabaseList()
	if err != nil {
		panic(err)
	}
	log.Printf("number of databases %d\n", len(dbs))
	var found = false
	for i := range dbs {
		for key, value := range dbs[i] {
			log.Printf("key:%s value:%s\n", key, value.(string))
			if value.(string) == "cpm" {
				found = true
			}
		}
	}
	if !found {
		logit.Info.Println("cpm database not found...creating it")

		if err := c.CreateDatabase("cpm"); err != nil {
			logit.Error.Println("cant create the cpm database")
			panic(err)
		}
	}

}
