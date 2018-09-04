package jwtservice

import (
	"github.com/KonstantinGrig/goapisql/config"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	config.Init()
	authorizationHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZSI6InBvc3RncmVzIn0.RiKyWr4Kw5TtFi9iGAkkqOYEtm284-2GNSt1oGHrTbg"
	res, err := Parse(authorizationHeader)
	if err != nil {
		t.Error("Parse error", err)
	}

	val := "postgres"
	roleClaims := res["role"].(string)
	if !strings.Contains(roleClaims, val) {
		t.Error("The claims should to contains", val)
	}
}
