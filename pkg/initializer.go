package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/iigou/town-diff/pkg/api"
	"github.com/iigou/town-diff/pkg/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type supplier func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error)

const apiPathName = "api"
const townPathName = "town"
const diffPathName = "diff"

var dbConn *gorm.DB
var townSvc api.ITownSvc = &internal.TownSvc{DbConnFn: GetDatabaseConnection}
var townDiff api.ITownDiffSvc = &internal.TownDiff{}

func registerTownSvcRouters(r *mux.Router) {
	r.
		Name("get towns").
		Methods("GET").
		PathPrefix("/" + townPathName).
		HandlerFunc(handler(
			func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error) {
				request := api.Town{}
				err := json.NewDecoder(body).Decode(&request)
				if err == io.EOF {
					log.Println("No body to deserialise... ")
					if len(queryParams) > 0 {
						if err = schema.NewDecoder().Decode(&request, queryParams); err != nil {
							log.Println("Error during the deserialization of queryParams: ", err)
						}
					}
				}
				return townSvc.Get(&request)
			}))

	r.
		Name("save town").
		Methods("POST").
		PathPrefix("/" + townPathName).
		HandlerFunc(handler(
			func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error) {
				request := api.Town{}
				err := json.NewDecoder(body).Decode(&request)
				if err == io.EOF {
					log.Println("No body to deserialise... ")
				}
				return townSvc.Save(&request)
			}))

	r.
		Name("update town").
		Methods("PATCH").
		PathPrefix("/" + townPathName + "/{id}").
		HandlerFunc(handler(
			func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error) {
				request := api.Town{}
				err := json.NewDecoder(body).Decode(&request)
				if err == io.EOF {
					log.Println("No body to deserialise... ")
				}
				return townSvc.Update(pathVars["id"], &request)
			}))

	r.
		Name("delete town").
		Methods("DELETE").
		PathPrefix("/" + townPathName + "/{id}").
		HandlerFunc(handler(
			func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error) {
				return townSvc.Delete(pathVars["id"])
			}))
}

func InitRouters() http.Handler {
	router := mux.NewRouter()
	apiRoute := router.NewRoute().Name(apiPathName).PathPrefix("/" + apiPathName)
	registerTownSvcRouters(apiRoute.Subrouter())
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(responseContentTypeMw)
	return handlers.LoggingHandler(os.Stdout, router)
}

func handler(supp supplier) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		response, err := supp(request.Body, mux.Vars(request), request.URL.Query())

		if err != nil {
			log.Println("Response is erroneus, Body:", response)
			w.WriteHeader(500)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(err.Error()))
		}

		if response != nil {
			w.WriteHeader(200)
			encoded, _ := json.MarshalIndent(response, "", "  ")
			log.Println("Response is 200, Body:", string(encoded))
			w.Write(encoded)
		}
	}
}

func responseContentTypeMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func CreateDBConnection() {

	config, err := getConnectionConfig()
	if err != nil {
		panic(err)
	}

	log.Println("Connecting to ", config["url"], " with user ", config["user"])

	//user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s%s", config["user"], config["pwd"], config["url"])), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Create the connection pool

	sqlDB, err := db.DB()

	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db

	log.Println("Migrating Town table")
	dbConn.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&api.Town{})

}

func GetDatabaseConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return dbConn, err
	}
	if err := sqlDB.Ping(); err != nil {
		return dbConn, err
	}
	return dbConn, nil
}

func getConnectionConfig() (map[string]string, error) {
	// Open our jsonFile
	configFile, err := os.Open("./config.json")
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully Opened configFile.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var config map[string]interface{}
	json.Unmarshal([]byte(byteValue), &config)
	//user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dbConfig := config["db"].(map[string]interface{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("@tcp(%s:%s)/", dbConfig["url"].(string), dbConfig["port"].(string)))
	if dbConfig["name"] != nil {
		sb.WriteString(dbConfig["name"].(string))
	}
	if dbConfig["args"] != nil {
		sb.WriteString("?")
		sb.WriteString(dbConfig["args"].(string))
	}

	return map[string]string{
		"user": dbConfig["user"].(string),
		"pwd":  decode64(dbConfig["pwd"].(string)),
		"url":  sb.String(),
	}, nil
}

func decode64(in string) string {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		log.Println("Error when decoding input. ", err)
		return string([]byte{})
	}

	return string(out)
}
