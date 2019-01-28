 CREATE TABLE "constats" (
  "id" bigserial NOT NULL,
  "type" text,
  "remarques" text,
  "auteur_id" bigint NOT NULL,
  "deleted_at" timestamptz,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id")
);

ALTER TABLE "point_de_controles"
  ADD COLUMN "constat_id" bigint,
  ADD FOREIGN KEY ("constat_id") REFERENCES "constats" ("id") ON DELETE SET NULL;
