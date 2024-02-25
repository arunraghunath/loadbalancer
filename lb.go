package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	"github.com/arunraghunath/loadbalancer/server"
)

type LBConfig struct {
	LBport    string           `json:"lbport"`
	Servers   []*server.Server `json:"servers"`
	RobinCnt  int
	LBLock    sync.Mutex
	Algorithm *string
}

var cfg LBConfig

func main() {
	algoType := flag.String("alg", "robin", "algorithm options (robin)")
	flag.Parse()
	cfg.Algorithm = algoType
	wg := new(sync.WaitGroup)
	ReadJson("config.json", &cfg)
	wg.Add(1)
	go func() {
		StartLB(cfg.LBport)
		wg.Done()
	}()

	for i := 0; i < len(cfg.Servers); i++ {
		wg.Add(1)
		lserver := cfg.Servers[i]
		go func(lserver *server.Server) {
			lserver.StartServer()
			wg.Done()
		}(lserver)
	}
	wg.Wait()
}

func HandleHC(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the health check service. About to print the status of each server.")
	fmt.Fprintf(w, "Status of Servers as below\n")
	for i := 0; i < len(cfg.Servers); i++ {
		fmt.Fprintf(w, "Name --> %s, Status --> %t\n", cfg.Servers[i].Name, cfg.Servers[i].IsHealthy())
	}
}

func HandleLB(w http.ResponseWriter, r *http.Request) {
	switch *cfg.Algorithm {
	case "robin":
		RoundRobin(w, r)
	}
}

func ReadJson(filename string, lbconfig *LBConfig) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(fileBytes, lbconfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func StartLB(port string) {
	fmt.Println("Starting Loadbalancer")
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleLB)
	mux.HandleFunc("/hc", HandleHC)
	http.ListenAndServe(port, mux)
}

func RoundRobin(w http.ResponseWriter, r *http.Request) {
	maxServers := len(cfg.Servers)
	cfg.LBLock.Lock()
	curr := cfg.RobinCnt % maxServers

	s := cfg.Servers[curr]
	if s.IsHealthy() {
		url, err := url.Parse(s.URL)
		if err != nil {
			log.Fatal(err.Error())
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(url)
		cfg.RobinCnt++
		reverseProxy.ServeHTTP(w, r)
	} else {
		cfg.RobinCnt++
		HandleLB(w, r)

	}

	cfg.LBLock.Unlock()
}
