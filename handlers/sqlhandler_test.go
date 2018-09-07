package handlers

import (
	"database/sql"
	"fmt"
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

var dbTest *sql.DB

const (
	jwtTokenOk                      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZSI6InBvc3RncmVzIn0.RiKyWr4Kw5TtFi9iGAkkqOYEtm284-2GNSt1oGHrTbg"
	jwtTokenExpire                  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoicG9zdGdyZXMiLCJleHAiOjE1MTYyMzkwMjJ9.5VyMx5na1V0K1EUBGyCqtkWvgD9Wu9Y95AYDUgwbg18"
	jwtTokenNoRole                  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.bqD_RIUSIOZlw38SKhrWqrM66dBXWGDAeF-IV62Qb0s"
	jwtTokenUnexpectedSigningMethod = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.Xbjh8W6mUQkhEWpMcVZ9PYKk_ezPn9eYHxlfm1_sEC9xCjUzhmCPh9BcHwTEGW5ThvfeMxI3Bdtj-NRAg3DMjA"
)

func setUp() {
	isFirst := false
	if dbTest == nil {
		isFirst = true
		os.Setenv("GOAPISQL_ENV", "test")
		config.Init()
		dbTest = config.GetDbConnection("postgres")
	}

	if isFirst {
		dropTable("customer")
		createTableCustomer()
	} else {
		dbTest.Exec("DELETE FROM customer;")
		dbTest.Exec("ALTER SEQUENCE customer_id_seq RESTART WITH 1;")
	}
}

func TestSqlHandlerJwtTokenOk(t *testing.T) {
	setUp()
	dbTest.Exec("INSERT INTO customer (age, first_name, last_name, dimension) VALUES (43, 'Konstantin', 'Savenkov', 15.3)")
	dbTest.Exec("INSERT INTO customer (age, first_name, last_name, dimension) VALUES (35, 'Oksana', 'Savenkova', 12.5)")
	port := 1234
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenOk)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 200 {
		t.Error("The response should be", 200)
	}
	val := "Konstantin"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
	val = "Oksana"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerJwtTokenExpire(t *testing.T) {
	setUp()
	port := 1235
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenExpire)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 403 {
		t.Error("The response should be", 403)
	}
	val := "Token is expired"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerJwtTokenNoRole(t *testing.T) {
	setUp()
	port := 1236
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenNoRole)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 403 {
		t.Error("The response should be", 403)
	}
	val := "No role in Authorization token"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerJwtTokenUnexpectedSigningMethod(t *testing.T) {
	setUp()
	port := 1237
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenUnexpectedSigningMethod)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 403 {
		t.Error("The response should be", 403)
	}
	val := "No role in Authorization token"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerJwtTokenShouldBePrefixBearer(t *testing.T) {
	setUp()
	port := 1238
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", jwtTokenOk)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 403 {
		t.Error("The response should be", 403)
	}
	val := "Error: Error in authorization header: should be prefix 'Bearer '"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerJwtTokenError(t *testing.T) {
	setUp()
	port := 1239
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "SELECT * FROM customer"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+"error token")
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 403 {
		t.Error("The response should be", 403)
	}
	val := "Error: Error in authorization token"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerSQLError(t *testing.T) {
	setUp()
	port := 1240
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "ERROR SQL STRING"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenOk)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 400 {
		t.Error("The response should be", 400)
	}
	val := "pq: syntax error at or near"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func TestSqlHandlerErrorOnlyPost(t *testing.T) {
	setUp()
	port := 1241
	defer startServerOnPort(t, port, SQLHandler).Close()

	sqlSting := "ERROR SQL STRING"
	body := strings.NewReader(sqlSting)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", "http://localhost:"+strconv.Itoa(port), body)
	req.Header.Set("Authorization", "Bearer "+jwtTokenOk)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	responseString := string(responseBody)

	if resp.StatusCode != 400 {
		t.Error("The response should be", 400)
	}
	val := "Only Post method allowed"
	if !strings.Contains(responseString, val) {
		t.Error("The string should to contains", val)
	}
}

func startServerOnPort(t *testing.T, port int, h fasthttp.RequestHandler) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %d: %s", port, err)
	}
	go fasthttp.Serve(ln, h)
	return ln
}

func dropTable(tableName string) {
	dbTest := config.GetDbConnection("postgres")
	sqlStatement := fmt.Sprintf(`DROP TABLE %s;`, tableName)
	_, err := dbTest.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

func createTableCustomer() {
	dbTest := config.GetDbConnection("postgres")
	sqlStatement := `
CREATE TABLE customer (
  id SERIAL PRIMARY KEY,
  age INT,
  dimension real,
  first_name TEXT,
  last_name TEXT
);
`
	_, err := dbTest.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
