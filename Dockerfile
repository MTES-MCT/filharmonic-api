FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY filharmonic-api /usr/local/bin/filharmonic-api

WORKDIR /filharmonic

# add migrations in a specific folder
# because https://github.com/go-pg/migrations/blob/094c91e7aafa796c34d824a346645a568a605081/collection.go#L121-L138
COPY database/migrations/*.sql /go/src/github.com/MTES-MCT/filharmonic-api/database/migrations/
COPY cron/templates/*.tmpl cron/templates/

EXPOSE 5000
CMD ["/usr/local/bin/filharmonic-api"]
