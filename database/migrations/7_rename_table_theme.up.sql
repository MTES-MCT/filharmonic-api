ALTER TABLE "theme_referentiels" RENAME TO "themes";
ALTER TABLE "themes" RENAME CONSTRAINT "theme_referentiels_pkey" TO "themes_pkey";
ALTER SEQUENCE "theme_referentiels_id_seq" RENAME TO "themes_id_seq";
