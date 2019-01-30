CREATE TABLE "piece_jointes" (
  "id" bigserial,
  "nom" text,
  "type" text,
  "taille" bigint,
  "storage_id" text UNIQUE,
  "message_id" bigint,
  "commentaire_id" bigint,
  "auteur_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("message_id") REFERENCES "messages" ("id"),
  FOREIGN KEY ("commentaire_id") REFERENCES "commentaires" ("id"),
  FOREIGN KEY ("auteur_id") REFERENCES "users" ("id")
);
