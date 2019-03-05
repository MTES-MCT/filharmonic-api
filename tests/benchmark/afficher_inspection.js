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
  const inspectionId = __ITER + 1
  const res = http.get(`http://localhost:5000/inspections/${inspectionId}`, httpParams)
  check(res, {
    'is status 200': r => r.status === 200
  })
}
