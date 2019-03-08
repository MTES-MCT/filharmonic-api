ALTER TABLE "constats"
  ADD COLUMN "notification_echeance_expiree_envoyee" boolean default false not null;
ALTER TABLE "constats"
  RENAME COLUMN "rappel_echeances_envoye" TO "notification_rappel_echeance_envoyee";
