package main

import (
	"flag"
	"net/http"
	"time"
	"runtime"
	"os"
    "encoding/json"
	"strconv"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/xmlpath.v2"
)

type MetricValue struct { 
	Name string
	ValuePath string
	Gauge prometheus.Gauge
}
type Metric struct {
	URI string
	Values []MetricValue
}

type Configuration struct {
	Miniserver	string
	Authentication string
	Metrics 	[]Metric
}

const (
	timeout = 10 * 1000
)

var (
	valueRegex = regexp.MustCompile("^(\\d+([\\.,]\\d+)?)")
	addr string
	config string
	configuration Configuration
)

func init() {
	flag.StringVar(&addr,"listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&config,"config", "config.json", "The configuration file which data should be requested")

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

func queryData(url string, auth string) (*xmlpath.Node,error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
          	return nil, err
  	}
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	req.Header.Add("Authorization", "Basic " + auth)
	req.Header.Add("Accept", "application/xml")
    		
   	resp, err := client.Do(req)
   	if err != nil {
		logrus.Errorf("Failure on getting request: %s", err)
    	return nil, err
  	}
   	defer resp.Body.Close()
	xmlData, err := xmlpath.Parse(resp.Body)
    if err != nil {
		logrus.Errorf("Failure on parsing request %s, %s", url, err)
        return nil, err
    }
    return xmlData, nil
}

func pushData(vals []MetricValue, data *xmlpath.Node) {

	for i:=0; i < len(vals); i++ {
		// logrus.Infof("%s",vals[i].ValuePath)
		path := xmlpath.MustCompile(vals[i].ValuePath)
		if value, ok := path.String(data); ok {

			if valueRegex.MatchString(value) == false {
				logrus.Errorf("not a valid value: %s", value)
				continue
			}

			parsedValue := valueRegex.FindAllString(value,-1)
			if len(parsedValue) == 0 {
				logrus.Errorf("Could not parse values: %s", value)
				continue
			}
			// logrus.Infof("current value: %s", parsedValue)
			f, err := strconv.ParseFloat(parsedValue[0], 64)

			if err != nil {
				logrus.Errorf("Failure on getting float value: %s", err)
				continue
			}
//		    logrus.Infof("%d", f)		
			vals[i].Gauge.Set(f)
		}
	}
}

func singleNode(m Metric, server string, auth string) {
	url := "http://" + server + "/" + m.URI
	for i := 0; i < len(m.Values); i++ {
		if m.URI == "" {
			logrus.Fatalf("No url defined in a configuration")
		}
		if m.Values[i].ValuePath == "" {
			m.Values[i].ValuePath = "/LL/@value"
		}

		m.Values[i].Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "Haus",
				Subsystem: "",
				Name: m.Values[i].Name,// .Replace(" ", "_",-1).ToLower(),
				Help: "Loxone Prometheus",
		})
		prometheus.MustRegister(m.Values[i].Gauge)
	}
	for {
		body,err := queryData(url, auth)
		if err == nil { 
			pushData(m.Values, body)
		}
		time.Sleep(time.Duration(timeout * time.Millisecond))
	}
}


func main() {
	runtime.GOMAXPROCS(6)
	flag.Parse()
	for i := 0; i < len(configuration.Metrics); i++ {
		go singleNode(configuration.Metrics[i], configuration.Miniserver, configuration.Authentication)
	}
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(addr, nil)
}