CREATE TABLE "user_to_favoris" (
  "inspection_id" bigserial,
  "user_id" bigserial,
  PRIMARY KEY ("inspection_id", "user_id"),
  FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
