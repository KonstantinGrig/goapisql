package config

import (
	"fmt"
	"os"
	"testing"
)

type testpair struct {
	key   string
	value string
}

var tests = []testpair{
	{"db-uri-template", "postgres://%s:%s@localhost:5433/goapisql_db_test?sslmode=disable"},
	{"db-anon-role", "web_anon"},
	{"jwt-secret", "G0naHHmCgbbfgfgUUUbvlaztyrOpppQ"},
}

func BenchmarkGetEnv(b *testing.B) {
	os.Setenv("GOAPISQL_ENV", "test")
	for i := 0; i < b.N; i++ {
		GetEnv()
	}
}

func TestGetNoCachedString(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")

	for _, pair := range tests {
		actual := GetNoCachedString(pair.key)
		if actual != pair.value {
			strErr := fmt.Sprintf("Expected: %s, got ", pair.value)
			t.Error(strErr, actual)
		}
	}
}

func TestGetString(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")

	for _, pair := range tests {
		actual := GetString(pair.key)
		actual = GetString(pair.key)
		actual = GetString(pair.key)
		actual = GetString(pair.key)
		if actual != pair.value {
			strErr := fmt.Sprintf("Expected: %s, got ", pair.value)
			t.Error(strErr, actual)
		}
	}
}

func TestGetNoCachedDbUsers(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")
	GetNoCachedDbUsers()
}

func TestGetDbUsers(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")
	GetNoCachedDbUsers()
	GetDbUsers()
	GetDbUsers()
	GetDbUsers()
	GetDbUsers()
}

func TestGetNoCachedDbConnection(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")
	GetNoCachedDbConnection("postgres")
}

func TestGetDbConnection(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")
	GetDbConnection("postgres")
	GetDbConnection("postgres")
	GetDbConnection("postgres")
	GetDbConnection("postgres")
}

func TestInitConfigFile(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	InitConfigFile("../config.json")
	str := GetNoCachedString("server-port")
	t.Log(str)
}
