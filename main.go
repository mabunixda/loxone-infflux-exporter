package main

import (
    "flag"
    "net/http"
    "runtime"
    "os"
        "encoding/json"

    "github.com/Sirupsen/logrus"

    "github.com/influxdata/influxdb/client/v2"
)

var (
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

    logrus.Infof("Miniserver: %s", configuration.Loxone.Address)
    logrus.Infof("Defined metrics: %d", len(configuration.Metrics))

}

func influxDBClient(config InfluxDbConfig) client.Client {
    c, err := client.NewHTTPClient(client.HTTPConfig {
        Addr:    config.Address,
        Username:    config.Username,
        Password:    config.Password,
    })
    if err != nil {
        logrus.Fatalf("Error at influx connection: %s", err);
    }
    return c
}
func main() {
    runtime.GOMAXPROCS(6)
    flag.Parse()
    c := influxDBClient(configuration.InfluxDb)
    for i := 0; i < len(configuration.Metrics); i++ {
        go singleNode(configuration.Metrics[i], configuration.Loxone, c)
    }
    router := NewRouter()
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
    http.ListenAndServe(addr, router)
}
