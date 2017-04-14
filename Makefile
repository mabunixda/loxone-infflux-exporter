
CONTAINER=r.nitram.at/loxprom
TODAY=`date +'%Y%m%d'`

container: golang
	docker build -t ${CONTAINER}:${TODAY} . 
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest

golang: goreq
	export GOPATH=${PWD}
	go build

all: container
	
goreq:
	export GOTPATH=${PWD}
	go get github.com/Sirupsen/logrus
	go get github.com/influxdata/influxdb/client/v2
	go get gopkg.in/xmlpath.v2
	go get github.com/gorilla/mux
