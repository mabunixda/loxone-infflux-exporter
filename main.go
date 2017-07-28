package main

import (
    "flag"
    "net/http"
    "runtime"
    "os"
    "net/url"
    "encoding/json"

    "github.com/Sirupsen/logrus"

)

var (
    timeout = 10
    addr string
    config string
    configuration Configuration
)

func init() {
    flag.StringVar(&addr,"publicUrl", "http://localhost:8080", "The address to listen on for HTTP requests.")
    flag.StringVar(&config,"config", "config.json", "The configuration file which data should be requested")
    flag.IntVar(&timeout, "interval", timeout, "the default interval of each metric in seconds")


}

func main() {
    runtime.GOMAXPROCS(6)
    flag.Parse()

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

	publicURL, err := url.Parse(addr)
	if err != nil {
		logrus.Fatalf("invalid public url:", err)
	}

    logrus.Infof("Miniserver: %s", configuration.Loxone.Address)
    logrus.Infof("Defined metrics: %d", len(configuration.Metrics))

    r, err := NewRouter(configuration, config, publicURL)

    go monitoring(configuration)

    http.ListenAndServe(":" + publicURL.Port(), r.router)
}
