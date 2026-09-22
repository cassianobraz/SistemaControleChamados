export interface Responsible {
  id: string
  name: string
  email: string
  active: boolean
  open_tickets: number
}

export interface CreateResponsiblePayload {
  name: string
  email: string
}

export interface UpdateResponsiblePayload {
  name: string
  email: string
  active: boolean
}
