#!/bin/bash -e

docker-compose exec db su - postgres -c "pg_dump --schema-only --create --no-owner --clean --exclude-table=gopg_migrations filharmonic" | sed '/^--/d' | perl -0pe 's/CREATE SEQUENCE public\.gopg_migrations_id[^;]*;//' | sed '/^\s*$/d' | grep -v gopg_migrations
