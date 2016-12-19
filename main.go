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
	Interval int
	Values []MetricValue
}

type Configuration struct {
	Miniserver	string
	Authentication string
	Metrics		[]Metric
}


var (
	valueRegex = regexp.MustCompile("^(-?\\d+([\\.,]\\d+)?)")
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

func queryData(url string, auth string) (*xmlpath.Node,error) {
	client := &http.Client{}
	logrus.Infof("Current url %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Could not make request to miniserver %s", err)
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
	logrus.Infof("Parsing values...")
	for i:=0; i < len(vals); i++ {
		path := xmlpath.MustCompile(vals[i].ValuePath)
		logrus.Infof("Compiled xlm path")
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
			logrus.Infof("current value: %s", parsedValue)
			f, err := strconv.ParseFloat(parsedValue[0], 64)

			if err != nil {
				logrus.Errorf("Failure on getting float value: %s", err)
				continue
			}
		    	logrus.Infof("%d", f)		
			vals[i].Gauge.Set(f)
		} else {
			logrus.Errorf("Could not parse path: %s", data)
}
	}
}

func singleNode(m Metric, server string, auth string) {
	url := "http://" + server +  m.URI
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
		logrus.Infof("Current url: %s", url)
		body,err := queryData(url, auth)
		if err == nil { 
			pushData(m.Values, body)
		} else {
			logrus.Errorf("Could not push datavalues %s", err)
		}
		interval := timeout 
		if m.Interval > 0 {
			interval = m.Interval
		}
		time.Sleep(time.Duration(interval * int(time.Second)))
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
