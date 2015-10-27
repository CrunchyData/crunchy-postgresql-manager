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

package admindb

import (
	"database/sql"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	_ "github.com/lib/pq"
	"strconv"
)

// GetHealthCheck return the current set of healthcheck metrics that are persisted
func GetHealthCheck(dbConn *sql.DB) ([]types.HealthCheck, error) {
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query(
		"select ID, ProjectName, ProjectID, ContainerName, ContainerID, " +
			"ContainerRole, ContainerImage, Status, to_char(UpdateDt, 'MM-DD-YYYY HH24:MI:SS') " +
			"from healthcheck order by ProjectName, ContainerName")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var checks []types.HealthCheck
	checks = make([]types.HealthCheck, 0)
	for rows.Next() {
		check := types.HealthCheck{}
		if err = rows.Scan(
			&check.ID,
			&check.ProjectName,
			&check.ProjectID,
			&check.ContainerName,
			&check.ContainerID,
			&check.ContainerRole,
			&check.ContainerImage,
			&check.Status, &check.UpdateDate); err != nil {
			return nil, err
		}

		checks = append(checks, check)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return checks, nil
}

// InsertHealthCheck insert a health check metric into the database, return the key
func InsertHealthCheck(dbConn *sql.DB, hc types.HealthCheck) (int, error) {
	queryStr := fmt.Sprintf(
		"insert into healthcheck ( "+
			"ProjectName, ProjectID, ContainerName, ContainerID, "+
			"ContainerRole, ContainerImage, Status, UpdateDt) values ("+
			"'%s', %s, '%s', %s, "+
			"'%s', '%s', '%s', now()) returning ID",
		hc.ProjectName, hc.ProjectID, hc.ContainerName,
		hc.ContainerID, hc.ContainerRole, hc.ContainerImage, hc.Status)

	logit.Info.Println("admindb:InsertHealthCheck:" + queryStr)
	var id int
	err := dbConn.QueryRow(queryStr).Scan(&id)
	switch {
	case err != nil:
		logit.Info.Println("InsertHealthCheck:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("InsertHealthCheck:inserted returned is " + strconv.Itoa(id))
	}

	return id, nil
}

// DeleteHealthCheck deleles the entire set of health check metrics that are persisted
func DeleteHealthCheck(dbConn *sql.DB) error {
	queryStr := fmt.Sprintf("delete from healthcheck")
	logit.Info.Println(queryStr)

	_, err := dbConn.Query(queryStr)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("DeleteHealthCheck:deleted ")
	return nil
}
