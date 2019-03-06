ALTER TABLE "inspections"
  ADD COLUMN "validation_rejetee" boolean default false not null,
  ADD COLUMN "motif_rejet_validation" text;
