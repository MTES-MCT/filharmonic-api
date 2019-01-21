CREATE TABLE "point_de_controles" (
  "id" bigserial,
  "sujet" text,
  "references_reglementaires" text[],
  "publie" boolean NOT NULL,
  "inspection_id" bigint NOT NULL,

  PRIMARY KEY ("id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id")
);

CREATE TABLE "commentaires" (
  "id" bigserial,
  "message" text,
  "date" timestamptz,
  "auteur_id" bigint NOT NULL,
  "inspection_id" bigint NOT NULL,

  PRIMARY KEY ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id")
);

CREATE TABLE "messages" (
  "id" bigserial,
  "message" text,
  "date" timestamptz,
  "lu" boolean NOT NULL,
  "interne" boolean NOT NULL,
  "auteur_id" bigint NOT NULL,
  "point_de_controle_id" bigint NOT NULL,

  PRIMARY KEY ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("point_de_controle_id") REFERENCES "point_de_controles" ("id")
);

ALTER TABLE "etablissements"
  ALTER COLUMN "iedmtd" SET NOT NULL;

ALTER TABLE "inspections"
  ALTER COLUMN "annonce" SET NOT NULL,
  ADD COLUMN "themes" text[];

DROP TABLE "theme_inspections";
