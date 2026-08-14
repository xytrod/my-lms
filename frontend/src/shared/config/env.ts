const fallbackApiUrl = 'http://localhost:3000'

export const env = {
  apiUrl: (import.meta.env.VITE_API_URL || fallbackApiUrl).replace(/\/$/, ''),
} as const
