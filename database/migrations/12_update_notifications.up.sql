ALTER TABLE "notifications"
  RENAME COLUMN "created_at" TO "read_at";
ALTER TABLE "notifications"
  ALTER COLUMN "lecteur_id" DROP NOT NULL;
