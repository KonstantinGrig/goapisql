package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
	"time"
)

var once sync.Once
var GoapisqlCache *cache.Cache

const (
	configName       = "config"
	configType       = "json"
	keyDbUsers       = "db_users"
	keyDbUriTemplate = "db-uri-template"
	keyPassword      = "password"
	dbDriverName     = "postgres"
)

func Init() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetCache() *cache.Cache {
	once.Do(func() {
		GoapisqlCache = cache.New(5*time.Minute, 10*time.Minute)
	})

	isFlushCache := os.Getenv("GOAPISQL_IS_CACHE_FLUSH")
	if isFlushCache == "yes" {
		GoapisqlCache.Flush()
		os.Setenv("GOAPISQL_IS_CACHE_FLUSH", "no")
	}

	return GoapisqlCache
}

func GetEnv() string {
	envVal := os.Getenv("GOAPISQL_ENV")
	if envVal == "" {
		envVal = "prod"
	}
	log.Println("Enviroment:", envVal)
	return envVal
}

func GetNoCachedString(key string) string {
	return viper.GetString(GetEnv() + "." + key)
}

func GetString(key string) string {
	var res string
	cacheApp := GetCache()
	resFromCache, ok := cacheApp.Get(key)
	if ok {
		res = resFromCache.(string)
	} else {
		res = GetNoCachedString(key)
		cacheApp.Set(key, res, cache.NoExpiration)
	}

	return res
}

func GetNoCachedDbUsers() map[string]map[string]string {
	resSlice := map[string]map[string]string{}
	var userName string
	resInterface := viper.Get(GetEnv() + "." + keyDbUsers).([]interface{})
	for _, item := range resInterface {
		itemTyped := item.(map[string]interface{})
		resMap := map[string]string{}
		for key, value := range itemTyped {
			valueTyped := value.(string)
			if key == "name" {
				userName = valueTyped
			}
			resMap[key] = valueTyped
		}
		resSlice[userName] = resMap
	}
	return resSlice
}

func GetDbUsers() map[string]map[string]string {
	var res map[string]map[string]string
	cacheApp := GetCache()
	resFromCache, ok := cacheApp.Get(keyDbUsers)
	if ok {
		res = resFromCache.(map[string]map[string]string)
	} else {
		res = GetNoCachedDbUsers()
		cacheApp.Set(keyDbUsers, res, cache.NoExpiration)
	}
	return res
}

func GetNoCachedDbConnection(name string) *sql.DB {
	dbUsers := GetNoCachedDbUsers()
	connStr := fmt.Sprintf(GetNoCachedString(keyDbUriTemplate), name, dbUsers[name][keyPassword])
	db, err := sql.Open(dbDriverName, connStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	db.SetMaxIdleConns(2)
	//defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

func GetDbConnection(name string) *sql.DB {
	var res *sql.DB
	var keyDbConnection = "db_connection_" + name
	cacheApp := GetCache()
	resFromCache, ok := cacheApp.Get(keyDbConnection)
	if ok {
		res = resFromCache.(*sql.DB)
	} else {
		res = GetNoCachedDbConnection(name)
		cacheApp.Set(keyDbConnection, res, cache.NoExpiration)
	}
	return res
}
