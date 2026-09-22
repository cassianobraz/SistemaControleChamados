import { http } from '@/api/http'
import type {
  AssignTicketPayload,
  CreateTicketPayload,
  ListTicketsParams,
  Ticket,
  TicketListResponse,
  UpdateTicketPayload,
} from '@/types/ticket'

export function listTickets(params: ListTicketsParams) {
  return http.get<TicketListResponse>('/api/v1/tickets', { params }).then((res) => res.data)
}

export function getTicket(id: string) {
  return http.get<Ticket>(`/api/v1/tickets/${id}`).then((res) => res.data)
}

export function createTicket(payload: CreateTicketPayload) {
  return http.post<Ticket>('/api/v1/tickets', payload).then((res) => res.data)
}

export function updateTicket(id: string, payload: UpdateTicketPayload) {
  return http.put<Ticket>(`/api/v1/tickets/${id}`, payload).then((res) => res.data)
}

export function assignTicket(id: string, payload: AssignTicketPayload) {
  return http.patch<Ticket>(`/api/v1/tickets/${id}/assign`, payload).then((res) => res.data)
}
