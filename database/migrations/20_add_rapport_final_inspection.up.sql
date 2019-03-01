CREATE TABLE "rapports" (
  "id" bigserial,
  "nom" text,
  "type" text,
  "taille" bigint,
  "storage_id" text UNIQUE,
  "auteur_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id")
);

ALTER TABLE "inspections"
  ADD COLUMN "rapport_id" bigint,
  ADD FOREIGN KEY ("rapport_id") REFERENCES "rapports" ("id") ON DELETE SET NULL;
