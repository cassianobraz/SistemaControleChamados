<script setup lang="ts">
import { reactive, watch } from 'vue'

import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import type { Responsible } from '@/types/responsible'
import type { Priority, Status, Ticket } from '@/types/ticket'
import { PRIORITY_OPTIONS, STATUS_OPTIONS } from '@/utils/ticketMeta'

export interface TicketFormValues {
  title: string
  description: string
  priority: Priority
  status: Status
  responsibleId: string
}

const props = defineProps<{
  ticket: Ticket | null
  responsibles: readonly Responsible[]
  submitting: boolean
  errorMessage: string
}>()

const open = defineModel<boolean>({ required: true })

const emit = defineEmits<{
  submit: [values: TicketFormValues]
}>()

function blankForm(): TicketFormValues {
  return { title: '', description: '', priority: 'media', status: 'aberto', responsibleId: '' }
}

const form = reactive<TicketFormValues>(blankForm())

watch(
  open,
  (isOpen) => {
    if (!isOpen) return

    if (props.ticket) {
      form.title = props.ticket.title
      form.description = props.ticket.description
      form.priority = props.ticket.priority
      form.status = props.ticket.status
      form.responsibleId = props.ticket.responsible_id
      return
    }

    Object.assign(form, blankForm())
  },
  { immediate: true },
)

function handleSubmit() {
  emit('submit', { ...form })
}
</script>

<template>
  <BaseModal v-model="open" :title="ticket ? 'Editar chamado' : 'Abrir novo chamado'">
    <form id="ticket-form" class="space-y-4" @submit.prevent="handleSubmit">
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-600" for="ticket-title">Título</label>
        <input
          id="ticket-title"
          v-model="form.title"
          type="text"
          required
          maxlength="200"
          placeholder="Ex.: Impressora do 2º andar não funciona"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        />
      </div>

      <div>
        <label class="mb-1 block text-xs font-medium text-slate-600" for="ticket-description">Descrição</label>
        <textarea
          id="ticket-description"
          v-model="form.description"
          required
          rows="4"
          placeholder="Descreva o que está acontecendo..."
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        />
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-600" for="ticket-priority">Prioridade</label>
          <select
            id="ticket-priority"
            v-model="form.priority"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
          >
            <option v-for="option in PRIORITY_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>

        <div v-if="ticket">
          <label class="mb-1 block text-xs font-medium text-slate-600" for="ticket-status">Status</label>
          <select
            id="ticket-status"
            v-model="form.status"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
          >
            <option v-for="option in STATUS_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>

        <div v-else>
          <label class="mb-1 block text-xs font-medium text-slate-600" for="ticket-responsible">Responsável</label>
          <select
            id="ticket-responsible"
            v-model="form.responsibleId"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
          >
            <option value="">Distribuir automaticamente</option>
            <option v-for="responsible in responsibles" :key="responsible.id" :value="responsible.id">
              {{ responsible.name }}
            </option>
          </select>
        </div>
      </div>

      <p v-if="errorMessage" class="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
    </form>

    <template #footer>
      <BaseButton variant="secondary" type="button" :disabled="submitting" @click="open = false">
        Cancelar
      </BaseButton>
      <BaseButton type="submit" form="ticket-form" :disabled="submitting">
        {{ submitting ? 'Salvando...' : 'Salvar' }}
      </BaseButton>
    </template>
  </BaseModal>
</template>
