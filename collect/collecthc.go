package collect

import (
	"database/sql"
	"github.com/crunchydata/crunchy-postgresql-manager/adminapi"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
)

const CLUSTERADMIN_DB = "clusteradmin"

func Collecthc() error {
	var err error

	logit.Info.Println("Collecthc called")

	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
	}
	defer dbConn.Close()


	//get all containers
	var containers []admindb.Container
	containers, err = admindb.GetAllContainers(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	//for each container, do a health check
	i := 0
	var credential adminapi.Credential
	var checks []admindb.HealthCheck
	checks = make([]admindb.HealthCheck, 0)

	var status string

	for i = range containers {
		hc := admindb.HealthCheck{}
		hc.ProjectID = containers[i].ProjectID
		hc.ProjectName = containers[i].ProjectName
		hc.ContainerName = containers[i].Name
		hc.ContainerID = containers[i].ID
		hc.ContainerRole = containers[i].Role
		hc.ContainerImage = containers[i].Image

		credential, err = adminapi.GetUserCredentials(dbConn, &containers[i])
                if err != nil {
                        logit.Error.Println(err.Error())
                } else {


			status, err = ping(&credential)
			hc.Status = status

			checks = append(checks, hc)
		}
		i++
	}

	//delete current health checks
	err = admindb.DeleteHealthCheck(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//persist health checks
	i = 0
	for i = range checks {
		_, err = admindb.InsertHealthCheck(dbConn, checks[i])
		if err != nil {
			logit.Error.Println(err.Error())
		}
		i++
	}

	return nil
}

func ping(credential *adminapi.Credential) (string, error) {
	var db *sql.DB
	var err error

	db, err = util.GetMonitoringConnection(credential.Host,
		credential.Username, credential.Port, credential.Database, 
		credential.Password)
	defer db.Close()
	if err != nil {
		logit.Error.Println("error in getting connectionto " + credential.Host)
		return "down", err
	}

	var result string
	err = db.QueryRow("select now()::text").Scan(&result)
	if err != nil {
		logit.Error.Println("could not ping db on " + credential.Host)
		return "down", err
	}
	return "up", nil

}
