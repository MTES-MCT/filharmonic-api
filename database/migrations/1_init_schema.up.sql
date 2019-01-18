CREATE TABLE "etablissements" (
  "id" bigserial,
  "s3ic" text UNIQUE,
  "nom" text,
  "raison" text,
  "adresse" text,
  "seveso" text,
  "activite" text,
  "iedmtd" boolean,
  PRIMARY KEY ("id")
);
CREATE TABLE "users" (
  "id" bigserial,
  "nom" text,
  "prenom" text,
  "email" text UNIQUE,
  "password" text,
  "profile" text,
  PRIMARY KEY ("id")
);
CREATE TABLE "etablissement_to_exploitants" (
  "etablissement_id" bigserial,
  "user_id" bigserial,
  PRIMARY KEY ("etablissement_id", "user_id"),
  FOREIGN KEY ("etablissement_id") REFERENCES "etablissements" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
