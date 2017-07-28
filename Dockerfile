FROM golang:1.8-alpine3.6 as builder
COPY *.go /go/src/
COPY Makefile /go/src/
WORKDIR /go/src/
RUN apk add --no-cache make git \
    && make golang


FROM alpine:3.6
MAINTAINER Martin Buchleitner "martin@nitram.at"
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/loxoneinfluxexporter /loxoneinfluxexporter
WORKDIR "/"
EXPOSE 8080
CMD ["/loxoneinfluxexporter"]
