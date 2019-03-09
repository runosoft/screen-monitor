package api

import (
	"log"
	"net/http"
	"encoding/json"

	"../stat"

	"github.com/gorilla/mux"
)

func Start() {
	apiListen := "0.0.0.0:8080"

	router := mux.NewRouter()
	log.Println("Starting the Api")
	log.Printf("Api Listen: %s\n", apiListen)

	router.HandleFunc("/api/osstats",
		collectSystemStats).Methods("GET")
	router.HandleFunc("/api/strosstats",
		collectStrSystemStats).Methods("GET")
	router.HandleFunc("/api/screens",
		collectScreenStats).Methods("GET")
	/*
	router.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode("Message": "An unexpected request. Check url"})
	})*/

	log.Println(http.ListenAndServe(apiListen, router))
}

func collectSystemStats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	reply := make(map[string]interface{})
	errorStruct := make(map[string]string)

	osStat, err := stat.CollectSystemStats()
	if err != nil {
		log.Println(err)
		errorStruct["status"] = "500 Internal Server Error"
		errorStruct["error"] = "unexpected error while listing addresses"
		json.NewEncoder(w).Encode(errorStruct)
		return
	}

	reply["osStat"] = osStat
	json.NewEncoder(w).Encode(reply)
	return
}

func collectStrSystemStats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	reply := make(map[string]interface{})
	errorStruct := make(map[string]string)

	strOsStat, err := stat.CollectStrSystemStats()
	if err != nil {
		log.Println(err)
		errorStruct["status"] = "500 Internal Server Error"
		errorStruct["error"] = "unexpected error while listing addresses"
		json.NewEncoder(w).Encode(errorStruct)
		return
	}

	reply["strOsStat"] = strOsStat
	json.NewEncoder(w).Encode(reply)
	return
}

func collectScreenStats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	reply := make(map[string]interface{})
	errorStruct := make(map[string]string)

	screenStat, err := stat.CollectScreenStats()
	if err != nil {
		log.Println(err)
		errorStruct["status"] = "500 Internal Server Error"
		errorStruct["error"] = "unexpected error while listing addresses"
		json.NewEncoder(w).Encode(errorStruct)
		return
	}

	reply["result"] = screenStat
	json.NewEncoder(w).Encode(reply)
	return
}




