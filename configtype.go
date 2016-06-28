package main;

import "github.com/prometheus/client_golang/prometheus"


type MetricValue struct { 
	Name string 			`json:"name"`
	ValuePath string		`json:"valuepath"`
	Gauge prometheus.Gauge
}
type Metric struct {
	URI string				`json:"URI"`
	Interval int			`json:"interval"`
	Values []MetricValue	`json:"Values"`
}

type Configuration struct {
	Miniserver	string		`json:"Miniserver"`
	Authentication string	`json:"Authentication"`
	Metrics 	[]Metric	`json:"Metrics"`
}
