#!/bin/bash -e

docker-compose down
docker-compose up -d
sleep 5
go run database/scripts/run_migrations/main.go
mkdir -p database/scripts/.tmp
./database/scripts/dump_pg_schema.sh > database/scripts/.tmp/migrations.dump.sql

docker-compose down
docker-compose up -d
sleep 5
go run database/scripts/init_schema/main.go
./database/scripts/dump_pg_schema.sh > database/scripts/.tmp/schema.dump.sql

if diff --color=auto database/scripts/.tmp/migrations.dump.sql database/scripts/.tmp/schema.dump.sql; then
    echo
    echo "Le schéma en base est identique à celui généré par les migrations !"
    echo
    exit 0
else
    echo "/!\\"
    echo "/!\\ Attention, le schéma en base est différent de celui généré par les migrations !"
    echo "/!\\"
    exit 1
fi
