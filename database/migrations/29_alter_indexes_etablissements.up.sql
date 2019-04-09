CREATE OR REPLACE FUNCTION f_unaccent(text)
  RETURNS text AS
$func$
SELECT public.unaccent('public.unaccent', $1)  -- schema-qualify function and dictionary
$func$  LANGUAGE sql IMMUTABLE;

DROP INDEX trgm_idx_etablissments_nom;
CREATE INDEX trgm_idx_etablissments_nom ON etablissements USING gin (f_unaccent(nom) gin_trgm_ops);
DROP INDEX trgm_idx_etablissments_raison;
CREATE INDEX trgm_idx_etablissments_raison ON etablissements USING gin (f_unaccent(raison) gin_trgm_ops);
DROP INDEX trgm_idx_etablissments_adresse1;
CREATE INDEX trgm_idx_etablissments_adresse1 ON etablissements USING gin (f_unaccent(adresse1) gin_trgm_ops);
DROP INDEX trgm_idx_etablissments_adresse2;
CREATE INDEX trgm_idx_etablissments_adresse2 ON etablissements USING gin (f_unaccent(adresse2) gin_trgm_ops);
DROP INDEX trgm_idx_etablissments_commune;
CREATE INDEX trgm_idx_etablissments_commune ON etablissements USING gin (f_unaccent(commune) gin_trgm_ops);
