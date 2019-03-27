#!/usr/bin/env node

const fs = require('fs')
const https = require('https')
const path = require('path')

// https://www.insee.fr/fr/information/2028040
// devrait être intégré dans les mois prochains à https://github.com/etalab/decoupage-administratif
// contient seulement les TOM/COM possédant des installations classées
const tomCom = `,
('975', 'Saint-Pierre-et-Miquelon', 'de ', '', ''),
('977', 'Saint-Barthélemy', 'de ', '', ''),
('978', 'Saint-Martin', 'de ', '', '')`

;(async function () {
  try {
    let regions = await fetchJSON('https://unpkg.com/@etalab/decoupage-administratif@0.3.0/data/regions.json')
    regions = regions.reduce((acc, region) => {
      acc[region.code] = {
        nom: region.nom,
        charniere: getCharniere(region.typeLiaison)
      }
      return acc
    }, {})
    const departements = await fetchJSON('https://unpkg.com/@etalab/decoupage-administratif@0.3.0/data/departements.json')
    let sourcecode = `INSERT INTO departements (code_insee, nom, charniere, region, charniere_region) VALUES ` +
      departements.map(departement => {
        return `\n('${departement.code}', '${escapeQuotes(departement.nom)}', '${escapeQuotes(getCharniere(departement.typeLiaison))}', '${escapeQuotes(regions[departement.region].nom)}', '${escapeQuotes(regions[departement.region].charniere)}')`
      }).join(', ') +
`${tomCom}
ON CONFLICT (code_insee)
DO UPDATE SET
  nom = EXCLUDED.nom,
  charniere = EXCLUDED.charniere,
  region = EXCLUDED.region,
  charniere_region = EXCLUDED.charniere_region
;`
    console.log(sourcecode)
    const targetFile = path.resolve(__dirname, 'data/import_departements.sql')
    fs.writeFileSync(targetFile, sourcecode)
  } catch (err) {
    console.log('Error:', err)
  }
})()

function escapeQuotes (value) {
  return value.replace("'", "''")
}

// reference: https://insee.fr/fr/information/3718969#tncc
const charnieres = {
  0: 'de ',
  1: "d'",
  2: 'du ',
  3: 'de la ',
  4: 'des ',
  5: "de l'",
  6: "des ",
  7: "de las ",
  8: "de los ",
}

function getCharniere (typeLiaison) {
  const charniere = charnieres[typeLiaison]
  if (!charniere) {
    throw new Error(`Le type de liaison ${typeLiaison} est inconnu`)
  }
  return charniere
}

function fetchJSON(url) {
  return new Promise((resolve, reject) => {
    https.get(url, resp => {
      let data = ''
      resp.on('data', chunk => {
        data += chunk
      })

      resp.on('end', () => {
        try {

          resolve(JSON.parse(data))
        } catch (err) {
          reject(err)
        }
      })
    })
    .on('error', err => {
      reject(err)
    })
  })
}
