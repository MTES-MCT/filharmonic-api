FROM alpine:latest

COPY filharmonic-api /usr/local/bin/filharmonic-api

EXPOSE 5000
CMD ["/usr/local/bin/filharmonic-api"]
