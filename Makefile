
CONTAINER=loxprom
TODAY=`date +'%Y%M%d'`

container: golang
	docker build -t ${CONTAINER}:${TODAY} . 
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest

golang:
	export GOPATH=${PWD}
	go build

all: container
	
