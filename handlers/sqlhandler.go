package handlers

import (
	"fmt"
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/KonstantinGrig/goapisql/goapisql"
	"github.com/KonstantinGrig/goapisql/jwtservice"

	//"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

//SQLHandler the function processes the SQL request coming from the body of the request post
func SQLHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) == "POST" {
		authorizationHeader := string(ctx.Request.Header.Peek("Authorization"))
		claims, err := jwtservice.Parse(authorizationHeader)
		if err != nil {
			ctx.Response.SetStatusCode(403)
			fmt.Fprintf(ctx, "Error: %s", err)
			return
		}
		roleClaims := claims["role"]
		if roleClaims == nil {
			ctx.Response.SetStatusCode(403)
			fmt.Fprintf(ctx, "Error: %s", "No role in Authorization token")
			return
		}
		role := roleClaims.(string)

		db := config.GetDbConnection(role)
		sqlString := string(ctx.Request.Body())

		res, err := goapisql.GetQueryResult(db, sqlString)
		if err != nil {
			ctx.Response.SetStatusCode(400)
			fmt.Fprintf(ctx, "Error: %s", err)
			return
		}
		ctx.WriteString(res)
	} else {
		ctx.Response.SetStatusCode(400)
		ctx.WriteString("Only Post method allowed")
	}
}
