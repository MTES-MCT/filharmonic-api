CREATE TABLE "departements" (
  "id" bigserial,
  "code_insee" text UNIQUE NOT NULL,
  "nom" text,
  "charniere" text,
  "region" text,
  "charniere_region" text,
  PRIMARY KEY ("id")
);

ALTER TABLE "etablissements"
  ADD COLUMN "departement_id" bigint,
  ADD FOREIGN KEY ("departement_id") REFERENCES "departements" ("id");
