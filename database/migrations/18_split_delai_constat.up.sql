ALTER TABLE "constats"
  ADD COLUMN "echeance_resolution" date,
  ADD COLUMN "delai_nombre" integer,
  ADD COLUMN "delai_unite" text;

UPDATE "constats" SET "delai_unite" = 'mois';
UPDATE "constats" SET "delai_nombre" = 1 where "delai" = '1 mois';
UPDATE "constats" SET "delai_nombre" = 3 where "delai" = '3 mois';
UPDATE "constats" SET "delai_nombre" = 6 where "delai" = '6 mois';

ALTER TABLE "constats" DROP COLUMN "delai";
