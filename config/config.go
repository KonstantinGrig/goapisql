package config

import (
	"database/sql"
	"fmt"
	//Import the PostgreSQL driver
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
	"time"
)

var once sync.Once
var goapisqlCache *cache.Cache

const (
	KEY_DB_USERS        = "db_users"
	KEY_DB_URI_TEMPLATE = "db-uri-template"
	KEY_PASSWORD        = "password"
	KEY_JWT_SECRET      = "jwt-secret"
	PREFIX_TOKEN        = "Bearer "
	DB_DRIVER_NAME      = "postgres"
)

//Init config file from the config paths
func InitConfigFile(path string) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

//GetCache retrieves the *cache.Cache for your Application
func GetCache() *cache.Cache {
	once.Do(func() {
		goapisqlCache = cache.New(5*time.Minute, 10*time.Minute)
	})

	isFlushCache := os.Getenv("GOAPISQL_IS_CACHE_FLUSH")
	if isFlushCache == "yes" {
		goapisqlCache.Flush()
		os.Setenv("GOAPISQL_IS_CACHE_FLUSH", "no")
	}

	return goapisqlCache
}

//GetEnv retrieves the name of current environment
//default "prod"
func GetEnv() string {
	envVal := os.Getenv("GOAPISQL_ENV")
	if envVal == "" {
		envVal = "prod"
	}
	log.Println("Enviroment:", envVal)
	return envVal
}

//GetNoCachedString returns without cache the value associated with the key as a string.
func GetNoCachedString(key string) string {
	return viper.GetString(GetEnv() + "." + key)
}

//GetString retrieves the string value of the config variable named by the key.
//If string value have in cache then retrieves from cache
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

//GetNoCachedDbUsers retrieves map with iformation about DB user.
//example: map{postgres: {name: postgres, password: postgres_user_password}}
//Cache doesn't use
func GetNoCachedDbUsers() map[string]map[string]string {
	resSlice := map[string]map[string]string{}
	var userName string
	resInterface := viper.Get(GetEnv() + "." + KEY_DB_USERS).([]interface{})
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

//GetDbUsers retrieves map with iformation about DB user.
//example: map{postgres: {name: postgres, password: postgres_user_password}}
//Cache use
func GetDbUsers() map[string]map[string]string {
	var res map[string]map[string]string
	cacheApp := GetCache()
	resFromCache, ok := cacheApp.Get(KEY_DB_USERS)
	if ok {
		res = resFromCache.(map[string]map[string]string)
	} else {
		res = GetNoCachedDbUsers()
		cacheApp.Set(KEY_DB_USERS, res, cache.NoExpiration)
	}
	return res
}

//GetNoCachedDbConnection retrieves *sql.DB
//Cache doesn't use
func GetNoCachedDbConnection(name string) *sql.DB {
	dbUsers := GetNoCachedDbUsers()
	connStr := fmt.Sprintf(GetNoCachedString(KEY_DB_URI_TEMPLATE), name, dbUsers[name][KEY_PASSWORD])
	db, err := sql.Open(DB_DRIVER_NAME, connStr)
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

//GetDbConnection retrieves *sql.DB
//Cache use
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
