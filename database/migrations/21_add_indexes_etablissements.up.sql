CREATE EXTENSION pg_trgm;
CREATE INDEX trgm_idx_etablissments_s3ic ON etablissements USING gin (s3ic gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_nom ON etablissements USING gin (nom gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_raison ON etablissements USING gin (raison gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_adresse1 ON etablissements USING gin (adresse1 gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_adresse2 ON etablissements USING gin (adresse2 gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_code_postal ON etablissements USING gin (code_postal gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_commune ON etablissements USING gin (commune gin_trgm_ops);

CREATE INDEX idx_users_profile ON users(profile);
CREATE INDEX idx_inspections_etat ON inspections(etat);
CREATE INDEX idx_point_de_controles_publie ON point_de_controles(publie);
CREATE INDEX idx_messages_interne ON messages(interne);
CREATE INDEX idx_messages_lu ON messages(lu);
