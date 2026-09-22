export type Priority = 'baixa' | 'media' | 'alta'

export type Status = 'aberto' | 'em_andamento' | 'resolvido' | 'fechado'

export type SortOrder = 'created_at_desc' | 'created_at_asc' | 'priority_desc'

export interface Ticket {
  id: string
  title: string
  description: string
  priority: Priority
  status: Status
  responsible_id: string
  responsible_name: string
  created_at: string
  updated_at: string
}

export interface TicketListResponse {
  data: Ticket[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface CreateTicketPayload {
  title: string
  description: string
  priority: Priority
  responsible_id?: string
}

export interface UpdateTicketPayload {
  title: string
  description: string
  priority: Priority
  status: Status
}

export interface AssignTicketPayload {
  responsible_id?: string
}

export interface ListTicketsParams {
  status?: Status
  priority?: Priority
  responsible_id?: string
  search?: string
  sort?: SortOrder
  page?: number
  page_size?: number
}
