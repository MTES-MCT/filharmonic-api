ALTER TABLE "messages"
  ADD COLUMN "etape_traitement_non_conformites" boolean default false not null,
  ALTER COLUMN "lu" set default false,
  ALTER COLUMN "interne" set default false;
ALTER TABLE "point_de_controles" ALTER COLUMN "publie" set default false;
ALTER TABLE "inspections" ALTER COLUMN "annonce" set default false;
UPDATE "suites" SET "penal_engage" = false;
ALTER TABLE "suites" ALTER COLUMN "penal_engage" set default false;
ALTER TABLE "suites" ALTER COLUMN "penal_engage" set not null;
ALTER TABLE "constats" ADD COLUMN "date_resolution" timestamptz;
