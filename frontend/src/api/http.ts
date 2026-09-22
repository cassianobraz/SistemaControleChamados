import axios from 'axios'

const baseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export const http = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const message = (error.response?.data as { error?: string } | undefined)?.error
    if (message) return message
  }

  return 'Não foi possível completar a operação. Tente novamente.'
}
