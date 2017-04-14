package main

import (
    "net/http"
    "time"
    "strconv"
    "regexp"
	"encoding/json"

    "github.com/Sirupsen/logrus"
    "github.com/influxdata/influxdb/client/v2"
	"github.com/oliveagle/jsonpath"
)


const (
    database = "loxone"
)
var (
    valueRegex = regexp.MustCompile("^(-?\\d+([\\.,]\\d+)?)")
)
func queryData(url string, auth string) (interface{},error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Add("Content-Type", "application/json; charset=utf-8")
    req.Header.Add("Authorization", "Basic " + auth)
    req.Header.Add("Accept", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        logrus.Errorf("Failure on getting request: %s", err)
    	return nil, err
    }
    defer resp.Body.Close()
	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
    logrus.Errorf("Failure on parsing request %s, %s", url, err)
    return nil, err

	}
	logrus.Infof("%s", jsonData)
    return jsonData, nil
}

func pushData(m Metric, data interface{}, bp client.BatchPoints) {

    // fields are values of a sensor
    tags := map[string]string {}
    fields := map[string]interface{} {}

    vals := m.Values
    for i:=0; i < len(vals); i++ {
		jsonValue, err := jsonpath.JsonPathLookup(data, vals[i].ValuePath)
		if err != nil {
			continue
		}
		value, ok := jsonValue.(string)
		if ok == false {
			continue
		}
        if valueRegex.MatchString(value) == false {
            logrus.Errorf("not a valid value: %s", value)
            continue
        }

        parsedValue := valueRegex.FindAllString(value,-1)
        if len(parsedValue) == 0 {
            logrus.Errorf("Could not parse values: %s", value)
            continue
        }
        f, err := strconv.ParseFloat(parsedValue[0], 64)

        if err != nil {
            logrus.Errorf("Failure on getting float value: %s", err)
            continue
        }
		logrus.Infof("Setting %s to %s", vals[i].Name, value)
	    fields[vals[i].Name] = f
    }
	logrus.Infof("Adding new datapoint")
    // use sensors name as newpoint
    pt, err := client.NewPoint(m.Name,
                          tags,
                          fields,
                          time.Now(),
                )
    if err != nil {
        logrus.Errorf("Could not add new point %s", err)
    }
    bp.AddPoint(pt)

}

func singleNode(m Metric, loxConfig LoxoneConfig, c client.Client) {
    url := "http://" + loxConfig.Address  + m.URI
    for i := 0; i < len(m.Values); i++ {
        if m.URI == "" {
            logrus.Fatalf("No url defined in a configuration")
        }
        if m.Values[i].ValuePath == "" {
            m.Values[i].ValuePath = "$.LL.value"
        }
    }
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
                    Database:  database,
                    Precision: "s",
    })

    if err != nil {
        logrus.Errorf("Error at push: %s", err)
    }
    for {
        body,err := queryData(url, loxConfig.Authentication)
        if err == nil {
            pushData(m, body, bp)
            if err := c.Write(bp); err != nil {
                logrus.Errorf("could not write bachpoints: %s", err)
            }
        } else {
            logrus.Errorf("%s", err)
        }
        interval := timeout
        if m.Interval > 0 {
            interval = m.Interval
        }
        time.Sleep(time.Duration(interval * int(time.Second)))
    }
}
