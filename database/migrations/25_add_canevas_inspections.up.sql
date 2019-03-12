CREATE EXTENSION unaccent;

CREATE TABLE "canevas" (
  "id" bigserial NOT NULL,
  "nom" text unique NOT NULL,
  "auteur_id" bigint NOT NULL,
  "data_version" bigint,
  "data" jsonb,
  "created_at" timestamptz,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id")
);
