package pkg

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/iigou/town-diff/pkg/api"
	"github.com/iigou/town-diff/pkg/internal"
)

const apiPathName = "api"
const townPathName = "town"
const diffPathName = "diff"

var townSvc api.ITownSvc = &internal.TownSvc{}
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

type supplier func(body io.ReadCloser, pathVars map[string]string, queryParams url.Values) (interface{}, error)

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
