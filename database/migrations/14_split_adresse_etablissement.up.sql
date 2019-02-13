ALTER TABLE "etablissements"
  DROP COLUMN "adresse",
  ADD COLUMN "adresse1" text,
  ADD COLUMN "adresse2" text,
  ADD COLUMN "code_postal" text,
  ADD COLUMN "commune" text;
