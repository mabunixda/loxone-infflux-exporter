package main;


type MetricValue struct {
    Name        string        `json:"name"`
    ValuePath    string        `json:"valuepath"`
}
type Metric struct {
    Name        string        `json:Name`
    URI        string        `json:"URI"`
    Interval    int        `json:"interval"`
    Values        []MetricValue    `json:"Values"`
}
type InfluxDbConfig struct {
    Address        string        `json:"Address"`
    Username    string        `json:"Username"`
    Password    string        `json:"Password"`
    Interval    int            `json:"Interval"`
}
type LoxoneConfig struct {
    Address        string        `json:"Address"`
    Authentication    string        `json:"Authentication"`
}
type Configuration struct {
    InfluxDb    InfluxDbConfig    `json:"Influx"`
    Loxone        LoxoneConfig    `json:"Loxone"`
    Metrics        []Metric    `json:"Metrics"`
}
