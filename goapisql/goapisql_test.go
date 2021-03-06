package goapisql

import (
	"database/sql"
	"fmt"
	"github.com/KonstantinGrig/goapisql/config"
	"log"
	"os"
	"strings"
	"testing"
)

var dbTest *sql.DB
var dbRole string

func setUp() {
	isFirst := false
	if dbTest == nil {
		log.Println("dbTest == nil")
		isFirst = true
		os.Setenv("GOAPISQL_ENV", "test")
		config.InitConfigFile("../config.json")
		dbRole = "postgres"
		dbTest = config.GetDbConnection(dbRole)
	} else {
		log.Println("dbTest not nil")
	}

	if isFirst {
		dropTable("customer")
		createTableCustomer()
	} else {
		dbTest.Exec("DELETE FROM customer;")
		dbTest.Exec("ALTER SEQUENCE customer_id_seq RESTART WITH 1;")
	}
}

func TestGetQueryResultSelect(t *testing.T) {
	setUp()
	dbTest.Exec("INSERT INTO customer (age, first_name, last_name, dimension) VALUES (43, 'Konstantin', 'Savenkov', 15.3)")
	dbTest.Exec("INSERT INTO customer (age, first_name, last_name, dimension) VALUES (35, 'Oksana', 'Savenkova', 12.5)")

	sqlString := "SELECT * FROM customer"
	res, err := GetQueryResult(dbRole, sqlString)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(res))

	val := "Konstantin"
	if !strings.Contains(string(res), val) {
		t.Error("The string should to contains", val)
	}
	val = "Oksana"
	if !strings.Contains(string(res), val) {
		t.Error("The string should to contains", val)
	}
	val = "user_role"
	if !strings.Contains(string(res), val) {
		t.Error("The string should to contains", val)
	}
}

func TestGetQueryResultInsert(t *testing.T) {
	setUp()

	sqlString := "INSERT INTO customer (age, first_name, last_name, dimension) VALUES (4, 'Kira', 'Mironova', 5.8) RETURNING *"
	res, err := GetQueryResult(dbRole, sqlString)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res)
	val := "Mironova"
	if !strings.Contains(string(res), val) {
		t.Error("The string should to contains", val)
	}
}

func TestGetQueryResult2(t *testing.T) {
	setUp()
	sqlString := ""
	res, err := GetQueryResult(dbRole, sqlString)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(res)
}

func TestGetQueryResult3(t *testing.T) {
	setUp()
	sqlString := ""
	res, err := GetQueryResult(dbRole, sqlString)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(res)
}

func dropTable(tableName string) {
	dbTest := config.GetDbConnection("postgres")
	sqlStatement := fmt.Sprintf(`DROP TABLE IF EXISTS %s; DROP TYPE IF EXISTS role_db_enum;`, tableName)
	_, err := dbTest.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

func createTableCustomer() {
	dbTest := config.GetDbConnection("postgres")
	sqlStatement := `
CREATE TYPE role_db_enum AS ENUM ('admin_role', 'manager_role', 'user_role');
CREATE TABLE customer (
  id SERIAL PRIMARY KEY,
  age INT,
  dimension real,
  first_name TEXT,
  last_name TEXT,
  role_db role_db_enum DEFAULT 'user_role'
);
`
	_, err := dbTest.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
