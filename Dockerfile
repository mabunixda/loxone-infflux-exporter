FROM alpine:latest
MAINTAINER Martin Buchleitner "martin@nitram.at"

RUN apk --no-cache add ca-certificates
ADD loxoneinfluxexporter /opt/loxoneinfluxexporter
RUN chmod 755 /opt/loxoneinfluxexporter
    

WORKDIR "/opt"
EXPOSE 8080
ENTRYPOINT ["/opt/loxoneinfluxexporter"]
