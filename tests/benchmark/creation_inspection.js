import http from 'k6/http'
import { sleep, check } from 'k6'
import config from './config.js'

const httpParams = {
  headers: {
    Authorization: config.authToken,
    'Content-Type': 'application/json',
  }
}

export default function() {
  const etablissementId = __ITER + 1
  const inspection = {
    date: '2019-04-01',
    type: 'ponctuel',
    annonce: true,
    origine: 'circonstancielle',
    circonstance: 'plainte',
    detail_circonstance: '',
    contexte: 'Contexte inspection',
    inspecteurs: [
      {
        id: 3,
      },
      {
        id: 4,
      }
    ],
    themes: [
      'Rejets Air',
      'COV'
    ],
    etablissement_id: etablissementId
  }

  let res = http.post('http://localhost:5000/inspections', JSON.stringify(inspection), httpParams)
  check(res, {
    'is status 200': r => r.status === 200
  })
  const inspectionId = res.json().id
  res = http.get(`http://localhost:5000/inspections/${inspectionId}`, httpParams)
  check(res, {
    'is status 200': r => r.status === 200
  })
}
