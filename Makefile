
CONTAINER=r.nitram.at/lox-influx-exporter
TODAY=`date +'%Y%m%d'`

all: container

ui: ui/build/bundle.js
	cd ui; npm run build

bindata: ui
	go-bindata ui/build/

container: bindata
	docker build -t ${CONTAINER}:${TODAY} .
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest
	docker push ${CONTAINER}:latest

goreq:
	export GOTPATH=${PWD}
	go get github.com/Sirupsen/logrus
	go get github.com/influxdata/influxdb/client/v2
	go get github.com/oliveagle/jsonpath
	go get github.com/gorilla/mux

golang: goreq
	go build -v -o loxoneinfluxexporter

