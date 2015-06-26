package collect

import (
	"github.com/crunchydata/crunchy-postgresql-manager/sec"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	_ "github.com/lib/pq"
	"database/sql"
)

func Collecthc() (error) {
	var err error

	logit.Info.Println("Collecthc called")

  	var dbConn *sql.DB
       	dbConn, err = util.GetConnection("clusteradmin")
       	if err != nil {
       		logit.Error.Println(err.Error())
	}
	defer dbConn.Close()

	var domain string
	domain, err = getDomain(dbConn)
	var pgport string
	pgport, err = getPort(dbConn)

	//get all containers
	var containers []admindb.Container
	admindb.SetConnection(dbConn)
	containers, err = admindb.GetAllContainers()
	if err != nil {
       		logit.Error.Println(err.Error())
	}

	//for each container, do a health check
	i := 0
	var checks []HealthCheck
	checks = make([]HealthCheck,0)

	var status string

	for i = range containers {
		hc := HealthCheck{}
		hc.ProjectID = containers[i].ProjectID
		hc.ProjectName = containers[i].ProjectName
		hc.ContainerName = containers[i].Name
		hc.ContainerID = containers[i].ID
		hc.ContainerRole = containers[i].Role
		hc.ContainerImage = containers[i].Image

		status, err = ping(dbConn, pgport, hc.ContainerName + "." + domain, hc.ContainerRole)
		hc.Status = status

		checks = append(checks, hc)
		i++
	}

	//delete current health checks
	err = DeleteHealthCheck(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//persist health checks
	i = 0
	for i = range checks {
		_, err = InsertHealthCheck(dbConn, checks[i])
		if err != nil {
			logit.Error.Println(err.Error())
		}
		i++
	}

	return nil
}

func ping(dbConn *sql.DB, port string, containerName string, containerRole string) (string, error) {
	var db *sql.DB
	var err error
	var userid, password, database string

	//get node credentials
	userid, password, database, err = getCredential(dbConn, containerName, containerRole)
	if err != nil {
		logit.Error.Println(err.Error())
		return "down", err
	}	

	db, err = util.GetMonitoringConnection(containerName,
		userid, port, database, password)
	defer db.Close()
	if err != nil {
		logit.Error.Println("error in getting connectionto " + containerName)
		return "down", err
	} 

	var result string
	err = db.QueryRow("select now()::text").Scan(&result)
	if err != nil {
		logit.Error.Println("could not ping db on " + containerName)
		return "down", err
	}
	return "up", nil	

}

func getDomain(dbConn *sql.DB) (string, error) {
	var err error
	var domain string

	admindb.SetConnection(dbConn)
	domain, err = admindb.GetDomain()
	if err != nil {
		logit.Error.Println(err.Error())
		return domain, err
	}
	return domain, nil
}

func getPort(dbConn *sql.DB) (string, error) {
	var err error
	var port admindb.Setting

	admindb.SetConnection(dbConn)
	port, err = admindb.GetSetting("PG-PORT")
	if err != nil {
		logit.Error.Println(err.Error())
		return port.Value, err
	}
	return port.Value, nil
}

//return id, password, database
func getCredential(dbConn *sql.DB, containerName string, containerRole string) (string, string, string, error) {
	var err error
	var userID  string
	var password string
	var database string

	admindb.SetConnection(dbConn)
	if containerRole == "pgpool" {
		userID = "cpmtest"
		database = "cpmtest"
	} else {
		userID = "postgres"
		database = "postgres"
	}

	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(containerName, userID)
	if err != nil {
		logit.Error.Println(err.Error())
		return userID, password, database, err
	}
	password, err = sec.DecryptPassword(nodeuser.Passwd)
	if err != nil {
		logit.Error.Println(err.Error())
		return userID, password, database, err
	}

	return userID, password, database, nil
}
