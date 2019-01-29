CREATE TABLE "suites" (
  "id" bigserial NOT NULL,
  "type" text,
  "synthese" text,
  "delai" bigint,
  PRIMARY KEY ("id")
);

ALTER TABLE "inspections"
  ADD COLUMN "suite_id" bigint,
  ADD FOREIGN KEY ("suite_id") REFERENCES "suites" ("id") ON DELETE SET NULL;

ALTER TABLE "constats"
  DROP COLUMN "deleted_at",
  DROP COLUMN "auteur_id";
