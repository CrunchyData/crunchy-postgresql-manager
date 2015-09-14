package util

import (
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
	"math/rand"
)

type Customer struct {
	ID       string
	Name     string
	Location string
}

type Product struct {
	ID          string
	CustomerID  string
	ProductName string
	ProductDesc string
}

func main() {

	fmt.Println("loading...")

	dbConn, err := util.GetMonitoringConnection("127.0.0.1", "postgres", "5432", "postgres", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	customer := Customer{}
	product := Product{}

	var MAX_ROWS = 10

	for i := 0; i < MAX_ROWS; i++ {
		customer.Name = randSeq(20)
		customer.Location = randSeq(100)
		queryStr := fmt.Sprintf("insert into customer ( name, location) values ( '%s', '%s') returning id", customer.Name, customer.Location)

		var customerid int
		err = dbConn.QueryRow(queryStr).Scan(&customerid)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		default:
			//fmt.Println("admindb:InsertCluster: cluster inserted returned is " + strconv.Itoa(clusterid))
		}

		product.ProductName = randSeq(20)
		product.ProductDesc = randSeq(100)
		queryStr = fmt.Sprintf("insert into product ( customerid, productname, productdesc) values ( %d, '%s', '%s') returning id", customerid, product.ProductName, product.ProductDesc)

		var productid int
		err = dbConn.QueryRow(queryStr).Scan(&productid)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		default:
			//fmt.Println("admindb:InsertCluster: cluster inserted returned is " + strconv.Itoa(clusterid))
		}
	}

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
