
CONTAINER=r.nitram.at/lox-influx-exporter
TODAY=`date +'%Y%m%d'`

container: 
	docker build -t ${CONTAINER}:${TODAY} .
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest
	docker push ${CONTAINER}:latest

golang: goreq
	export GOPATH=${PWD}
	go build -v -o loxoneinfluxexporter

all: container

goreq:
	export GOTPATH=${PWD}
	go get github.com/Sirupsen/logrus
	go get github.com/influxdata/influxdb/client/v2
	go get github.com/oliveagle/jsonpath
	go get github.com/gorilla/mux
