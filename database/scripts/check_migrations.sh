#!/bin/bash -e

cleanDatabase() {
  echo 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;' | PGPASSWORD=filharmonic psql -h localhost -p 5432 -U filharmonic filharmonic
  sleep 5
}
mkdir -p database/scripts/.tmp

cleanDatabase
go run database/scripts/run_migrations/main.go
./database/scripts/dump_pg_schema.sh > database/scripts/.tmp/migrations.dump.sql

cleanDatabase
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
