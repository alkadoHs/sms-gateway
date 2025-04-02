import axios from 'axios'

const client = axios.create({
  baseURL: process.env.API_BASE_URL || 'http://localhost:8000',
  headers: {
    'X-Client-Version': process.env.GITHUB_SHA || 'dev',
  },
})

export function setAuthHeader(token: string | null) {
  client.defaults.headers.common.Authorization = 'Bearer ' + token
}

export function setApiKey(apiKey: string | null) {
  client.defaults.headers.common['x-api-key'] = apiKey ?? ''
}

export default client
