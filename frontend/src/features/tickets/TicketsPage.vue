<script setup lang="ts">
import { onMounted, ref } from 'vue'

import BaseButton from '@/components/ui/BaseButton.vue'
import PaginationControls from '@/components/ui/PaginationControls.vue'
import { useResponsibles } from '@/composables/useResponsibles'
import { useTickets } from '@/composables/useTickets'
import ResponsibleFormModal, { type ResponsibleFormValues } from '@/features/tickets/ResponsibleFormModal.vue'
import ResponsibleSidebar from '@/features/tickets/ResponsibleSidebar.vue'
import TicketDetailModal from '@/features/tickets/TicketDetailModal.vue'
import TicketFilters from '@/features/tickets/TicketFilters.vue'
import TicketFormModal, { type TicketFormValues } from '@/features/tickets/TicketFormModal.vue'
import TicketList from '@/features/tickets/TicketList.vue'
import type { Responsible } from '@/types/responsible'
import type { Ticket } from '@/types/ticket'

const {
  tickets,
  total,
  totalPages,
  page,
  loading,
  listError,
  filters,
  fetchTickets,
  goToPage,
  applyFilters,
  createTicket,
  updateTicket,
  assignTicket,
} = useTickets()

const {
  responsibles,
  loading: responsiblesLoading,
  fetchResponsibles,
  createResponsible,
  updateResponsible,
  deleteResponsible,
} = useResponsibles()

const isFormOpen = ref(false)
const editingTicket = ref<Ticket | null>(null)
const formSubmitting = ref(false)
const formError = ref('')

const isDetailOpen = ref(false)
const viewingTicket = ref<Ticket | null>(null)

const isResponsibleFormOpen = ref(false)
const editingResponsible = ref<Responsible | null>(null)
const responsibleFormSubmitting = ref(false)
const responsibleFormError = ref('')

onMounted(() => {
  fetchTickets()
  fetchResponsibles()
})

function openCreateForm() {
  editingTicket.value = null
  formError.value = ''
  isFormOpen.value = true
}

function openEditForm(ticket: Ticket) {
  editingTicket.value = ticket
  formError.value = ''
  isFormOpen.value = true
}

function openTicketDetail(ticket: Ticket) {
  viewingTicket.value = ticket
  isDetailOpen.value = true
}

function editFromDetail(ticket: Ticket) {
  isDetailOpen.value = false
  openEditForm(ticket)
}

async function handleFormSubmit(values: TicketFormValues) {
  formSubmitting.value = true
  formError.value = ''

  const result = editingTicket.value
    ? await updateTicket(editingTicket.value.id, {
        title: values.title,
        description: values.description,
        priority: values.priority,
        status: values.status,
      })
    : await createTicket({
        title: values.title,
        description: values.description,
        priority: values.priority,
        responsible_id: values.responsibleId || undefined,
      })

  formSubmitting.value = false

  if (!result.ok) {
    formError.value = result.message
    return
  }

  isFormOpen.value = false
  fetchResponsibles()
}

async function handleAssign(ticketId: string, responsibleId: string) {
  const result = await assignTicket(ticketId, responsibleId)
  if (result.ok) fetchResponsibles()
}

function openCreateResponsibleForm() {
  editingResponsible.value = null
  responsibleFormError.value = ''
  isResponsibleFormOpen.value = true
}

function openEditResponsibleForm(responsible: Responsible) {
  editingResponsible.value = responsible
  responsibleFormError.value = ''
  isResponsibleFormOpen.value = true
}

async function handleResponsibleFormSubmit(values: ResponsibleFormValues) {
  responsibleFormSubmitting.value = true
  responsibleFormError.value = ''

  const result = editingResponsible.value
    ? await updateResponsible(editingResponsible.value.id, values)
    : await createResponsible({ name: values.name, email: values.email })

  responsibleFormSubmitting.value = false

  if (!result.ok) {
    responsibleFormError.value = result.message
    return
  }

  isResponsibleFormOpen.value = false
}

async function handleDeleteResponsible(responsible: Responsible) {
  if (!window.confirm(`Remover o responsável "${responsible.name}"?`)) return

  const result = await deleteResponsible(responsible.id)
  if (!result.ok) window.alert(result.message)
}
</script>

<template>
  <main class="mx-auto max-w-6xl space-y-4 px-4 py-6 sm:px-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h2 class="text-lg font-semibold text-slate-900">Chamados</h2>
      <BaseButton @click="openCreateForm">+ Abrir chamado</BaseButton>
    </div>

    <TicketFilters
      v-model:status="filters.status"
      v-model:priority="filters.priority"
      v-model:responsible-id="filters.responsibleId"
      v-model:search="filters.search"
      v-model:sort="filters.sort"
      :responsibles="responsibles"
      @apply="applyFilters"
    />

    <p v-if="listError" class="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ listError }}</p>

    <div class="grid grid-cols-1 gap-4 lg:grid-cols-4">
      <div class="space-y-4 lg:col-span-3">
        <TicketList
          :tickets="tickets"
          :responsibles="responsibles"
          :loading="loading"
          @view="openTicketDetail"
          @edit="openEditForm"
          @assign="handleAssign"
        />
        <p class="text-sm text-slate-500">{{ total }} chamado(s) encontrado(s)</p>
        <PaginationControls :page="page" :total-pages="totalPages" :disabled="loading" @change="goToPage" />
      </div>

      <ResponsibleSidebar
        :responsibles="responsibles"
        :loading="responsiblesLoading"
        @create="openCreateResponsibleForm"
        @edit="openEditResponsibleForm"
        @delete="handleDeleteResponsible"
      />
    </div>

    <TicketDetailModal v-model="isDetailOpen" :ticket="viewingTicket" @edit="editFromDetail" />

    <TicketFormModal
      v-model="isFormOpen"
      :ticket="editingTicket"
      :responsibles="responsibles"
      :submitting="formSubmitting"
      :error-message="formError"
      @submit="handleFormSubmit"
    />

    <ResponsibleFormModal
      v-model="isResponsibleFormOpen"
      :responsible="editingResponsible"
      :submitting="responsibleFormSubmitting"
      :error-message="responsibleFormError"
      @submit="handleResponsibleFormSubmit"
    />
  </main>
</template>
