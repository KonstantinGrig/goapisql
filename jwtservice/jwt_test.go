package jwtservice

import (
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	config.InitConfigFile("../config.json")
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

func TestCreateToken(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"role": "publisher",
		"exp":  time.Now().Add(1 * time.Second).Unix(),
	}
	tokenString, err := CreateToken(mapClaims)
	if err != nil {
		t.Error("CreateToken error", err)
	}
	res, err := Parse(config.PREFIX_TOKEN + tokenString)
	if err != nil {
		t.Error("Parse error", err)
	}
	if res["role"] != "publisher" {
		t.Error("Wrong token")
	}
}
