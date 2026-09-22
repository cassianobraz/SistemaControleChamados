import { reactive, readonly, ref } from 'vue'

import * as ticketsApi from '@/api/tickets'
import { extractErrorMessage } from '@/api/http'
import type {
  CreateTicketPayload,
  ListTicketsParams,
  Priority,
  SortOrder,
  Status,
  Ticket,
  UpdateTicketPayload,
} from '@/types/ticket'

export interface TicketFilters {
  status: Status | ''
  priority: Priority | ''
  responsibleId: string
  search: string
  sort: SortOrder
}

export type ActionResult = { ok: true } | { ok: false; message: string }

const DEFAULT_PAGE_SIZE = 10

export function useTickets() {
  const _tickets = ref<Ticket[]>([])
  const _total = ref(0)
  const _totalPages = ref(0)
  const _page = ref(1)
  const _loading = ref(false)
  const _listError = ref('')

  const filters = reactive<TicketFilters>({
    status: '',
    priority: '',
    responsibleId: '',
    search: '',
    sort: 'created_at_desc',
  })

  function buildParams(): ListTicketsParams {
    const params: ListTicketsParams = {
      page: _page.value,
      page_size: DEFAULT_PAGE_SIZE,
      sort: filters.sort,
    }

    if (filters.status) params.status = filters.status
    if (filters.priority) params.priority = filters.priority
    if (filters.responsibleId) params.responsible_id = filters.responsibleId
    if (filters.search.trim()) params.search = filters.search.trim()

    return params
  }

  async function fetchTickets(): Promise<void> {
    _loading.value = true
    _listError.value = ''

    try {
      const response = await ticketsApi.listTickets(buildParams())
      _tickets.value = response.data
      _total.value = response.total
      _totalPages.value = response.total_pages
      _page.value = response.page
    } catch (error) {
      _listError.value = extractErrorMessage(error)
    } finally {
      _loading.value = false
    }
  }

  function goToPage(nextPage: number): Promise<void> {
    _page.value = nextPage
    return fetchTickets()
  }

  function applyFilters(): Promise<void> {
    _page.value = 1
    return fetchTickets()
  }

  async function runAction(action: () => Promise<unknown>): Promise<ActionResult> {
    try {
      await action()
      await fetchTickets()
      return { ok: true }
    } catch (error) {
      return { ok: false, message: extractErrorMessage(error) }
    }
  }

  function createTicket(payload: CreateTicketPayload): Promise<ActionResult> {
    return runAction(() => ticketsApi.createTicket(payload))
  }

  function updateTicket(id: string, payload: UpdateTicketPayload): Promise<ActionResult> {
    return runAction(() => ticketsApi.updateTicket(id, payload))
  }

  function assignTicket(id: string, responsibleId: string): Promise<ActionResult> {
    return runAction(() => ticketsApi.assignTicket(id, { responsible_id: responsibleId }))
  }

  return {
    tickets: readonly(_tickets),
    total: readonly(_total),
    totalPages: readonly(_totalPages),
    page: readonly(_page),
    loading: readonly(_loading),
    listError: readonly(_listError),
    filters,
    fetchTickets,
    goToPage,
    applyFilters,
    createTicket,
    updateTicket,
    assignTicket,
  }
}
