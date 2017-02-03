FROM alpine:latest
MAINTAINER Martin Buchleitner "martin@nitram.at"
ENV HTTP_PROXY http://squid.avl.com:8080
ENV HTTPS_PROXY http://squid.avl.com:8080

RUN apk --no-cache add ca-certificates
ADD loxonegoprometheus /opt/loxonegoprometheus
RUN chmod 755 /opt/loxonegoprometheus
    
VOLUME "/opt/config"

WORKDIR "/opt"
EXPOSE 8080
ENTRYPOINT ["./loxonegoprometheus"]
