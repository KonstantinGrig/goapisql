package handlers

import (
	"errors"
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/KonstantinGrig/goapisql/goapisql"
	"github.com/KonstantinGrig/goapisql/jwtservice"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
)

//FastHttpHandler the function processes the SQL request coming from the body of the request post
func FastHttpHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != "POST" {
		errorHandler(ctx, errors.New("Only Post method allowed"), 400)
		return
	}
	authorizationHeader := string(ctx.Request.Header.Peek("Authorization"))
	sqlString := string(ctx.Request.Body())
	res, statusCode, err := ProcSqlRequest(authorizationHeader, sqlString)
	if err != nil {
		errorHandler(ctx, err, statusCode)
		return
	}
	_, err = ctx.Write(res)
	if err != nil {
		errorHandler(ctx, err, 400)
		return
	}
}

//NetHttpHandler the function processes the SQL request coming from the body of the request post
func NetHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		w.Write([]byte("Only Post method allowed"))
		return
	}
	authorizationHeader := string(r.Header.Get("Authorization"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	sqlString := string(body)
	res, statusCode, err := ProcSqlRequest(authorizationHeader, sqlString)
	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(res)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
}

func ProcSqlRequest(authorizationHeader string, sqlString string) ([]byte, int, error) {
	claims, err := jwtservice.Parse(authorizationHeader)
	if err != nil {
		return nil, 403, err
	}
	roleClaims := claims["role"]
	if roleClaims == nil {
		return nil, 403, errors.New("No role in Authorization token")
	}
	role := roleClaims.(string)
	db := config.GetDbConnection(role)
	res, err := goapisql.GetQueryResult(db, sqlString)
	if err != nil {
		return nil, 400, err
	}
	return res, 200, err
}

func errorHandler(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.Error("Error: "+err.Error(), statusCode)
}
