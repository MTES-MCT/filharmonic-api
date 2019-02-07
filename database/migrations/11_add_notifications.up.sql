CREATE TABLE "evenements" (
  "id" bigserial,
  "type" text,
  "created_at" timestamptz,
  "data" json,
  "auteur_id" bigint NOT NULL,
  "inspection_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id")
);

CREATE TABLE "notifications" (
  "id" bigserial,
  "created_at" timestamptz,
  "lue" boolean NOT NULL,
  "evenement_id" bigint NOT NULL,
  "lecteur_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("lecteur_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("evenement_id") REFERENCES "evenements" ("id")
);
