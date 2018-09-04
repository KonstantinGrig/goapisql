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
	{"db-default-role-password", "hu8jmn3"},
	{"server-port", "3000"},
	{"jwtservice-secret", "G0naHHmCgbbLPoK4rdTnGhl30B3VKBkD"},
}

func BenchmarkGetEnv(b *testing.B) {
	os.Setenv("GOAPISQL_ENV", "test")
	for i := 0; i < b.N; i++ {
		GetEnv()
	}
}

func TestGetNoCachedString(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	Init()

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
	Init()

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
	Init()
	GetNoCachedDbUsers()
}

func TestGetDbUsers(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	Init()
	GetNoCachedDbUsers()
	GetDbUsers()
	GetDbUsers()
	GetDbUsers()
	GetDbUsers()
}

func TestGetNoCachedDbConnection(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	Init()
	GetNoCachedDbConnection("postgres")
}

func TestGetDbConnection(t *testing.T) {
	os.Setenv("GOAPISQL_ENV", "test")
	Init()
	GetDbConnection("postgres")
	GetDbConnection("postgres")
	GetDbConnection("postgres")
	GetDbConnection("postgres")
}
