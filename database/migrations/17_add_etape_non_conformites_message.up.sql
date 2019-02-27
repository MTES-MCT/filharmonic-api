AlTER TABLE "messages" ADD COLUMN "etape_traitement_non_conformites" boolean not null;
AlTER TABLE "constats" ADD COLUMN "date_resolution" timestamptz;
