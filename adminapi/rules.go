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

package adminapi

import (
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

const CLUSTERADMIN_DB = "clusteradmin"

type Rule struct {
	ID          string
	Token       string
	Name        string
	Type        string
	Database    string
	User        string
	Address     string
	Method      string
	Description string
	CreateDate  string
	UpdateDate  string
}

type ContainerAccessRule struct {
	ID             string
	Token          string
	ContainerID    string
	AccessRuleID   string
	AccessRuleName string
	Selected       string
	CreateDate     string
}

func RulesGet(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("RulesGet: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	rule, err := GetAccessRule(ID)
	if err != nil {
		logit.Error.Println("RulesGet:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&rule)
}

func RulesGetAll(w rest.ResponseWriter, r *rest.Request) {
	var err error
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("RulesGetAll: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rules, err := GetAllRules()
	if err != nil {
		logit.Error.Println("RulesGetAll:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&rules)
}

func RulesDelete(w rest.ResponseWriter, r *rest.Request) {
	var err error
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("RulesDelete: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	err = DeleteRule(ID)
	if err != nil {
		logit.Error.Println("RulesGetAll:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func RulesUpdate(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("RulesUpdate: in RulesUpdate")
	rule := Rule{}
	err := r.DecodeJsonPayload(&rule)
	if err != nil {
		logit.Error.Println("RulesUpdate: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(rule.Token, "perm-container")
	if err != nil {
		logit.Error.Println("RulesUpdate: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if rule.ID == "" {
		logit.Error.Println("RulesUpdate: error in ID")
		rest.Error(w, "rule ID required", http.StatusBadRequest)
		return
	}

	if rule.Name == "" {
		logit.Error.Println("RulesUpdate: error in Name")
		rest.Error(w, "rule name required", http.StatusBadRequest)
		return
	}

	err = UpdateRule(rule)
	if err != nil {
		logit.Error.Println("RulesUpdate: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}
func RulesInsert(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("RulesInsert: in RulesInsert")
	rule := Rule{}
	err := r.DecodeJsonPayload(&rule)
	if err != nil {
		logit.Error.Println("RulesInsert: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(rule.Token, "perm-container")
	if err != nil {
		logit.Error.Println("RulesInsert: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if rule.Name == "" {
		logit.Error.Println("RulesInsert: error in Name")
		rest.Error(w, "rule name required", http.StatusBadRequest)
		return
	}

	err = InsertRule(rule)
	if err != nil {
		logit.Error.Println("RulesUpdate: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func UpdateRule(rule Rule) error {
	queryStr := fmt.Sprintf(
		"update accessrule set ( name, ruletype, database, ruleuser, address, method, description, updatedt) = ('%s', '%s', '%s', '%s', '%s', '%s', '%s', now()) where id = %s returning id",
		rule.Name,
		rule.Type,
		rule.Database,
		rule.User,
		rule.Address,
		rule.Method,
		rule.Description,
		rule.ID)
	logit.Info.Println(queryStr)

	var ruleid int
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	err = dbConn.QueryRow(queryStr).Scan(&ruleid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("rule updated " + rule.Name)
	}
	return nil

}

func InsertRule(rule Rule) error {
	queryStr := fmt.Sprintf(
		"insert into accessrule ( name, ruletype, database, ruleuser, address, method, description, createdt, updatedt) values ( '%s', '%s', '%s', '%s', '%s', '%s', '%s', now(), now()) returning id",
		rule.Name,
		rule.Type,
		rule.Database,
		rule.User,
		rule.Address,
		rule.Method,
		rule.Description)

	logit.Info.Println(queryStr)
	var ruleid int
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	err = dbConn.QueryRow(queryStr).Scan(&ruleid)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("accessrule inserted id " + strconv.Itoa(ruleid))
	}

	return nil

}

func DeleteRule(ID string) error {
	queryStr := fmt.Sprintf("delete from accessrule where  id=%s returning id", ID)
	logit.Info.Println(queryStr)

	var ruleid int
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	err = dbConn.QueryRow(queryStr).Scan(&ruleid)
	switch {
	case err != nil:
		logit.Error.Println(err)
		return err
	default:
		logit.Info.Println("deleted  accessrule id " + ID)
	}
	return nil

}

func GetAccessRule(ID string) (Rule, error) {
	rule := Rule{}

	queryStr := fmt.Sprintf("select ID, NAME, RULETYPE, DATABASE, RULEUSER, ADDRESS, METHOD, DESCRIPTION, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS'),  to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from accessrule where id = %s", ID)

	logit.Info.Println(queryStr)

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return rule, err
	}
	err = dbConn.QueryRow(queryStr).Scan(
		&rule.ID,
		&rule.Type,
		&rule.Database,
		&rule.User,
		&rule.Address,
		&rule.Method,
		&rule.Description,
		&rule.CreateDate,
		&rule.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("no accessrule found id " + ID)
		return rule, err
	case err != nil:
		return rule, err
	}

	return rule, nil

}

func GetAllRules() ([]Rule, error) {

	var rules []Rule
	var rows *sql.Rows
	var err error
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return rules, err
	}
	rows, err = dbConn.Query(
		"select id, name, ruletype, database, ruleuser, address, method, description, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS'), to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from accessrule order by name")
	if err != nil {
		return rules, err
	}
	defer rows.Close()
	rules = make([]Rule, 0)
	for rows.Next() {
		rule := Rule{}
		if err = rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Type,
			&rule.Database,
			&rule.User,
			&rule.Address,
			&rule.Method,
			&rule.Description,
			&rule.CreateDate,
			&rule.UpdateDate,
		); err != nil {
			return rules, err
		}
		rules = append(rules, rule)
	}
	if err = rows.Err(); err != nil {
		return rules, err
	}
	return rules, nil
}

//
// containeraccessrules logic
//

/*
func ContainerAccessRuleGet(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerRulesGet: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	car, err := GetContainerAccessRule(ID)
	if err != nil {
		logit.Error.Println("ContainerRulesGet:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&car)
}
*/

func ContainerAccessRuleGetAll(w rest.ResponseWriter, r *rest.Request) {
	var err error
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerRulesGetAll: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ContainerID := r.PathParam("ID")
	if ContainerID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	cars, err := GetAllContainerAccessRule(ContainerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&cars)
}

/*
func ContainerAccessRuleDelete(w rest.ResponseWriter, r *rest.Request) {
	var err error
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerRuleDelete: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	containerRuleID := r.PathParam("ID")
	if containerRuleID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	err = DeleteContainerAccessRule(containerRuleID)
	if err != nil {
		logit.Error.Println("ContainerRuleDelete:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func ContainerAccessRuleInsert(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("ContainerAccessRuleInsert")
	car := ContainerAccessRule{}
	err := r.DecodeJsonPayload(&car)
	if err != nil {
		logit.Error.Println("ContainerAccessRuleInsert: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(car.Token, "perm-container")
	if err != nil {
		logit.Error.Println("ContainerAccessRuleInsert: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = InsertContainerAccessRule(car)
	if err != nil {
		logit.Error.Println("ContainerAccessRuleUpdate: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}
*/
func ContainerAccessRuleUpdate(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("ContainerAccessRuleInsert")
	cars := make([]ContainerAccessRule, 0)
	err := r.DecodeJsonPayload(&cars)
	if err != nil {
		logit.Error.Println("ContainerAccessRuleInsert: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(cars) > 0 {
		err = secimpl.Authorize(cars[0].Token, "perm-container")
		if err != nil {
			logit.Error.Println("ContainerAccessRuleInsert: authorize error " + err.Error())
			rest.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	for i := range cars {
		if cars[i].Selected == "true" {
			if cars[i].ID != "-1" {
				//do nothing
				logit.Info.Println("doing nothing on " + cars[i].ContainerID + " " + cars[i].AccessRuleID + " " + cars[i].AccessRuleName)
			} else {
				//insert
				logit.Info.Println("doing insert on " + cars[i].ContainerID + " " + cars[i].AccessRuleID + " " + cars[i].AccessRuleName)
				err = InsertContainerAccessRule(cars[i])
				if err != nil {
					logit.Error.Println("error doing insert on " + cars[i].ContainerID + " " + cars[i].AccessRuleID + " " + cars[i].AccessRuleName)
					logit.Error.Println(err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
			}
		} else {
			if cars[i].ID != "-1" {
				//delete
				logit.Info.Println("doing delete on " + cars[i].ContainerID + " " + cars[i].AccessRuleID + " " + cars[i].AccessRuleName)
				err = DeleteContainerAccessRule(cars[i].ID)
				if err != nil {
					logit.Error.Println("error doing delete on " + cars[i].ContainerID + " " + cars[i].AccessRuleID + " " + cars[i].AccessRuleName)
					logit.Error.Println(err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
			}
		}
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

//
// database logic for containeraccessrules
//

func DeleteContainerAccessRule(ID string) error {
	queryStr := fmt.Sprintf("delete from containeraccessrule where id=%s returning id", ID)
	logit.Info.Println(queryStr)

	var id int
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	err = dbConn.QueryRow(queryStr).Scan(&id)
	switch {
	case err != nil:
		logit.Error.Println(err)
		return err
	default:
		logit.Info.Println("deleted containeraccessrule id " + ID)
	}
	return nil

}

func InsertContainerAccessRule(car ContainerAccessRule) error {
	queryStr := fmt.Sprintf(
		"insert into containeraccessrule ( containerid, accessruleid, createdt ) values ( %s, %s, now()) returning id",
		car.ContainerID,
		car.AccessRuleID)

	logit.Info.Println(queryStr)
	var id int
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	err = dbConn.QueryRow(queryStr).Scan(&id)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("containeraccessrule inserted id " + strconv.Itoa(id))
	}

	return nil

}

func GetAllContainerAccessRule(containerID string) ([]ContainerAccessRule, error) {

	var cars []ContainerAccessRule
	var rows *sql.Rows
	var err error
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return cars, err
	}
	rows, err = dbConn.Query(
		"select coalesce(car.id, -1), coalesce(car.containerid, -1), ar.id, ar.name, to_char(coalesce(ar.createdt, now()), 'MM-DD-YYYY HH24:MI:SS') from accessrule ar left outer join containeraccessrule car on (car.accessruleid = ar.id and car.containerid = " + containerID + ")")
	if err != nil {
		return cars, err
	}
	defer rows.Close()
	cars = make([]ContainerAccessRule, 0)
	for rows.Next() {
		car := ContainerAccessRule{}
		if err = rows.Scan(
			&car.ID,
			&car.ContainerID,
			&car.AccessRuleID,
			&car.AccessRuleName,
			&car.CreateDate,
		); err != nil {
			return cars, err
		}
		if car.ID == "-1" {
			car.Selected = "false"
		} else {
			car.Selected = "true"
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return cars, err
	}
	return cars, nil
}

func GetContainerAccessRule(ID string) (ContainerAccessRule, error) {
	car := ContainerAccessRule{}

	queryStr := fmt.Sprintf("select ID, CONTAINERID, ACCESSRULEID, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from containeraccessrule where id = %s", ID)

	logit.Info.Println(queryStr)

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	defer dbConn.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return car, err
	}
	err = dbConn.QueryRow(queryStr).Scan(
		&car.ID,
		&car.ContainerID,
		&car.AccessRuleID,
		&car.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("no containeraccessrule found id " + ID)
		return car, err
	case err != nil:
		return car, err
	}

	return car, nil

}
