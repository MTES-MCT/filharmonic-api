COPY ( select inspection.etablissement_id, point_de_controle.id as point_de_controle_id, constat.type as constat_type, constat.date_resolution, constat.echeance_resolution
from inspections as inspection
join point_de_controles as point_de_controle on point_de_controle.inspection_id = inspection.id
left join constats as constat on point_de_controle.constat_id = constat.id ) TO STDOUT WITH CSV HEADER;
