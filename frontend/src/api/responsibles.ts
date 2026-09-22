import { http } from '@/api/http'
import type { CreateResponsiblePayload, Responsible, UpdateResponsiblePayload } from '@/types/responsible'

export function listResponsibles() {
  return http.get<Responsible[]>('/api/v1/responsibles').then((res) => res.data)
}

export function createResponsible(payload: CreateResponsiblePayload) {
  return http.post<Responsible>('/api/v1/responsibles', payload).then((res) => res.data)
}

export function updateResponsible(id: string, payload: UpdateResponsiblePayload) {
  return http.put<Responsible>(`/api/v1/responsibles/${id}`, payload).then((res) => res.data)
}

export function deleteResponsible(id: string) {
  return http.delete<void>(`/api/v1/responsibles/${id}`).then((res) => res.data)
}
