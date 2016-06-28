package main

import (
	"flag"
	"net/http"
	"runtime"
	"os"
    "encoding/json"
	"regexp"

	"github.com/Sirupsen/logrus"	
	"github.com/prometheus/client_golang/prometheus"
)


var (
	valueRegex = regexp.MustCompile("^(\\d+([\\.,]\\d+)?)")
	timeout = 10
	addr string
	config string	
	configuration Configuration
)

func init() {
	flag.StringVar(&addr,"listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&config,"config", "config.json", "The configuration file which data should be requested")
	flag.IntVar(&timeout, "interval", timeout, "the default interval of each metric in seconds")
	if config == "" {
		logrus.Fatalf("You must pass a configuration file")
	}
	file, _ := os.Open(config)
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
  		logrus.Fatalf("error: %s", err)
	}

	logrus.Infof("Miniserver: %s", configuration.Miniserver)
	logrus.Infof("Defined metrics: %d", len(configuration.Metrics))

}

func main() {
	runtime.GOMAXPROCS(6)
	flag.Parse()
	for i := 0; i < len(configuration.Metrics); i++ {
		go singleNode(configuration.Metrics[i], configuration.Miniserver, configuration.Authentication)
	}
	router := NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	router.Handle("/metrics", prometheus.Handler())
	
 	http.ListenAndServe(addr, router)
}