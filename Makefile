
CONTAINER="loxprom"
TODAY=`date +'%Y%M%D'`

container: golang
	docker build -t ${CONTAINER}:${TODAY} . 

golang:
	export GOPATH=${PWD}
	go build
all:
	
