CREATE TABLE "inspections" (
  "id" bigserial,
  "date" date,
  "type" text,
  "annonce" boolean,
  "origine" text,
  "circonstance" text,
  "detail_circonstance" text,
  "contexte" text,
  "etat" text,
  "etablissement_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("etablissement_id") REFERENCES "etablissements" ("id")
);
CREATE TABLE "theme_inspections" (
  "id" bigserial,
  "nom" text,
  "inspection_id" bigint,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id")
);
CREATE TABLE "theme_referentiels" (
  "id" bigserial,
  "nom" text,
  PRIMARY KEY ("id")
);
CREATE TABLE "inspection_to_inspecteurs" (
  "inspection_id" bigserial,
  "user_id" bigserial,
  PRIMARY KEY ("inspection_id", "user_id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
