package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/arsmine/screen-monitor/config"
	"github.com/arsmine/screen-monitor/stat"

	"github.com/gorilla/mux"
)

func Start(mainCfg *config.MainConfig) {
	apiListen := mainCfg.Listen

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
	isAllowed := false
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	for _, value := range stat.ActiveScreensStruct.AllowedIPs {
		if ip == value {
			isAllowed = true
			log.Printf("%s ip is allowed", ip)
		}
	}

	if !isAllowed {
		log.Printf("%s ip is not allowed.\n", ip)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	log.Println(stat.ActiveScreensStruct)

	osStat := stat.ReturnSystemStats()

	json.NewEncoder(w).Encode(osStat)
	return
}

func collectStrSystemStats(w http.ResponseWriter, req *http.Request) {
	isAllowed := false
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	for _, value := range stat.ActiveScreensStruct.AllowedIPs {
		if ip == value {
			isAllowed = true
			log.Printf("%s ip is allowed", ip)
		}
	}

	if !isAllowed {
		log.Printf("%s ip is not allowed.\n", ip)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	strOsStat := stat.ReturnStrSystemStats()

	json.NewEncoder(w).Encode(strOsStat)
	return
}

func collectScreenStats(w http.ResponseWriter, req *http.Request) {
	isAllowed := false
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	for _, value := range stat.ActiveScreensStruct.AllowedIPs {
		if ip == value {
			isAllowed = true
			log.Printf("%s ip is allowed", ip)
		}
	}

	if !isAllowed {
		log.Printf("%s ip is not allowed.\n", ip)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	screenStat := stat.ReturnScreenStats()
	/*
		if err != nil {
			log.Println(err)
			errorStruct["status"] = "500 Internal Server Error"
			errorStruct["error"] = "unexpected error while listing addresses"
			json.NewEncoder(w).Encode(errorStruct)
			return
		}*/

	json.NewEncoder(w).Encode(screenStat)
	return
}
