ALTER TABLE "point_de_controles" ADD COLUMN "order" bigint;

UPDATE "point_de_controles" SET "order" = "id";
