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
  const res = http.get('http://localhost:5000/etablissements?s3ic=123&nom=tar&adresse=tar&page=2531', httpParams)
  check(res, {
    'is status 200': r => r.status === 200
  })
}
