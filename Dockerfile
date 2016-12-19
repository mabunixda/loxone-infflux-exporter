FROM debian:jessie
MAINTAINER Martin Buchleitner "martin@nitram.at"
RUN apt-get update
RUN apt-get install -y ca-certificates

COPY loxonegoprometheus /opt/loxonegoprometheus
RUN chmod 755 /opt/loxonegoprometheus

WORKDIR "/opt"
EXPOSE 8080
CMD ["/opt/loxonegoprometheus"]
